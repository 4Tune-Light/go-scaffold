package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/rizky/go-scaffold/internal/config"
	"github.com/rizky/go-scaffold/internal/greeter"
	"github.com/rizky/go-scaffold/internal/health"
	"github.com/rizky/go-scaffold/internal/middleware"
	"github.com/rizky/go-scaffold/internal/server"
	"github.com/rizky/go-scaffold/internal/telemetry"
	"github.com/rizky/go-scaffold/pkg/database"
	"github.com/rizky/go-scaffold/pkg/idempotency"
)

func main() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg := config.MustLoad("")

	ctx := context.Background()

	tp, mp, err := telemetry.Init(ctx, telemetry.Config{
		ServiceName: cfg.OTel.ServiceName,
		Environment: cfg.OTel.Environment,
		Endpoint:    cfg.OTel.Endpoint,
		Insecure:    cfg.OTel.Insecure,
		TraceRatio:  cfg.OTel.TraceRatio,
	})
	if err != nil {
		log.Warn().Err(err).Msg("telemetry not available")
	} else {
		defer func() { _ = tp.Shutdown(ctx); _ = mp.Shutdown(ctx) }()
	}

	var checker *health.Checker
	var idempStore *idempotency.Store
	var greeterHandler *greeter.Handler
	var greeterSvc greeter.Service

	if pgPool, err := database.NewPostgresPool(ctx, cfg.PostgresConfig()); err != nil {
		log.Warn().Err(err).Msg("PostgreSQL not available")
	} else {
		defer pgPool.Close()
		q := database.NewQuerier(pgPool)
		tx := database.NewTransactor(pgPool)
		greeterSvc = greeter.NewService(greeter.NewRepository(q), tx)
		greeterHandler = greeter.NewHandler(greeterSvc)

		if rdb, err := database.NewRedisClient(ctx, cfg.RedisConfig()); err != nil {
			log.Warn().Err(err).Msg("Redis not available")
			checker = health.New(pgPool, nil)
		} else {
			defer rdb.Close()
			checker = health.New(pgPool, rdb)
			idempStore = idempotency.NewStore(rdb)
		}
	}

	if checker == nil {
		checker = health.New(nil, nil)
	}
	if idempStore == nil {
		idempStore = idempotency.NewStore(nil)
	}

	httpSrv := server.NewHTTPServer("api",
		cfg.Server.HTTP.Host, cfg.Server.HTTP.Port, cfg.Server.HTTP.ReadTimeout)

	registerRoutes(httpSrv.Router(), checker, idempStore, cfg, greeterHandler)

	grpcSrv, err := server.NewGRPCServer("grpc", cfg.Server.GRPC.Host, cfg.Server.GRPC.Port)
	if err != nil {
		log.Warn().Err(err).Msg("gRPC not available")
		grpcSrv = nil
	}

	svr := []server.Server{httpSrv}
	if grpcSrv != nil {
		svr = append(svr, grpcSrv)
		if greeterSvc != nil {
			greeter.RegisterGRPC(grpcSrv.Server(), greeterSvc)
		}
	}

	mgr := server.NewManager(svr...)

	log.Info().Msg("API server started")
	if err := mgr.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("server stopped with error")
	}
}

func registerRoutes(r chi.Router, checker *health.Checker, idempStore *idempotency.Store, cfg *config.Config, greeterH *greeter.Handler) {
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORS([]string{"*"}))
	r.Use(middleware.OTelHTTP(cfg.OTel.ServiceName))
	r.Use(middleware.Idempotency(idempStore))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		status := checker.Check(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	})

	if greeterH != nil {
		greeter.RegisterRoutes(r, greeterH, cfg.JWT.Secret)
	}
}
