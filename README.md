# go-scaffold

> Production-ready Go backend scaffold with OpenTelemetry, PostgreSQL, Redis, and gRPC.

[![Go Version](https://img.shields.io/badge/go-1.26-blue)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Features

- **HTTP server** — chi router with full middleware stack (logger, CORS, recovery, request ID)
- **gRPC server** — with reflection support for development tools
- **Config** — viper-based YAML config with automatic environment variable override
- **Database** — PostgreSQL via pgx connection pool, Redis via go-redis
- **Observability** — OpenTelemetry traces + metrics, exportable to Prometheus and Jaeger
- **Graceful shutdown** — clean signal handling for SIGINT/SIGTERM/SIGQUIT
- **Docker** — multi-stage builds + docker-compose for local development
- **Standard response** — consistent JSON API response format

## Project Structure

```
├── cmd/
│   ├── api/              HTTP + gRPC server entrypoint
│   └── worker/           Async worker entrypoint
├── internal/
│   ├── config/           Config loader and type definitions
│   ├── middleware/        HTTP middleware (logger, CORS, recovery, request ID)
│   ├── server/            Server abstractions (HTTP, gRPC, Manager)
│   └── telemetry/         OpenTelemetry SDK initialization
├── pkg/
│   ├── database/          PostgreSQL and Redis client factories
│   └── response/          Standard JSON API response helpers
├── configs/               YAML configuration files
├── deploy/                Docker, Compose, and observability configs
├── migrations/            SQL migration files
├── .env.example           Environment variable template
├── Makefile               Build automation
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

## Usage as a Template

1. Fork or copy this repository
2. Update the module name:
   ```bash
   go mod edit -module github.com/your-username/your-project
   ```
3. Add your domain packages under `internal/` (e.g., `internal/user/`, `internal/order/`)
4. Replace sample migrations under `migrations/` with your schema
5. Wire your handlers in `cmd/api/main.go` with your domain packages
6. Update `configs/config.yaml` with your settings
7. Build and deploy

## Deployment

### Docker Compose (single VM)

```bash
# Set production environment
export APP_ENV=production

# Deploy full stack
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
- [viper](https://github.com/spf13/viper) — Configuration management
- [zerolog](https://github.com/rs/zerolog) — Structured logging
- [OpenTelemetry](https://opentelemetry.io/) — Observability SDK
- [gRPC](https://grpc.io/) — RPC framework

## License

MIT
