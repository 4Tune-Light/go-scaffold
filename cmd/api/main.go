package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/rizky/go-scaffold/internal/config"
	"github.com/rizky/go-scaffold/internal/middleware"
	"github.com/rizky/go-scaffold/internal/server"
	"github.com/rizky/go-scaffold/internal/telemetry"
	"github.com/rizky/go-scaffold/pkg/database"
)

func main() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg := loadConfig()

	ctx := context.Background()

	tp, mp, err := telemetry.Init(ctx, telemetry.Config{
		ServiceName: cfg.OTel.ServiceName,
		Environment: cfg.OTel.Environment,
		Endpoint:    cfg.OTel.Endpoint,
		Insecure:    cfg.OTel.Insecure,
		TraceRatio:  cfg.OTel.TraceRatio,
	})
	if err != nil {
		log.Warn().Err(err).Msg("telemetry not available, continuing without it")
	} else {
		defer func() {
			_ = tp.Shutdown(ctx)
			_ = mp.Shutdown(ctx)
		}()
	}

	initPostgres(ctx, cfg)
	initRedis(ctx, cfg)

	httpSrv := server.NewHTTPServer(
		"api",
		cfg.Server.HTTP.Host,
		cfg.Server.HTTP.Port,
		cfg.Server.HTTP.ReadTimeout,
	)
	registerRoutes(httpSrv.Router())

	grpcSrv, err := server.NewGRPCServer(
		"grpc",
		cfg.Server.GRPC.Host,
		cfg.Server.GRPC.Port,
	)
	if err != nil {
		log.Warn().Err(err).Msg("gRPC server not available, continuing without it")
		grpcSrv = nil
	}

	svr := []server.Server{httpSrv}
	if grpcSrv != nil {
		svr = append(svr, grpcSrv)
	}

	mgr := server.NewManager(svr...)

	log.Info().Msg("API server started")
	if err := mgr.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("server stopped with error")
	}
}

func loadConfig() *config.Config {
	cfgPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	return cfg
}

func buildPostgresDSN(cfg *config.Config) database.PostgresConfig {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password,
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
		cfg.Database.Postgres.DBName,
		cfg.Database.Postgres.SSLMode,
	)
	return database.PostgresConfig{
		DSN:             dsn,
		MaxOpenConns:    cfg.Database.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Database.Postgres.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.Postgres.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.Postgres.ConnMaxIdleTime,
	}
}

func initPostgres(ctx context.Context, cfg *config.Config) {
	pgCfg := buildPostgresDSN(cfg)
	if _, err := database.NewPostgresPool(ctx, pgCfg); err != nil {
		log.Warn().Err(err).Msg("PostgreSQL not available, continuing without it")
	}
}

func initRedis(ctx context.Context, cfg *config.Config) {
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	if _, err := database.NewRedisClient(ctx, database.RedisConfig{
		Addr:         redisAddr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	}); err != nil {
		log.Warn().Err(err).Msg("Redis not available, continuing without it")
	}
}

func registerRoutes(r chi.Router) {
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", healthHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
