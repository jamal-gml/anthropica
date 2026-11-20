---
name: test-performance
description: Measure and fix latency, throughput, memory, database and frontend performance with baselines, percentiles, profilers, EXPLAIN ANALYZE, Core Web Vitals and load tests (k6, Locust, JMeter), reporting before/after numbers. Use when something is slow, times out or uses too much CPU/memory, before a launch needs capacity numbers, or to add performance budgets to CI. Not for alert design (use `operate-observability`).
---

# Performance Investigation and Load Testing

Measure first, change one thing, measure again under the same conditions, and report numbers.

## Procedure
1. **Define the question.** Symptom, affected operation/page, workload (data size, concurrency, request mix), environment, and a target: "p95 of `GET /search` < 300 ms at 200 req/s with error rate < 0.1%". No target = no finish line; agree one.
2. **Baseline.** Record p50/p95/p99, throughput, error rate and resource use (CPU, memory, GC, DB time, pool waits) from production telemetry or a repeatable benchmark. Note version, hardware, dataset and config. Repeat >= 3 runs; report the spread.
3. **Locate the bottleneck top-down.** Trace a slow request (where is time spent: app, DB, downstream, queueing?), then USE on the saturated resource, then profile the hot component. Read `references/profiling.md` for per-stack profilers and DB tools.
4. **Classify** the bottleneck: CPU-bound, memory/GC, I/O wait, lock/contention, pool exhaustion, N+1 or chatty calls, missing index, payload size, render/main-thread (frontend).
5. **Hypothesis and smallest fix** with a predicted effect ("removing N+1 cuts queries per request from 101 to 2; expect p95 -60%").
6. **Re-measure** under identical conditions (same data, warm-up, run length, machine). Keep the change only if it beats run-to-run noise and has no correctness regression (`test-validation`).
7. **Load test** when capacity or concurrency behaviour matters; read `references/load-testing.md` for design, scripts and thresholds.
8. **Guard against regressions:** add a benchmark or load-test threshold to CI and a performance budget (see Rules).
9. **Report** with the Output template.

## Rules

**Measurement**
- Percentiles, never averages alone: averages hide tail latency. Report p50, p95, p99 (and max for small samples) with sample counts. Do not average percentiles across instances; merge histograms.
- Warm up before measuring (JIT, caches, connection pools); discard the warm-up window. Run steady state >= 5-10 min for load tests.
- Benchmark noise: report median of >= 5 runs and spread; a change smaller than the noise is not a result. Use proper harnesses (JMH, BenchmarkDotNet, `go test -bench -count=10` + `benchstat`, `pytest-benchmark`, `hyperfine`) not ad-hoc timers.
- Coordinated omission: closed-loop tools under-report latency when the system stalls. Use open-model/arrival-rate executors (k6 `constant-arrival-rate`, wrk2, Gatling `constantUsersPerSec`) for latency SLO tests.
- Measure on production-like data volumes and indexes; a 1,000-row dev DB hides every query problem.

**Little's law and capacity**
- `L = lambda x W`: concurrency = throughput x time in system. 200 req/s x 0.25 s = 50 in flight. Use it to size pools, worker counts and load-test virtual users (VUs = target rps x (response time + think time)).
- Pool sizing: DB connections needed ~= rps x avg query time per request; total across replicas must stay under the DB's max connections.
- Latency rises sharply past ~70-80% utilisation of the bottleneck resource; plan capacity for peak at <= 70%.

**Common fixes by bottleneck**
| Bottleneck | Evidence | Fix first |
|---|---|---|
| N+1 queries | query count grows with result size | eager load / batch (`IN`, dataloader) |
| Missing/unused index | `Seq Scan` on large table, rows removed by filter high | composite index matching WHERE + ORDER BY; verify with EXPLAIN |
| Lock contention | waits in `pg_stat_activity`, thread dumps blocked | shorter transactions, row-level locking, queue writes |
| Pool exhaustion | time waiting for connection, timeouts at load | fix slow queries first, then size pool |
| Chatty downstream | many sequential calls in trace | batch, parallelise, cache with TTL |
| CPU hot loop / serialisation | flame graph wide frame | algorithmic fix, avoid repeated work, faster codec |
| GC / allocation | high GC time, sawtooth heap | reduce allocations, pool buffers, right-size heap |
| Large payloads | high bytes per response | pagination, compression, field selection |
| Frontend LCP/INP | Lighthouse/RUM | see `references/profiling.md` frontend section |
- Caching is a last-mile fix with invalidation cost; prefer removing the work. Database changes (indexes, migrations) follow `build-database`.

**Frontend budgets** (Core Web Vitals, field data at p75)
| Metric | Good | Poor |
|---|---|---|
| LCP | <= 2.5 s | > 4.0 s |
| INP | <= 200 ms | > 500 ms |
| CLS | <= 0.1 | > 0.25 |
| TTFB | <= 0.8 s | > 1.8 s |
- Lab (Lighthouse) for debugging, field (RUM/CrUX) for truth. Budget JS per route (e.g. < 170 KB compressed critical JS on mobile). Frontend implementation detail lives in `build-frontend`.

**Regression budgets**
- Set per endpoint/page: p95 latency, error rate, memory per request, bundle size. CI fails when a PR exceeds the budget by more than noise (e.g. > 10% p95 or > 5% bundle size).
- Track trends per release; compare against the same baseline environment.

## Safety gates
- Never load test production, shared staging, or third-party APIs (payment, email, SMS, maps) without explicit approval, a time window, owners notified, and a kill switch. Stub or sandbox third parties.
- Ask before running profilers or `EXPLAIN ANALYZE` on production (it executes the query; wrap writes in `BEGIN; ... ROLLBACK;`), adding indexes to large production tables (use `CREATE INDEX CONCURRENTLY`), or changing pool sizes, instance sizes or autoscaling limits.
- Load-test data must be synthetic or masked; no real customer PII. Clean up generated data afterwards.

## References
- Read `references/profiling.md` when profiling a specific stack, analysing SQL plans, finding memory leaks, or debugging Core Web Vitals.
- Read `references/load-testing.md` when designing or running a load, stress, soak or spike test, or adding performance checks to CI.

## Output
```
Question:     <operation>, workload, target (e.g. p95 < 300 ms @ 200 rps, errors < 0.1%)
Environment:  <env, version/commit, instance sizes, dataset size>
Baseline:     p50 / p95 / p99 / rps / errors / CPU / mem  (n runs, spread)
Bottleneck:   <component> - evidence (trace, profile, EXPLAIN excerpt)
Change:       <what changed, why>
After:        p50 / p95 / p99 / rps / errors / CPU / mem  (delta %)
Capacity:     max sustainable rps at SLO, breaking point, limiting resource
Trade-offs:   <memory, complexity, consistency, cost>
Regression guard: <benchmark/threshold/budget added>
Next:         <remaining bottlenecks in order>
```

---

**Phase:** Test. **Next:** when this step is done, offer Review: `/anthropica:review-pull-request`, `/anthropica:review-seo`, `/anthropica:review-ui-ux`. Phase map: `/anthropica:workflow`.
