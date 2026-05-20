# ADR 002: Use pgx/v5 Instead of database/sql

**Status:** Accepted

**Context:**
The scaffold needs PostgreSQL access. The standard library provides `database/sql`, but pgx offers native PostgreSQL features.

**Options considered:**
1. **database/sql + lib/pq** — Standard approach but lacks native PostgreSQL types, connection pooling, and COPY support.
2. **database/sql + pgx stdlib** — Bridges pgx through database/sql, losing pgx-specific features.
3. **pgx/v5 directly** — Native PostgreSQL driver with connection pooling (`pgxpool`), PostgreSQL type support, `pgx.Batch`, `COPY`, and `LISTEN/NOTIFY`.

**Decision:**
Use `github.com/jackc/pgx/v5` directly via `pgxpool.Pool`.

**Consequences:**
- Direct access to PostgreSQL types (UUID, JSONB, arrays, numeric).
- Built-in connection pool with configurable `MaxConns`, `MinConns`, `MaxConnLifetime`.
- Support for `pgx.Batch` for batch operations.
- Transaction support via `pgx.Tx`.
- Slight learning curve for teams used to `database/sql`.
