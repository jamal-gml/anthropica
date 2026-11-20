# Database migrations in the pipeline

Goal: every deploy and every rollback works while two app versions share one database. Schema authoring detail lives in `build-database`; this file covers ordering with deploys.

## Expand / contract (parallel change)
| Release | Schema step | App step | Rollback safe? |
|---|---|---|---|
| 1. Expand | add new column/table/index (nullable or defaulted), no drops or renames | old code still works unchanged | yes |
| 2. Migrate | backfill in batches; dual-write old and new | new code writes both, reads old (or new with fallback) | yes |
| 3. Switch | - | read new; still write both if rollback to 2 must stay possible | yes, to 2 |
| 4. Contract | drop old column/table after N-1 no longer uses it | stop writing old | only by restore; do it after bake time |

Each row is a separate deploy. Never combine expand and contract in one release.

## Ordering rules
- **Additive migrations run before the app deploy** (pre-deploy job or init step that blocks rollout on failure).
- **Destructive migrations run after** the code that stops using the object is fully deployed and past the bake time.
- **One migration runner at a time:** use the tool's lock (Flyway, Liquibase, Alembic, Prisma, Rails, EF, golang-migrate) or a pipeline concurrency group. Never run migrations from every app replica at startup without a lock.
- **Migrations are versioned, immutable once applied** in any shared environment. Fix with a new migration.
- **Same migration path for every environment;** never hand-apply SQL in production.

## Dangerous operations and safe alternatives
| Operation | Risk | Safer approach |
|---|---|---|
| Rename column/table | breaks N-1 code | add new, dual-write, backfill, switch, drop old |
| Drop column | N-1 code fails, irreversible | stop using, deploy, bake, then drop |
| Add NOT NULL column without default | fails on existing rows / rewrites table | add nullable, backfill, then add constraint (validate separately where supported) |
| Change column type | table rewrite, lock | new column + backfill + switch |
| Add index on large table | long write lock | concurrent/online index creation (`CREATE INDEX CONCURRENTLY` in PostgreSQL, online DDL in MySQL) |
| Add foreign key / check constraint | full scan under lock | add as not-validated, validate later |
| Large backfill in one transaction | lock contention, replication lag, long rollback | batches of 1k-10k rows with pauses; resumable job |
| Enum value removal | old code writes it | stop writing, migrate data, then remove |

Set `lock_timeout` / `statement_timeout` (or equivalent) on migration sessions so a blocked DDL fails fast instead of queueing all traffic behind it.

## Rollback
- Prefer roll-forward migrations over down migrations in production; down migrations that drop data are restores in disguise.
- App rollback must not require schema rollback: that is what expand/contract guarantees.
- Before contract steps and risky backfills: verified backup or point-in-time recovery window, and a tested restore time.

## Pipeline checks
- Lint migrations for destructive statements (`squawk`, `atlas migrate lint`, custom grep for `DROP|RENAME|ALTER .* TYPE`) and require explicit approval when found.
- Apply migrations to a production-like copy (size and data shape) in staging and record duration and locks.
- Run the N-1 app version's test suite against the new schema for expand steps.
- Record which migration version is applied per environment alongside the app version.

## Data stores without schemas
Document stores and event streams need the same discipline: readers tolerate unknown and missing fields, version events, never repurpose a field's meaning.
