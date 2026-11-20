# Stack notes

Only idioms and gotchas that commonly go wrong. Always prefer the project's existing patterns.

## Node.js / TypeScript (Express, Fastify, NestJS)
- Validate with a schema library already in the project (zod, valibot, class-validator, TypeBox, AJV); infer TS types from the schema, do not duplicate.
- Express 4 does not catch rejected promises in handlers: use the project's async wrapper or Express 5. Always end with one central error middleware `(err, req, res, next)`.
- Never block the event loop: no sync crypto/fs/JSON of huge payloads in handlers; offload CPU work to worker threads or a queue.
- `fetch` has no default timeout: pass `signal: AbortSignal.timeout(ms)`.
- Use `process.on('SIGTERM')` to close the server and pools; `unhandledRejection` should log and exit, not be swallowed.
- NestJS: guards for authz, pipes for validation (`whitelist: true, forbidNonWhitelisted: true`), interceptors for logging/mapping, exception filters for problem details.
- Config: validate `process.env` at startup with a schema; export a typed config object.

## Python (FastAPI, Django, Flask)
- FastAPI: Pydantic models for request/response (`response_model` prevents leaking fields); dependencies (`Depends`) for auth, DB session, current user. Do not call blocking I/O inside `async def` routes; use `def` routes or async drivers.
- Django: `select_related`/`prefetch_related` against N+1; `transaction.atomic()`; `transaction.on_commit()` to enqueue jobs after commit; DRF serializers with explicit `fields` (never `__all__` for writes); object-level permissions in `get_queryset` filtered by user/tenant.
- Flask: application factory, blueprints, explicit error handlers returning problem JSON.
- `requests`/`httpx` have no default timeout in many cases: always pass `timeout=`.
- Use `Decimal` for money; timezone-aware `datetime` (`datetime.now(timezone.utc)`).

## Go
- Pass `context.Context` as the first parameter through every layer; use `http.NewRequestWithContext`; set `http.Client{Timeout: ...}` (the default client has none).
- Server timeouts: set `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout` on `http.Server`.
- Wrap errors with `%w`; check with `errors.Is/As`; map to HTTP in one place.
- `defer rows.Close()` and check `rows.Err()`; `defer tx.Rollback()` after `BeginTx` (no-op after commit).
- Close response bodies (`defer resp.Body.Close()`), and drain them for connection reuse.
- Run tests with `-race` in CI; avoid goroutines without cancellation (leaks).

## Java / Kotlin (Spring Boot)
- `@Valid` on request bodies with Bean Validation annotations; `@RestControllerAdvice` + `ProblemDetail` (Spring 6+) for RFC 9457 errors.
- `@Transactional` only works on public methods called through the proxy (not self-invocation); keep it on service methods, not controllers; avoid remote calls inside.
- Use `@PreAuthorize` or explicit service checks for object-level authz, not only URL patterns.
- JPA: watch lazy loading N+1 (use fetch joins/entity graphs), disable open-session-in-view (`spring.jpa.open-in-view=false`), never expose entities directly as DTOs.
- Configure timeouts on `RestClient`/`WebClient` and connection pools (HikariCP size modest).
- Use Testcontainers for integration tests against the real DB.

## .NET (ASP.NET Core)
- Minimal APIs or controllers per project convention; validation with DataAnnotations or FluentValidation; `AddProblemDetails()` + exception handler middleware.
- `IHttpClientFactory` for outbound calls (avoids socket exhaustion); configure timeouts and resilience handlers there.
- Pass `CancellationToken` through async calls; never `.Result`/`.Wait()` (deadlocks, thread starvation).
- EF Core: `AsNoTracking()` for reads, explicit `Include` against N+1, migrations reviewed as SQL (`dotnet ef migrations script`).
- Options pattern with `ValidateOnStart()` for configuration; secrets from environment or a secret store, never `appsettings.json`.
- Authorization policies (`[Authorize(Policy=...)]`) plus resource-based checks via `IAuthorizationService` for object-level access.

## Cross-stack test tips
- Integration tests through the HTTP layer catch serialisation, validation and middleware ordering bugs that unit tests miss.
- Use real database engines in tests (containers), not SQLite/H2 substitutes, when queries use engine-specific features.
- Freeze time and seed randomness for deterministic tests; stub outbound HTTP at the client boundary (mock server), not deep inside code.
