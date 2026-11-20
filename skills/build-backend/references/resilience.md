# Resilience: outbound calls, retries, jobs, messaging

## Outbound HTTP/RPC calls
- Reuse one client per dependency (connection pooling); never create a client per request.
- Set connect timeout (~1-3 s) and overall/read timeout from the caller's latency budget (e.g. caller SLO 500 ms -> dependency budget 200-300 ms). No timeout means a hung worker.
- Propagate deadlines/cancellation (context, `AbortSignal`, `CancellationToken`) so work stops when the caller gives up.
- Wrap each dependency in a small client module: auth, timeouts, retries, error mapping, metrics in one place.
- Map dependency failures to your own errors (502/503/504), never pass their bodies through to clients.

## Retry policy
| Retry | Do not retry |
|---|---|
| connect errors, timeouts on idempotent calls | 4xx other than 408/429 |
| 429 (honour `Retry-After`) | validation or auth errors |
| 502, 503, 504 | non-idempotent POST without an idempotency key |

- Exponential backoff with full jitter: `sleep = random(0, min(cap, base * 2^attempt))`, e.g. base 100 ms, cap 5 s, max 3 attempts.
- Retry at one layer only. Retries at client, gateway and service multiply load (3 x 3 x 3 = 27 calls).
- Keep total retry time inside the caller's deadline.
- Circuit breaker for dependencies that fail slowly: open after an error-rate threshold, fail fast, probe with half-open.
- Bulkhead: separate pools/semaphores per dependency so one slow service cannot exhaust all workers.

## Idempotency of side effects
- Use natural unique keys or a client-supplied idempotency key with a unique constraint.
- Upserts (`INSERT ... ON CONFLICT DO NOTHING/UPDATE`) instead of check-then-insert (race).
- External APIs that support idempotency keys (payment providers, email services): always send one derived from your operation ID.

## Save and publish (outbox)
Problem: writing to the DB and publishing to a broker are two systems; either can fail after the other succeeded.
1. In the same DB transaction, insert the business change and an `outbox` row (id, aggregate, type, payload, created_at).
2. A relay (poller or CDC) publishes unsent rows and marks them sent; it may publish twice, so consumers dedupe.
3. Never publish inside the transaction before commit, and never "publish then save".

## Queue consumers
- Assume at-least-once delivery and possible reordering.
- Idempotent handler: record processed message IDs (unique constraint) in the same transaction as the effect, or make the effect naturally idempotent.
- Acknowledge only after the effect is committed.
- Bounded retries with backoff; then dead-letter queue with alerting and a documented replay procedure.
- Poison messages: validate schema first; malformed messages go straight to DLQ, not retry loops.
- Visibility timeout / lease longer than the max processing time; extend for long jobs.
- Version message schemas; consumers tolerate unknown fields; producers make additive changes only.

## Background and scheduled jobs
- Request handlers enqueue; workers execute. Return 202 with a status resource for user-visible work.
- Every job: idempotent, timeout, max attempts, failed state visible in a dashboard or alert, structured logs with job ID.
- Scheduled jobs on multiple instances need a lock (DB advisory lock, lease row, or scheduler with single execution) to avoid double runs.
- Process large datasets in batches with checkpoints so a restart resumes rather than restarts.
- Store job payloads by reference (IDs), not large blobs; re-read current state when the job runs.

## Webhooks
Receiving:
- Verify the signature (HMAC over raw body, constant-time compare) and timestamp tolerance (e.g. 5 min) before parsing.
- Respond 2xx fast after durable enqueue; process asynchronously.
- Dedupe by event ID; tolerate out-of-order events (compare versions/timestamps).
Sending:
- Sign payloads, include event ID and timestamp, retry with backoff over hours, allow consumers to replay, and disable endpoints after sustained failure with notice.

## Graceful shutdown and health
- On SIGTERM: stop accepting new work, finish in-flight requests/jobs within a deadline, close pools.
- Liveness checks only the process; readiness checks critical dependencies cheaply. Do not fail liveness when the DB is down (causes restart storms).

## Rate limiting and load shedding
- Per-client/tenant limits at the edge or middleware; return 429 with `Retry-After`.
- Bound internal queues; when full, reject quickly (503) rather than queue indefinitely.

## Caching
- Cache-aside with TTL; key includes every input that changes the result (tenant, locale, version).
- Invalidate on write or keep TTL short; never cache per-user data in shared keys.
- Stampede protection: request coalescing or jittered TTLs for hot keys.
