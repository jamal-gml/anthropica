# HTTP API design

Follow the existing API's conventions first. These are defaults for new APIs or gaps.

## Resources and naming
- Plural nouns, lowercase, hyphenated: `/orders`, `/orders/{orderId}/line-items`. Nest at most one level; beyond that use filters (`/line-items?orderId=...`).
- Actions that are not CRUD: sub-resource noun (`POST /orders/{id}/cancellation`) or explicit verb (`POST /orders/{id}:cancel`); pick one style per API.
- JSON field case: follow the existing API (camelCase or snake_case), never mix.
- IDs are opaque strings to clients (UUID/ULID). Do not expose sequential IDs where enumeration is a risk.

## Methods and semantics
| Method | Safe | Idempotent | Notes |
|---|---|---|---|
| GET | yes | yes | never changes state; cacheable |
| PUT | no | yes | full replace; client supplies whole representation |
| PATCH | no | not by default | partial update; JSON Merge Patch (RFC 7396) or JSON Patch (RFC 6902) |
| DELETE | no | yes | repeat returns 204 or 404 consistently |
| POST | no | no | create or command; make retry-safe with `Idempotency-Key` |

## Status codes
| Code | Use |
|---|---|
| 200 | success with body |
| 201 | created; `Location` header to new resource; return the resource |
| 202 | accepted for async work; return a status resource URL |
| 204 | success, no body |
| 304 | conditional GET, not modified (`ETag`/`If-None-Match`) |
| 400 | malformed JSON, wrong types, schema violation |
| 401 | missing/invalid/expired credentials |
| 403 | authenticated, lacks permission |
| 404 | not found, or hidden for authorisation reasons |
| 405 | method not allowed (send `Allow`) |
| 409 | conflict with current state (duplicate, invalid transition) |
| 410 | permanently removed (sunset endpoints) |
| 412 | `If-Match` precondition failed (optimistic concurrency) |
| 413 / 415 | payload too large / unsupported media type |
| 422 | well-formed but semantically invalid (if the API distinguishes it from 400) |
| 428 | precondition required (enforce `If-Match` on updates) |
| 429 | rate limited; `Retry-After` |
| 500 | unexpected server error |
| 502 / 503 / 504 | upstream bad response / unavailable or overloaded (`Retry-After`) / upstream timeout |

## Error format (RFC 9457 problem details)
Content type `application/problem+json`:
```json
{
  "type": "https://example.com/problems/validation-error",
  "title": "Request validation failed",
  "status": 422,
  "detail": "2 fields are invalid",
  "instance": "/orders",
  "traceId": "4bf92f3577b34da6a3ce929d0e0e4736",
  "errors": [
    {"pointer": "/items/0/quantity", "code": "min", "detail": "must be at least 1"},
    {"pointer": "/email", "code": "format", "detail": "must be a valid email"}
  ]
}
```
- `type` is a stable identifier clients can switch on; `title` does not vary per occurrence; `detail` is human-readable and safe to show.
- Add a machine-readable `code` per error for clients; never put stack traces, SQL or internal hostnames in any field.
- Include the trace/correlation ID so support can find the log.

## Pagination
| Style | Use | Shape |
|---|---|---|
| Cursor (keyset) | large, frequently changing, infinite scroll | `?limit=50&cursor=<opaque>` -> `{ "data": [...], "nextCursor": "..." }` |
| Offset | small, stable sets needing page numbers | `?limit=50&offset=100`; deep offsets are slow |

- Enforce a default and maximum `limit` server-side.
- Cursor encodes the last sort key plus ID (opaque, base64 or signed); sort must be deterministic.
- Avoid returning total counts on large tables unless required (expensive); offer `hasMore` instead.

## Filtering, sorting, fields
- Filters as query params with explicit names (`status=paid&createdAfter=2024-01-01T00:00:00Z`).
- Sort: `sort=-createdAt,id`; whitelist sortable fields (they must be indexed).
- Sparse fields/expansions only if clients need them; whitelist.

## Idempotency keys
For POSTs that create money movements, orders or messages:
1. Client sends `Idempotency-Key: <uuid>` per logical operation.
2. Server stores key + caller identity + request hash + response (status and body) with a TTL (often 24h), atomically with the side effect or with an in-progress lock.
3. Same key, same payload: return the stored response. Same key, different payload: 422 (or 409). Concurrent duplicate while in progress: 409, client retries later.
4. Scope keys per client/tenant.

## Concurrency control
- Return `ETag` (version or hash) on GET; require `If-Match` on PUT/PATCH/DELETE for contested resources; 412 on mismatch.
- Internally: version column with `UPDATE ... WHERE id = ? AND version = ?` and check affected rows.

## Versioning and evolution
- Additive changes are non-breaking: new optional request fields, new response fields, new endpoints, new enum values only if clients were told to tolerate unknown values.
- Breaking: removing/renaming fields, changing types or formats, making a field required, changing defaults, status codes, error `type`s or pagination.
- Version in the URL (`/v2/`) or a header/media type; follow existing practice. Run old and new in parallel; announce deprecation with `Deprecation` and `Sunset` headers and a date.
- Clients and servers: be strict in what you send, tolerant of unknown fields in what you receive.

## Long-running operations
`POST /exports` -> 202 with `Location: /exports/{id}`; `GET /exports/{id}` returns `{status: pending|running|succeeded|failed, result|error}`. Offer webhooks if clients need push.

## Security headers and limits (API)
- Request body size limits, timeouts, rate limits per client/tenant.
- CORS: explicit allowed origins, never `*` with credentials.
- `Cache-Control: no-store` for responses containing personal or sensitive data.
- Content type checks: reject unexpected `Content-Type` with 415.

## Other API styles
- gRPC: use canonical status codes (`INVALID_ARGUMENT`, `NOT_FOUND`, `ALREADY_EXISTS`, `PERMISSION_DENIED`, `UNAUTHENTICATED`, `FAILED_PRECONDITION`, `UNAVAILABLE`, `DEADLINE_EXCEEDED`); set deadlines on every call; never reuse or renumber proto field numbers (mark `reserved`).
- GraphQL: authorise in resolvers or the domain layer, not only at the gateway; batch loads with a DataLoader to avoid N+1; limit query depth and complexity; disable introspection in production if the API is not public.
