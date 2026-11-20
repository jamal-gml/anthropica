# Profiling by stack

Profile under realistic load; a profile of an idle process shows startup code. Prefer sampling profilers (low overhead) in shared environments. Read flame graphs by width (time), not height; look for wide plateaus.

## CPU and wall-clock profilers
| Stack | CPU / wall | Memory / allocations | Notes |
|---|---|---|---|
| Linux (any) | `perf record -F 99 -g -p <pid>`; `perf top` | `bpftrace`, `memleak` (bcc) | needs frame pointers or debug symbols |
| Java/JVM | async-profiler `asprof -d 30 -e cpu -f out.html <pid>`; JFR `jcmd <pid> JFR.start duration=60s filename=rec.jfr` | `asprof -e alloc`; heap dump `jcmd <pid> GC.heap_dump h.hprof` -> Eclipse MAT | thread dump `jcmd <pid> Thread.print` for locks |
| Go | `go tool pprof http://host/debug/pprof/profile?seconds=30` | `.../debug/pprof/heap`, `allocs`; `-base` to diff | `goroutine` and `mutex`/`block` profiles for contention |
| Python | `py-spy record -o out.svg --pid <pid>`; `py-spy top`; `cProfile` + snakeviz locally | `memray run`, `tracemalloc` snapshots diff | `py-spy dump` shows stuck threads; asyncio: check blocking calls in event loop |
| Node.js | `node --cpu-prof app.js`; `0x`; clinic.js `clinic flame` | `--heapsnapshot-signal=SIGUSR2`, Chrome DevTools heap compare | `clinic doctor` for event-loop delay; `--inspect` |
| .NET | `dotnet-trace collect -p <pid>`; `dotnet-counters monitor` | `dotnet-gcdump`, `dotnet-dump` | |
| Ruby | `rbspy record`, `stackprof` | `memory_profiler`, `derailed_benchmarks` | |
| PHP | `xdebug` profiler (dev), `spx`, `blackfire` | | |
| Continuous | Pyroscope / Parca / vendor continuous profiling | | compare before/after deploy |

## Memory leak workflow
1. Confirm growth: RSS/heap trending up under steady load over hours (not just warm-up or cache fill to a cap).
2. Take two heap snapshots/dumps separated by load; diff them (retained size by type/allocation site).
3. Usual suspects: unbounded caches/maps, listeners/subscriptions not removed, global collections, closures holding requests, connection/stream not closed, thread-locals in pools.
4. Fix, then soak test (see `load-testing.md`) to prove flat memory.
Container OOMKilled with a healthy heap: native memory, thread stacks, or heap limit not aligned to container limit.

## Database
**PostgreSQL**
```sql
EXPLAIN (ANALYZE, BUFFERS, VERBOSE, SETTINGS) <query>;   -- executes the query
-- writes: BEGIN; EXPLAIN ANALYZE ...; ROLLBACK;
SELECT query, calls, mean_exec_time, total_exec_time, rows
FROM pg_stat_statements ORDER BY total_exec_time DESC LIMIT 20;
SELECT * FROM pg_stat_activity WHERE state <> 'idle';          -- running / waiting
SELECT relname, seq_scan, idx_scan, n_live_tup FROM pg_stat_user_tables ORDER BY seq_scan DESC;
```
Plan reading:
| Signal | Meaning |
|---|---|
| `Seq Scan` on big table with selective filter | missing or unusable index (function on column, type mismatch, leading wildcard `LIKE '%x'`) |
| estimated rows vs actual off by 10x+ | stale stats (`ANALYZE`), correlated columns (extended statistics) |
| `Rows Removed by Filter` large | index does not cover the predicate |
| `Sort Method: external merge` | `work_mem` too small or missing index for ORDER BY |
| `Nested Loop` with large outer side | bad estimate; check join columns indexed |
| high `Buffers: shared read` | cold cache / I/O bound |
- Tools: `auto_explain` for slow queries in logs, explain.dalibo.com / explain.depesz.com to visualise.

**MySQL/MariaDB**: `EXPLAIN ANALYZE` (8.0.18+), `type: ALL` = full scan, check `key`, `rows`, `Extra: Using filesort/temporary`; slow query log with `long_query_time`; `performance_schema` digest tables; `pt-query-digest`.

**SQL Server**: actual execution plan, Query Store top resource consumers, `sys.dm_exec_query_stats`, missing index DMVs (validate before creating).

**MongoDB**: `.explain("executionStats")` - `COLLSCAN`, `totalDocsExamined` vs `nReturned`; profiler `db.setProfilingLevel(1, {slowms: 100})`.

**Index rules**
- Composite index order: equality columns first, then range, then sort. Cover the query (`INCLUDE`) for hot reads.
- Every index slows writes and uses memory; drop unused ones (`idx_scan = 0` over a full business cycle).
- Production: `CREATE INDEX CONCURRENTLY` (Postgres), online DDL (MySQL `ALGORITHM=INPLACE, LOCK=NONE`); ask first.

**N+1 detection**: count queries per request (ORM logging, APM span count, `django-debug-toolbar`, `bullet` gem, Hibernate statistics). Fix with eager loading (`select_related`/`prefetch_related`, `includes`, `JOIN FETCH`, `Include()`), batching, or dataloaders for GraphQL.

## HTTP / service level
- `curl -w '@fmt.txt' -o /dev/null -s <url>` with `time_namelookup time_connect time_appconnect time_starttransfer time_total` to split DNS/TCP/TLS/server/transfer.
- `hey`/`oha`/`wrk` for quick single-endpoint checks (not a substitute for a designed load test).
- Distributed trace of a slow request: find the widest span and sequential gaps (time not covered by child spans = in-process work or queueing).

## Frontend (Core Web Vitals)
| Metric | Usual causes | Fixes |
|---|---|---|
| LCP | slow TTFB, render-blocking CSS/JS, late-discovered hero image, client-side rendering | SSR/SSG, `<link rel=preload>`/`fetchpriority="high"` on LCP image, inline critical CSS, CDN, modern image formats sized correctly |
| INP | long tasks (> 50 ms) on main thread, heavy event handlers, hydration | split work (`scheduler.yield`/`setTimeout`), defer non-critical JS, web workers, reduce re-renders |
| CLS | images/ads/embeds without dimensions, late fonts, injected banners | `width`/`height` or `aspect-ratio`, reserve space, `font-display: optional`/size-adjust |
| TTFB | server time, no caching, redirects | cache HTML at edge, fix backend, remove redirect chains |
Tools: Lighthouse / PageSpeed Insights (lab + CrUX field), Chrome DevTools Performance panel (long tasks, layout shifts), WebPageTest (filmstrip, waterfall, throttled mobile), `web-vitals` library for RUM, bundle analysers (`source-map-explorer`, `webpack-bundle-analyzer`, `vite-bundle-visualizer`).
Always test throttled mid-tier mobile (4x CPU slowdown, slow 4G); desktop numbers flatter.

## Lighthouse CI budget example
```json
{ "ci": { "assert": { "assertions": {
  "largest-contentful-paint": ["error", { "maxNumericValue": 2500 }],
  "cumulative-layout-shift": ["error", { "maxNumericValue": 0.1 }],
  "total-blocking-time": ["warn", { "maxNumericValue": 200 }],
  "resource-summary:script:size": ["error", { "maxNumericValue": 170000 }]
} } } }
```
