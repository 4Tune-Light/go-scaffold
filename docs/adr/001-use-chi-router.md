# ADR 001: Use chi as HTTP Router

**Status:** Accepted

**Context:**
The scaffold needed an HTTP router that is idiomatic, composable, and compatible with the standard `net/http` library.

**Options considered:**
1. **gorilla/mux** — Popular but archived (v1.x is in maintenance mode).
2. **gin** — Fast but uses `gin.Context` instead of standard `http.Handler`, making it harder to swap or reuse middleware.
3. **echo** — Fast but similar coupling issue as gin.
4. **chi** — Lightweight, idiomatic, fully compatible with `net/http`, supports middleware chaining, context propagation, and URL parameters.

**Decision:**
Use `go-chi/chi/v5` as the HTTP router.

**Consequences:**
- Middleware is standard `func(http.Handler) http.Handler` — reusable without chi dependency.
- Easy to test with `httptest.NewRecorder` and `chi.NewRouter`.
- URL parameters via `chi.URLParam(r, "name")`.
- No framework lock-in — handlers are standard `http.HandlerFunc`.
