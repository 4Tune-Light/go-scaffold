# go-scaffold

> Production-ready Go backend scaffold with OpenTelemetry, PostgreSQL, Redis, and gRPC.

[![Go Version](https://img.shields.io/badge/go-1.26-blue)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Features

- **HTTP server** — chi router with full middleware stack (logger, CORS, recovery, request ID, security headers)
- **JWT Auth** — Bearer token authentication with role-based access control
- **Rate Limiter** — token bucket algorithm per-IP rate limiting
- **Idempotency** — `Idempotency-Key` header support via Redis-backed store
- **gRPC server** — with reflection support for development tools
- **Config** — viper-based YAML config with automatic environment variable override
- **Database** — PostgreSQL via pgx connection pool, Redis via go-redis
- **Transactions** — context-aware Queryer/Transactor pattern for multi-step operations
- **Retry** — exponential backoff utility for transient failures
- **Observability** — OpenTelemetry traces + metrics, exportable to Prometheus and Jaeger
- **Architecture Decisions** — ADR docs explaining key design choices
- **Graceful shutdown** — clean signal handling for SIGINT/SIGTERM/SIGQUIT
- **Docker** — multi-stage builds + docker-compose for local development
- **Standard response** — consistent JSON API response format
- **Domain Pattern** — per-domain entity, dto, errors, repository with concrete implementation

## Project Structure

```
├── cmd/
│   ├── api/                  HTTP + gRPC server entrypoint
│   └── worker/               Async worker entrypoint (config + DB/Redis)
├── docs/
│   └── adr/                  Architecture Decision Records
├── internal/
│   ├── config/               Config loader and type definitions
│   ├── greeter/              Example domain (handler → service → repository)
│   │   ├── dto/              Request/response types
│   │   ├── entity.go         Domain model with business methods
│   │   ├── errors.go         Sentinel errors
│   │   ├── handler.go        HTTP handler
│   │   ├── repository.go     Concrete Postgres repository
│   │   └── service.go        Business logic
│   ├── middleware/
│   │   ├── auth.go           JWT authentication + RBAC
│   │   ├── cors.go           CORS wrapper
│   │   ├── idempotency.go    Idempotency-Key middleware
│   │   ├── logger.go         Structured request logging
│   │   ├── otel.go           OpenTelemetry HTTP tracing
│   │   ├── ratelimit.go      Token bucket rate limiter
│   │   ├── recovery.go       Panic recovery
│   │   ├── request_id.go     X-Request-Id propagation
│   │   └── security.go       Security headers (CSP, XSS, frame options)
│   ├── server/               Server abstractions (HTTP, gRPC, Manager)
│   └── telemetry/            OpenTelemetry SDK initialization
├── pkg/
│   ├── database/
│   │   ├── postgres.go       PostgreSQL pool factory
│   │   ├── redis.go          Redis client factory
│   │   ├── querier.go        Context-aware query wrapper (tx detection)
│   │   └── transactor.go     Transaction helper (WithinTransaction)
│   ├── idempotency/          Redis-backed idempotency key store
│   ├── response/             Standard JSON API response helpers
│   └── retry/                Exponential backoff utility
├── configs/                  YAML configuration files
├── deploy/                   Docker, Compose, and observability configs
├── migrations/               SQL migration files
├── .env.example              Environment variable template
├── Makefile                  Build automation
└── README.md
```

## Quick Start

```bash
# 1. Clone and enter
git clone https://github.com/rizky/go-scaffold.git
cd go-scaffold

# 2. Copy environment
cp .env.example .env

# 3. Start dependencies (PostgreSQL, Redis, OTel)
docker compose -f deploy/docker-compose.yml up postgres redis otel-collector -d

# 4. Run API server
make run-api

# 5. Verify
curl http://localhost:8080/health
# {"status":"ok"}
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make build` | Build all service binaries |
| `make run-api` | Run API server locally |
| `make run-worker` | Run worker locally |
| `make dev` | Start full stack with Docker Compose |
| `make dev-down` | Stop Docker Compose services |
| `make test` | Run all tests with race detection |
| `make lint` | Run golangci-lint |
| `make tidy` | Run go mod tidy |

## Middleware Stack

```
chi RequestID → RequestID → Logger → Recovery → SecurityHeaders → CORS → OTelHTTP → Idempotency → Routes
```

All middleware can be enabled/disabled per-route group via `r.Group()`.

## Domain Pattern

Each domain (`internal/{domain}/`) follows a consistent structure:

```go
entity.go      // Domain model with validation/business methods
errors.go      // Sentinel errors (ErrNotFound, ErrInvalidInput)
repository.go  // Data access layer with *database.Querier
service.go     // Business logic with interface for testability
handler.go     // HTTP handler with error mapping
dto/           // Request/response DTO types
```

### Transaction Support

Use `database.Transactor.WithinTransaction` for multi-step operations:

```go
err := tx.WithinTransaction(ctx, func(txCtx context.Context) error {
    repo.Save(txCtx, ...)       // automatically uses the transaction
    repo.Update(txCtx, ...)     // same transaction
    return nil
})
```

The `Querier` automatically detects the transaction from context — no need to pass `pgx.Tx` explicitly.

## Architecture Decisions

Key design choices are documented as ADRs in `docs/adr/`:

| ADR | Decision |
|-----|----------|
| 001 | Use chi as HTTP router |
| 002 | Use pgx/v5 over database/sql |
| 003 | Use OpenTelemetry for observability |
| 004 | Modular monolith with layered architecture |

## Usage as a Template

1. Fork or copy this repository
2. Update the module name:
   ```bash
   go mod edit -module github.com/your-username/your-project
   ```
3. Add your domain packages under `internal/` following the greeter pattern
4. Replace sample migrations under `migrations/` with your schema
5. Wire your handlers in `cmd/api/main.go`
6. Update `configs/config.yaml` with your settings
7. Build and deploy

## Deployment

### Docker Compose (single VM)

```bash
export APP_ENV=production
docker compose -f deploy/docker-compose.yml up --build -d
```

### Production Checklist

- [ ] Update `JWT_SECRET` and other secrets via environment variables
- [ ] Set `APP_ENV=production` and `APP_DEBUG=false`
- [ ] Configure PostgreSQL with proper credentials and SSL
- [ ] Set up Prometheus + Grafana for monitoring
- [ ] Add reverse proxy (nginx / Caddy) for TLS termination
- [ ] Configure database backups

## Dependencies

- [chi](https://github.com/go-chi/chi) — HTTP router
- [pgx](https://github.com/jackc/pgx) — PostgreSQL driver
- [go-redis](https://github.com/redis/go-redis) — Redis client
- [golang-jwt](https://github.com/golang-jwt/jwt) — JWT authentication
- [viper](https://github.com/spf13/viper) — Configuration management
- [zerolog](https://github.com/rs/zerolog) — Structured logging
- [OpenTelemetry](https://opentelemetry.io/) — Observability SDK
- [gRPC](https://grpc.io/) — RPC framework

## License

MIT
