---
name: build-backend
description: Implement backend APIs and services to production standard (status codes, RFC 9457 errors, pagination, versioning, idempotency keys, validation, authz, transactions, timeouts and retries, background jobs, tests). Use when adding or changing endpoints, handlers, services, webhooks, queue consumers or jobs. Not for schema or migration work (use `build-database`) or security review (use `secure-code-review`).
---

# Backend Development

Ship server-side changes that fit the existing codebase, keep the API contract stable, and fail safely.

## Procedure
1. **Trace the existing path** for a similar endpoint or job: route -> middleware (auth, validation, rate limit) -> handler -> service -> repository/client -> response mapping. Copy its structure, DI, error types, logging and test style. If none exists, use `plan-codebase-analysis` first.
2. **Define the contract before code**: method and path, request schema, response schema, status codes, error codes, auth requirement, idempotency, pagination. If an OpenAPI/proto/GraphQL schema is the source of truth, change it first and regenerate. Read `references/api-design.md` when designing or changing an HTTP API.
3. **Classify the change**: additive (new field, new endpoint) or breaking (removed/renamed field, type change, new required input, changed status code or semantics, tighter validation). Breaking changes need a new version or an expand/contract rollout and the user's approval.
4. **Validate at the boundary** with the project's schema library: types, ranges, lengths, formats, enums, unknown-field policy. Parse into typed domain objects once; inner layers trust them. Return 400/422 with field-level details.
5. **Authenticate in middleware, authorise at the resource.** Check that the caller may perform this action on this specific object (ownership, tenant, role) in the service or query, not only by route. Not found and forbidden for another tenant's object both return 404 to avoid leaking existence (follow project convention if it differs).
6. **Implement business logic** in the service layer, free of transport types. Wrap multi-write operations in one transaction; keep network calls out of open transactions.
7. **Make side effects safe**: idempotency keys for non-idempotent POSTs that clients may retry (payments, orders); unique constraints as the last line of defence; outbox pattern for "save and publish". Read `references/resilience.md` when calling other services, adding retries, webhooks, queues or background jobs.
8. **Handle errors deliberately**: map domain errors to status codes in one place; never leak stack traces, SQL, or internal hostnames; log the full error server-side with a request/correlation ID.
9. **Config via environment**: read and validate all config at startup, fail fast on missing required values, no secrets in code, defaults safe for production.
10. **Observability**: structured logs with request ID, user/tenant ID (not personal data), operation and outcome; metrics for rate, errors, duration on new endpoints; propagate trace context to downstream calls (see `operate-observability`).
11. **Test**: success, validation failure, unauthenticated (401), unauthorised (403/404), not found, conflict/idempotent replay, and dependency failure (timeout, 5xx). Prefer integration tests through the HTTP layer with a real DB (containers) for persistence logic; unit-test pure domain logic.
12. **Run** targeted tests, lint, typecheck; then report contract changes and affected consumers.

## Rules
- Status codes (detail in `references/api-design.md`):
  | Situation | Code |
  |---|---|
  | Created | 201 + `Location` |
  | Accepted for async processing | 202 + status URL |
  | Success, no body | 204 |
  | Malformed or schema-invalid input | 400 (or 422 if the project uses it for semantic validation) |
  | No/invalid credentials | 401 + `WWW-Authenticate` |
  | Authenticated but not allowed | 403 (404 for cross-tenant objects) |
  | State conflict, duplicate, version mismatch | 409 (412 for failed `If-Match`) |
  | Rate limited | 429 + `Retry-After` |
  | Our bug | 500; dependency down/timeout 502/503/504 |
- Error bodies use one format across the API; default to RFC 9457 `application/problem+json` (`type`, `title`, `status`, `detail`, `instance`, plus an `errors` array for fields) unless the project has its own.
- Lists are always paginated with a server-enforced max page size; prefer cursor pagination for large or changing data; sort order must be deterministic (tie-break on ID).
- Never trust client-supplied IDs for ownership, roles, prices or tenant; derive from the session and the database.
- Mass assignment: map allowed fields explicitly; never bind request bodies directly onto ORM entities.
- Every outbound call has a timeout (connect and total). Retry only idempotent operations, with exponential backoff and jitter, max ~3 attempts, and only on retryable errors (timeouts, 429, 502-504).
- Do not hold a DB transaction open across HTTP calls, queue publishes or user input.
- Money as integer minor units or decimal type, never float. Timestamps stored and sent in UTC ISO 8601 with offset.
- Background jobs must be idempotent, have a timeout, bounded retries and a dead-letter/failed state that someone sees. Never do long work in the request thread; return 202.
- Logs: no passwords, tokens, full card numbers, secrets or full request bodies with personal data.
- New dependency only when the standard library or an existing dependency cannot reasonably do it; ask first.
- Read `references/stack-notes.md` when working in Node/TypeScript, Python, Go, Java/Spring or .NET and the framework's idiom is unclear.

## Safety gates
- Ask before: breaking an API contract, adding a dependency, changing auth/authz middleware, changing rate limits or CORS, enabling a new public endpoint, or running anything against shared or production environments.
- Never commit secrets or real customer data in fixtures; use generated data.
- Do not disable validation, auth or TLS verification to make a test pass.

## Output
```
## Change
<endpoint/job and behaviour in one or two lines>

## Contract
| Method + path | Request | Responses | Auth | Breaking? |
|---|---|---|---|---|

## Implementation notes
- Validation: ...
- Authorisation: <rule and where enforced>
- Transactions / idempotency: ...
- Outbound calls: <timeout, retry policy>

## Tests
<cases added; commands run and results>

## Consumers affected / rollout
<clients, docs, versioning, feature flag, migration order>
```

---

**Phase:** Build. **Next:** when this step is done, offer Create: `/anthropica:create-marketing-copy`, `/anthropica:create-seo-content`. Phase map: `/anthropica:workflow`.
