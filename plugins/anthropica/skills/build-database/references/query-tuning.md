# Query tuning and indexing

Measure first: find the slow or frequent query (slow query log, `pg_stat_statements`, Performance Schema, Query Store, APM traces), then explain it with realistic data volumes. Plans on a near-empty dev database are meaningless.

## Getting a plan
| Engine | Estimated plan | Actual execution |
|---|---|---|
| PostgreSQL | `EXPLAIN <query>` | `EXPLAIN (ANALYZE, BUFFERS) <query>` |
| MySQL 8 | `EXPLAIN <query>` / `EXPLAIN FORMAT=TREE` | `EXPLAIN ANALYZE <query>` |
| SQL Server | estimated execution plan / `SET SHOWPLAN_XML ON` | actual execution plan / `SET STATISTICS IO, TIME ON` |
| SQLite | `EXPLAIN QUERY PLAN <query>` | - |

`EXPLAIN ANALYZE` executes the statement: for `INSERT/UPDATE/DELETE` wrap it in `BEGIN; ... ROLLBACK;` and never run it on production without approval.

## Reading a plan
- Read from the innermost/most indented node outward; cost and time accumulate upward.
- Compare estimated rows vs actual rows per node. Off by 10x or more means stale statistics or correlated predicates: run `ANALYZE`/update statistics, consider extended statistics.
- Find the node where most time or buffers go; tune only that.

| Plan node / note | Meaning | Usual fix |
|---|---|---|
| Seq Scan / full table scan / `type: ALL` on large table with selective filter | no usable index | add index matching the predicate; check functions or casts on the column |
| Index Scan with many heap fetches | index finds rows, table lookups are costly | covering index (`INCLUDE`), or accept |
| Index Only Scan with high "Heap Fetches" (Postgres) | visibility map stale | vacuum; check autovacuum |
| Bitmap Heap Scan with "Recheck" and lossy blocks | many rows matched | more selective index or `work_mem` |
| Nested Loop with large outer input | O(n x m) | index on inner join key, or hash join via better estimates |
| Sort / Hash with "external merge" / "Disk" | memory spill | index that provides the order, reduce rows earlier, raise memory per query carefully |
| `Using filesort`, `Using temporary` (MySQL) | sort/group not served by index | composite index matching ORDER BY/GROUP BY |
| Rows Removed by Filter very high | index not selective for the filter | add filtered column to composite index or partial index |

## Index design rules
- Composite order: equality predicates first, then one range predicate, then ORDER BY columns. `WHERE tenant_id = ? AND status = ? AND created_at > ? ORDER BY created_at` -> `(tenant_id, status, created_at)`.
- Leftmost prefix: `(a, b, c)` serves `a`, `a,b`, `a,b,c`; not `b` or `c` alone.
- Partial / filtered index for hot subsets: `CREATE INDEX ... ON jobs (run_at) WHERE status = 'pending'`.
- Expression index when querying a function: `ON users (lower(email))`.
- Covering: `CREATE INDEX ... ON orders (customer_id, created_at) INCLUDE (total)`.
- Multi-tenant tables: lead with `tenant_id` on most indexes.
- Remove redundant indexes: `(a)` is redundant when `(a, b)` exists (unless uniqueness or size matters).
- Check usage before dropping: Postgres `pg_stat_user_indexes.idx_scan`, MySQL `sys.schema_unused_indexes`, SQL Server `sys.dm_db_index_usage_stats` (stats reset on restart; check replicas too).
- Write-heavy tables: every index costs on each insert/update; justify each one.

## Query patterns that defeat indexes
- Functions or casts on the indexed column: `WHERE date(created_at) = ?` -> use a range `created_at >= ? AND created_at < ?`.
- Leading wildcard `LIKE '%term'` -> trigram/full-text index or search engine.
- Implicit type conversion (string column compared to number, collation mismatch).
- `OR` across different columns -> `UNION ALL` of indexed queries, or separate indexes the planner can combine.
- `NOT IN` with nullable subquery (also a correctness trap) -> `NOT EXISTS`.
- `OFFSET` deep pagination -> keyset: `WHERE (created_at, id) < (?, ?) ORDER BY created_at DESC, id DESC LIMIT 50`.
- `SELECT *` prevents index-only scans and ships unused data.
- `COUNT(*)` on huge tables for UI -> estimates, cached counters, or "has more".

## N+1 and ORM traps
- Symptom: many identical queries differing only by ID in logs/traces.
- Detect in tests: assert query count for a list endpoint (e.g. Django `assertNumQueries`, Rails Bullet, Hibernate statistics, custom counter hooks).
- Fix: eager load (`select_related`/`prefetch_related`, `include`, `JOIN FETCH`/entity graphs, `Include`), or batch with `WHERE id IN (...)`/DataLoader.
- Beware eager loading that multiplies rows (cartesian explosion with multiple collections): split queries.
- Lazy loading outside a transaction/session causes errors or extra queries; load what the view needs explicitly.
- Bulk operations: use set-based `UPDATE`/`INSERT ... SELECT`/bulk insert APIs instead of row-by-row loops.

## Concurrency patterns
- Queue tables: `SELECT ... FOR UPDATE SKIP LOCKED LIMIT n` to let workers claim rows without blocking.
- Counters and balances: `UPDATE accounts SET balance = balance - ? WHERE id = ? AND balance >= ?` and check affected rows, instead of read-then-write.
- Upsert: `INSERT ... ON CONFLICT (...) DO UPDATE` (Postgres/SQLite), `INSERT ... ON DUPLICATE KEY UPDATE` (MySQL), `MERGE` carefully (SQL Server `MERGE` has concurrency caveats; use `HOLDLOCK` or update-then-insert).
- Optimistic locking: `version` column checked in the `WHERE` of updates.
- Advisory locks for "run once" jobs (Postgres `pg_try_advisory_lock`, MySQL `GET_LOCK`).

## Lock and health diagnostics (read-only)
- Postgres: `pg_stat_activity` (long transactions, `idle in transaction`), `pg_locks` joined to activity for blockers, `pg_stat_user_tables` (seq scans, dead tuples, last autovacuum).
- MySQL: `SHOW ENGINE INNODB STATUS` (latest deadlock), `performance_schema.data_locks`, `sys.innodb_lock_waits`.
- SQL Server: `sys.dm_exec_requests`, `sys.dm_tran_locks`, Query Store regressed queries.

## Caching and read scaling
- Fix the query and index before adding a cache.
- Read replicas lag: do not read-your-own-write from a replica; route those reads to the primary.
- Materialised views/summary tables for heavy aggregates, with a defined refresh strategy and staleness tolerance.
