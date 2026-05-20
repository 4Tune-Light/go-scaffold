# ADR 003: Use OpenTelemetry for Observability

**Status:** Accepted

**Context:**
The scaffold needs distributed tracing and metrics to support debugging, monitoring, and performance analysis in production.

**Options considered:**
1. **Prometheus client directly** — Provides metrics but no tracing support.
2. **Jaeger client directly** — Provides tracing but no metrics support.
3. **OpenTelemetry** — Unified standard for traces, metrics, and logs; vendor-agnostic; can export to Prometheus, Jaeger, Datadog, etc.

**Decision:**
Use OpenTelemetry SDK with OTLP gRPC exporter.

**Consequences:**
- Traces and metrics are exported via OTLP to an OpenTelemetry Collector.
- The collector can fan-out to Prometheus (metrics) and Jaeger (traces).
- OTel SDK is initialized at startup with graceful degradation (warns on failure, server still starts).
- Service name and environment are configurable via config.
- Sampling rate is configurable (default 10%).
