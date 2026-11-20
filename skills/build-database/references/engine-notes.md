# Engine notes

Version matters: confirm `SELECT version();` (or equivalent) before relying on any behaviour below.

## PostgreSQL
- Transactional DDL: most DDL can run in a transaction and roll back (exceptions: `CREATE/DROP INDEX CONCURRENTLY`, `VACUUM`, some `ALTER TYPE ... ADD VALUE` cases).
- Lock queue: an `ACCESS EXCLUSIVE` request waits behind running queries and blocks every later query. Always set `lock_timeout` for DDL.
- Lock levels to remember: `ALTER TABLE ... ADD COLUMN`, `DROP COLUMN`, `SET NOT NULL`, `ALTER TYPE` take `ACCESS EXCLUSIVE`; `VALIDATE CONSTRAINT` and `CREATE INDEX CONCURRENTLY` take `SHARE UPDATE EXCLUSIVE` (reads and writes continue).
- Foreign key columns are not indexed automatically; add indexes for joins and for `ON DELETE` performance.
- `timestamptz` stores an instant (UTC internally); `timestamp` has no zone - avoid for events.
- Identity columns (`GENERATED ALWAYS AS IDENTITY`) over `serial`; `bigint` keys for growing tables.
- MVCC: updates create new row versions; long transactions and `idle in transaction` sessions block vacuum and bloat tables. Set `idle_in_transaction_session_timeout`.
- Autovacuum: high-churn tables may need per-table tuning; monitor dead tuples and transaction ID age (wraparound).
- Useful extensions (ask before adding): `pg_stat_statements` (query stats), `pg_trgm` (LIKE/ILIKE search), `citext` (case-insensitive text), `pgcrypto` / `gen_random_uuid()` (built in from 13).
- Row-level security can enforce tenant isolation; note that table owners and superusers bypass it unless `FORCE ROW LEVEL SECURITY`.
- Isolation: Read Committed default. `SERIALIZABLE` requires retrying SQLSTATE `40001`; deadlocks are `40P01`; unique violations `23505`; FK violations `23503`.
- Poolers: PgBouncer in transaction mode does not keep session state (session `SET`, session advisory locks, `LISTEN/NOTIFY`, temp tables); check driver prepared-statement settings for compatibility with the pooler version.
- Backups: logical (`pg_dump`, per database, slow restore for large DBs) vs physical with WAL archiving for point-in-time recovery. Test `pg_restore`/PITR into a scratch instance.

## MySQL / MariaDB (InnoDB)
- DDL is not transactional: each DDL statement commits implicitly. Keep one DDL per migration.
- Online DDL: specify `ALGORITHM=INSTANT` or `ALGORITHM=INPLACE, LOCK=NONE` so the statement fails instead of silently copying the table. INSTANT add column (8.0.12+; any position from 8.0.29). Many type changes still require `COPY`.
- Metadata lock (MDL): DDL waits for open transactions on the table and then blocks new queries. Set `lock_wait_timeout` low and kill or wait for long transactions.
- Large tables: gh-ost or pt-online-schema-change for changes that cannot be INSTANT/INPLACE; check foreign key and trigger compatibility of the tool.
- Default isolation Repeatable Read with gap/next-key locks: range updates can lock more than expected and cause deadlocks; index the predicate columns of `UPDATE/DELETE`.
- Use `utf8mb4` (not `utf8`/`utf8mb3`) and an explicit collation; index key length limits apply to long `VARCHAR`.
- `DATETIME` has no zone; `TIMESTAMP` converts via session time zone and has a 2038 limit. Store UTC consistently.
- Secondary indexes include the primary key implicitly; keep primary keys short (bigint) and monotonic to avoid page splits (random UUIDv4 primary keys hurt insert performance; prefer time-ordered IDs).
- Deadlock details: `SHOW ENGINE INNODB STATUS`; error 1213 (deadlock) and 1205 (lock wait timeout) should be retried by the app.
- Replication lag affects read replicas; large single-transaction updates replicate as one big event - batch them.

## SQL Server
- DDL is transactional. Online index operations (`WITH (ONLINE = ON)`) depend on edition; resumable index operations available on newer versions.
- Adding a NOT NULL column with a default is metadata-only on Enterprise edition (2012+); otherwise it can update every row.
- Default Read Committed uses locking reads; Read Committed Snapshot Isolation (RCSI) reduces reader/writer blocking - check `is_read_committed_snapshot_on` before reasoning about blocking.
- `MERGE` has known concurrency and correctness pitfalls; prefer explicit update-then-insert with appropriate locking hints (`UPDLOCK, HOLDLOCK`).
- Use `datetime2`/`datetimeoffset`, not `datetime`; `nvarchar` for Unicode.
- Parameter sniffing can produce plans that are bad for other values; diagnose with Query Store before adding hints.
- Lock escalation to table level happens around 5,000 locks on one object: batch large updates below that.

## SQLite
- Limited `ALTER TABLE` (add column, rename column/table, drop column in 3.35+); other changes require the 12-step table rebuild (create new table, copy, drop, rename) inside a transaction with foreign keys disabled and re-checked (`PRAGMA foreign_key_check`).
- Foreign keys are off unless `PRAGMA foreign_keys = ON` per connection.
- Dynamic typing: declare `STRICT` tables (3.37+) to enforce column types.
- Single writer: use WAL mode (`PRAGMA journal_mode = WAL`) and a busy timeout for concurrent access; keep write transactions short.
- Not a stand-in for the production engine in tests when queries use engine-specific features.

## Cloud-managed and distributed SQL
- Managed services handle backups but you still own restore testing, retention settings and cross-region copies.
- Failover changes the primary endpoint or DNS: clients must reconnect and retry; keep DNS TTL caching in drivers low.
- Distributed SQL engines often run Serializable by default and require client retry on transaction conflicts; sequential primary keys create hotspots - prefer random or hashed keys. Check which DDL runs online.
