# ADR 004: Modular Monolith with Layered Architecture

**Status:** Accepted

**Context:**
The scaffold needs an architecture that is maintainable, testable, and can scale from simple CRUD to complex business logic without premature decomposition into microservices.

**Options considered:**
1. **Flat structure** — Simple but leads to god packages and circular dependencies.
2. **Hexagonal (ports and adapters)** — Clean separation but introduces abstraction overhead for simple domains.
3. **Modular monolith** — Domain-based packages with clear boundaries; each domain has its own handler, service, repository, entity, and errors. Can extract domains into separate services later if needed.
4. **Microservices** — Premature for a scaffold; adds operational complexity.

**Decision:**
Use a modular monolith with layered architecture per domain:

```
internal/{domain}/
    handler.go      — HTTP layer (chi)
    service.go      — Business logic
    repository.go   — Data access (pgx)
    entity.go       — Domain model with business methods
    errors.go       — Sentinel errors
    dto/            — Request/response types
```

**Consequences:**
- Each domain is self-contained with clear boundaries.
- Domains communicate via interfaces (not direct coupling).
- Easy to test — each layer can be mocked independently.
- Future extraction to microservices is possible by reusing the domain package.
- Shared infrastructure lives in `pkg/` (database, response, retry, idempotency).
