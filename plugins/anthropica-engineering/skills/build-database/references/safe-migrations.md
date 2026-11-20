# Production-safe migrations

Assume: the table has data and traffic, old and new app versions run at the same time during deploy, and any statement waiting on a lock blocks every query queued behind it.

## Always
- Set timeouts at the top of the migration so it fails fast instead of queueing traffic behind it:
  - PostgreSQL: `SET lock_timeout = '5s'; SET statement_timeout = '15min';` (per migration, tuned)
  - MySQL: `SET SESSION lock_wait_timeout = 5;`
  - SQL Server: `SET LOCK_TIMEOUT 5000;`
- Retry on lock timeout with backoff (run during low traffic for large tables).
- Check for long-running transactions before DDL (Postgres `pg_stat_activity`, MySQL `information_schema.innodb_trx`); a stuck transaction makes even fast DDL wait.
- One DDL concern per migration; separate schema migrations from data backfills.
- Estimate duration on a production-sized copy for anything touching more than ~1M rows.

## Operation playbook
| Operation | Risk | Safe approach |
|---|---|---|
| Add nullable column | low | single migration |
| Add column with constant default | Postgres 11+, MySQL 8.0.12+ (INSTANT), SQL Server 2012+ Enterprise: metadata only; older: full rewrite | check version; else add nullable, backfill, then set default |
| Add column with volatile default (`now()`, random) | table rewrite | add nullable, backfill in batches, then default for new rows |
| Add NOT NULL to existing column | full scan under exclusive lock | Postgres: `ADD CONSTRAINT c CHECK (col IS NOT NULL) NOT VALID`; `VALIDATE CONSTRAINT c`; `ALTER COLUMN SET NOT NULL` (12+ skips scan); drop check |
| Add foreign key | validates all rows under lock | Postgres: `ADD CONSTRAINT ... NOT VALID` then `VALIDATE CONSTRAINT` in a later migration |
| Add CHECK constraint | same as FK | `NOT VALID` then `VALIDATE` (Postgres); MySQL validates immediately - test duration |
| Add unique constraint | builds index under lock | create unique index concurrently/online first, then attach constraint using it |
| Create index | blocks writes for the build | Postgres `CREATE INDEX CONCURRENTLY` (not inside a transaction; on failure leaves INVALID index - drop and retry); MySQL `ALGORITHM=INPLACE, LOCK=NONE`; SQL Server `WITH (ONLINE = ON)` where edition supports it |
| Drop index | brief exclusive lock | Postgres `DROP INDEX CONCURRENTLY`; confirm no query plans rely on it |
| Rename column/table | breaks running app instantly | expand/contract: add new, dual-write, backfill, switch reads, stop old writes, drop old. Or a view/synonym bridge |
| Change column type | usually full rewrite + breaks app | new column + backfill + switch; exceptions: widening `varchar(n)` or to `text` in Postgres is metadata only |
| Drop column | running app/ORM may still select it | deploy code that no longer reads or writes it (and ORM ignores it) first; drop in a later release |
| Drop table | irreversible | rename or revoke access first, wait a release, take a backup, then drop with approval |
| Add enum value | Postgres: `ALTER TYPE ... ADD VALUE` cannot be used in the same transaction it is added (pre-12 cannot run in a transaction at all) | own migration; deploy readers that tolerate the value first |
| Large data backfill | long transactions, lock and replication lag | batched job (below), not inside the schema migration |
| Partition existing large table | rewrite | create new partitioned table, dual write or CDC, backfill, swap; plan as a project |

## Expand / contract sequence
1. **Expand**: add new column/table/index (nullable, no constraint enforcement yet). Deploy.
2. **Dual write**: app writes both old and new. Deploy.
3. **Backfill** old rows into new in batches. Verify counts/checksums.
4. **Enforce**: add NOT NULL/constraints using the NOT VALID + VALIDATE pattern.
5. **Switch reads** to new. Deploy. Monitor.
6. **Stop writing old**. Deploy.
7. **Contract**: drop old column/table in a later migration, after confirming nothing reads it.
Each step is independently deployable and reversible until step 7.

## Batched backfill
```sql
-- Postgres example; run repeatedly from a job until 0 rows updated
UPDATE orders
SET    total_minor = round(total * 100)
WHERE  id IN (
  SELECT id FROM orders
  WHERE  total_minor IS NULL
  ORDER  BY id
  LIMIT  5000
  FOR UPDATE SKIP LOCKED
);
```
- Batch by primary key ranges or `LIMIT` with an indexed predicate; 1k-10k rows per batch, each batch its own transaction.
- Sleep between batches; pause when replication lag or p99 latency rises.
- Make it resumable and idempotent (the `IS NULL` predicate is the checkpoint).
- Log progress (rows done, rows remaining); estimate total time from the first batches.
- Prefer running backfills as an application job or script, not in the migration runner (which may hold a transaction or time out deploys).

## Migration tooling gotchas
- Some tools wrap each migration in a transaction: concurrent index builds must opt out (e.g. a per-migration "no transaction" flag).
- MySQL DDL causes an implicit commit; a failed multi-statement migration leaves partial changes. Keep one DDL per migration.
- Large MySQL tables: use online schema change tools (gh-ost, pt-online-schema-change) when `INSTANT`/`INPLACE` is not available.
- ORMs that generate migrations may emit drop+add for renames or full table rebuilds (SQLite); always read the generated SQL (`migrate --sql`, `sqlmigrate`, `ef migrations script`).
- Down migrations that drop data are not rollbacks; document a roll-forward plan instead when data would be lost.
- Run migrations with a dedicated, least-privilege role; the app runtime role should not have DDL rights.

## Deploy ordering
| Change | Order |
|---|---|
| Additive schema (new column/table/index) | migrate first, then deploy app |
| Removal (drop column/table) | deploy app that stops using it first, then migrate in a later release |
| Constraint tightening | backfill and fix data, deploy app that always satisfies it, then add constraint |

## Pre-flight checklist
- [ ] Generated SQL read and reviewed
- [ ] Table sizes known; lock level and rewrite behaviour checked for this engine version
- [ ] Timeouts set; retry plan for lock timeout
- [ ] Compatible with old and new app versions
- [ ] Tested on production-sized copy (duration measured)
- [ ] Backup / PITR confirmed for destructive steps
- [ ] Monitoring during rollout: locks, replication lag, error rate, latency
- [ ] Rollback or roll-forward documented
