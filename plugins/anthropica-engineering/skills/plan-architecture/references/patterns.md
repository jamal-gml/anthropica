# Architecture patterns and reference numbers

## Back-of-envelope numbers (order of magnitude)
| Item | Rough value |
|---|---|
| Seconds per day | ~86,400 (use 100k); 1M requests/day is ~12 req/s average |
| Peak to average | plan 3-10x average unless measured |
| Main-memory read | ~100 ns |
| SSD random read | ~100 us |
| Same-region network round trip | ~0.5-1 ms |
| Cross-region round trip | ~50-150 ms |
| Single relational DB primary (well indexed, modest rows) | thousands to low tens of thousands of simple queries/s |
| Single app instance (typical web API) | hundreds to a few thousand req/s, workload dependent |
| Availability 99.9% | ~43 min downtime/month; 99.99% ~4.3 min/month |

Always replace these with measured numbers from the real system when available.

## Decomposition
| Signal | Response |
|---|---|
| One team, one product, < ~20 engineers | Modular monolith: enforce module boundaries in code (packages, lint rules), one deployable |
| Part needs very different scaling, availability or runtime | Extract that part only |
| Teams block each other on deploys | Split along team-owned business capabilities, not technical layers |
| Two services always change together | Merge them |
| Shared database between services | Assign one owner; others move to API/events/read replicas |

Boundaries follow business capabilities (bounded contexts): where the language changes ("order" means different things to billing and shipping), a boundary is likely.

## Integration styles
| Style | Use for | Watch out |
|---|---|---|
| REST/HTTP JSON | CRUD-like, public APIs, broad client support | chatty calls, versioning |
| gRPC | internal low-latency, typed contracts, streaming | browser support, debugging tooling |
| GraphQL | many clients needing different shapes of the same graph | N+1 resolvers, caching, authz per field |
| Message queue (point-to-point) | work distribution, load levelling, retries | poison messages, ordering, DLQ |
| Event stream / pub-sub | many consumers of facts, replay, audit | schema evolution, consumer lag, retention |
| Webhooks | notifying external parties | signing, retries, idempotency on receiver |
| Batch / file | large periodic exchanges | partial failure, reprocessing |

Contracts: version them, make additive changes only, use schema registry or contract tests (consumer-driven) where multiple teams depend on them.

## Consistency patterns
- **Transactional outbox**: write the business row and an outbox row in the same DB transaction; a relay publishes outbox rows to the broker. Prevents "saved but event lost" and "event sent but save rolled back".
- **Inbox / idempotent consumer**: record processed message IDs (unique constraint) in the same transaction as the side effect.
- **Saga**: sequence of local transactions with compensating actions. Orchestrated (central coordinator, easier to reason about) or choreographed (events, fewer moving parts but harder to trace). Use for multi-service business processes.
- **CQRS / read models**: separate write model from denormalised read views when read and write shapes or loads differ greatly. Adds eventual consistency; do not adopt by default.
- **Event sourcing**: store events as source of truth. Powerful for audit and temporal queries, high complexity; only with a clear need.
- Read-your-writes: route a user's reads to the primary for a short window after their write, or return the written state in the response.

## Resilience
| Concern | Default |
|---|---|
| Timeouts | every network call has one, shorter than the caller's budget |
| Retries | only idempotent operations; exponential backoff with jitter; capped attempts; retry budget |
| Circuit breaker | for dependencies that fail slowly; fail fast and degrade |
| Bulkheads | separate pools/queues so one slow dependency cannot exhaust all workers |
| Backpressure | bounded queues; reject or shed load with 429/503 rather than collapse |
| Degradation | define what still works when a dependency is down (cached data, read-only mode, queued writes) |
| Rate limits | per client/tenant at the edge |

## Data
- Choose the store by access pattern: relational by default (transactions, constraints, ad hoc queries); key-value for simple high-throughput lookups; document for self-contained aggregates; search engine for full-text/faceting (as a derived copy); time-series for metrics; object storage for blobs (store references in the DB).
- Multi-tenancy: shared schema with tenant_id (cheapest, needs strict filtering, consider row-level security), schema per tenant, database per tenant (strongest isolation, most ops). Decide early.
- Define RPO/RTO and verify backups by restoring. Retention and deletion rules for personal data are design inputs.

## Migration patterns
- **Strangler fig**: put a routing layer (proxy, gateway, or facade in code) in front of the old system; move one capability at a time; delete old paths as you go.
- **Branch by abstraction**: introduce an interface in the existing code, implement the new version behind it, switch via flag, remove the old one.
- **Parallel run / shadow traffic**: send real requests to both, compare outputs, serve from old until the diff rate is acceptable.
- **Data migration**: dual write is risky (no atomicity); prefer change data capture or outbox from the old system, backfill, verify counts/checksums, then switch reads, then writes.
- **Expand and contract** for schema and API changes: add new, migrate clients, remove old (see `build-database`).
- Each step needs a rollback path and a metric that says it worked. Avoid big-bang cutovers unless the system is small and downtime is accepted.

## Phase plan heuristics
- Phase 0: spike the riskiest unknown (throughput, third-party limits, library capability) with a measurable success criterion.
- Walking skeleton next: thin end-to-end slice through every layer, deployed to a real environment with CI, logging and a health check.
- Then vertical slices by use case, each with acceptance criteria phrased as observable behaviour and non-functional checks ("p95 < 300 ms at 200 req/s").
- Put decommissioning of the old path in the plan, with an owner.
