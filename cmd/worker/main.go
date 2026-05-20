package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/rizky/go-scaffold/internal/config"
	"github.com/rizky/go-scaffold/internal/telemetry"
	"github.com/rizky/go-scaffold/pkg/database"
)

func main() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg := config.MustLoad("")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, mp, err := telemetry.Init(ctx, telemetry.Config{
		ServiceName: cfg.OTel.ServiceName + "-worker",
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

	pgPool, err := database.NewPostgresPool(ctx, cfg.PostgresConfig())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to PostgreSQL")
	}
	defer pgPool.Close()

	rdb, err := database.NewRedisClient(ctx, cfg.RedisConfig())
	if err != nil {
		log.Warn().Err(err).Msg("Redis not available")
	}

	q := database.NewQuerier(pgPool)
	_ = q

	log.Info().
		Str("service", cfg.OTel.ServiceName).
		Msg("Worker started")

	if rdb != nil {
		log.Info().Msg("Redis available — worker can subscribe to streams")
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Info().Msg("shutting down worker")
}
