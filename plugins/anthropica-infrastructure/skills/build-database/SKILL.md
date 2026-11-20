---
name: build-database
description: Design schemas and make production-safe database changes (constraints, indexes, EXPLAIN, N+1, expand/contract migrations, online index builds, batched backfills, isolation, pooling, backups) with approval gates for destructive SQL. Use when writing migrations, changing a schema, adding indexes, tuning slow queries or ORM code. Not for API logic (use `build-backend`).
---

# Database Engineering

Change schemas and queries so data stays correct, production stays online, and every change can be rolled forward or back.

## Procedure
1. **Identify the setup**: engine and version, migration tool and folder, schema source of truth (migrations vs ORM models vs schema file), ORM/query builder, how migrations run in deploys (before/after app rollout, which user), connection pooler, table sizes for touched tables (ask or query catalog stats in a non-production copy).
2. **Read the current schema** for affected tables: columns, types, nullability, defaults, constraints, indexes, foreign keys, triggers, and every query that touches them (search repositories, ORM models, raw SQL, reports, other services).
3. **Design the change.** Model constraints in the database, not only in code: `NOT NULL`, `CHECK`, `UNIQUE`, foreign keys, correct types. Read `references/query-tuning.md` when adding indexes or a query is slow.
4. **Classify migration risk** for the engine and table size. Read `references/safe-migrations.md` before writing any migration on a table that has production data. Anything that rewrites the table, takes an exclusive lock for longer than a moment, or breaks the running app version needs the expand/contract sequence.
5. **Write the migration**: one logical change per migration, idempotent where the tool allows (`IF NOT EXISTS`), with lock and statement timeouts, a down/rollback step or a written roll-forward plan, and schema changes separate from data backfills.
6. **Update data-access code**: parameterised queries only; select explicit columns; eager-load or batch to avoid N+1; transactions around multi-statement writes; handle unique-violation and serialization/deadlock errors explicitly.
7. **Test**: run migrations up (and down if supported) against a fresh database and against a copy with realistic data; run repository/integration tests on the real engine (containers), not a substitute; for query changes, compare `EXPLAIN` before and after with production-like row counts.
8. **Plan the rollout**: order of migration vs app deploy, backfill duration estimate, monitoring (locks, replication lag, error rates), rollback trigger. Hand release coordination to `ship-readiness`.

## Rules
- Never edit or delete a migration that has been applied anywhere shared; add a new one. Never let an ORM auto-sync the schema in shared environments.
- Every migration must be compatible with both the currently deployed app version and the new one (deploys are not atomic). Renames and type changes therefore become: add new -> dual write/backfill -> switch reads -> stop writing old -> drop old, over separate deploys.
- Indexing:
  | Rule | Detail |
  |---|---|
  | Index for real queries | WHERE, JOIN, ORDER BY columns of frequent or slow queries; not "every column" |
  | Composite order | equality columns first, then range, then sort; leftmost-prefix rule applies |
  | Selectivity | low-cardinality columns (booleans, status) alone rarely help; use partial/filtered indexes |
  | Covering | include selected columns (`INCLUDE`) to enable index-only scans for hot reads |
  | Foreign keys | index the referencing column (Postgres does not do it automatically) |
  | Cost | each index slows writes and uses memory; drop unused and duplicate indexes (after checking usage stats) |
  | Functions | `WHERE lower(email) = ?` needs an expression index or a case-insensitive type |
- Read plans from the inside out; red flags: sequential/full scans on large tables in hot paths, estimated vs actual rows off by 10x+, nested loops over large outer sets, sorts or hashes spilling to disk, "Using filesort/temporary".
- N+1: one query per row in a loop. Fix with joins, `IN (...)` batching, ORM eager loading, or DataLoader-style batching; verify by counting queries in a test.
- Keep transactions short; no network calls or user waits inside. Lock rows in a consistent order to avoid deadlocks. Use `SELECT ... FOR UPDATE` (or `SKIP LOCKED` for queues) or optimistic version columns for read-modify-write. Retry on serialization failure and deadlock with backoff.
- Know the default isolation level (Postgres and SQL Server: Read Committed; MySQL InnoDB: Repeatable Read). Lost updates and write skew are possible below Serializable; prove the invariant with a constraint or lock.
- Pools: size per instance small (start near 2x CPU cores of the DB host divided across app instances); total connections across all instances and workers must stay below the server limit with headroom; set acquire and statement timeouts; use a pooler for serverless or many small instances. Transaction-mode poolers break session state (session `SET`, session advisory locks, `LISTEN`, some prepared-statement setups).
- Types: `timestamptz`/UTC for instants, `numeric`/decimal for money, `text`/`varchar` with CHECK for bounded strings, UUID/ULID or bigint identity keys; avoid `float` for exact values and enums that change often.
- Soft delete only if required; it needs partial unique indexes and filters everywhere. Personal data needs a deletion/retention path.
- Backups: confirm automated backups and point-in-time recovery exist, know RPO/RTO, and test restores on a schedule; a backup never restored is unverified.
- Read `references/engine-notes.md` when the engine is PostgreSQL, MySQL/MariaDB, SQL Server or SQLite and lock behaviour or syntax matters.

## Safety gates
- Never run against production or shared databases without explicit approval in the current conversation. Default to local or disposable databases.
- Destructive statements (`DROP`, `TRUNCATE`, `DELETE`/`UPDATE` without a tight `WHERE`, column drops, type changes that lose data, disabling constraints/triggers): before running, show the exact SQL, the target environment, a `SELECT count(*)` with the same predicate, the backup/snapshot status, and the rollback plan; then wait for a "yes".
- Run manual data fixes inside an explicit transaction, check affected row counts before `COMMIT`, and keep them in a reviewed script, not ad hoc.
- Never print connection strings, passwords or production rows containing personal data; use masked or sampled output.
- Ask before adding extensions, changing pool or server settings, or adding a new datastore.

## Output
```
## Change
<schema/query change and why>

## Migration plan
| Step | Statement / action | Lock / duration risk | Deploy order | Reversible? |
|---|---|---|---|---|

## Query impact
<queries affected; EXPLAIN before/after summary; indexes added/removed>

## Data safety
Backfill: <batch size, throttle, estimate> | Constraints: <added/validated> | Backup: <status>

## Tests
<migration up/down, repository tests, query count checks; commands and results>

## Rollback
<down migration or roll-forward; trigger to use it>
```

---

**Phase:** Build. **Next:** when this step is done, offer Test: `/anthropica-infrastructure:test-performance`. Phase map: `/anthropica-infrastructure:workflow`.
