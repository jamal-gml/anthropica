# Load testing

## Test types
| Type | Shape | Question | Duration |
|---|---|---|---|
| Smoke | 1-5 VUs | does the script work? | 1-2 min |
| Load | expected peak rate | do we meet SLO at peak? | 15-60 min steady |
| Stress | ramp past peak until failure | where and how does it break? | ramp in steps |
| Spike | sudden 5-10x jump | does autoscaling/queueing cope? recovery time? | short bursts |
| Soak | ~70-80% of peak | leaks, slow degradation, pool/disk growth | 4-24 h |
| Breakpoint/capacity | step ramp | max sustainable rps at SLO | until SLO breach |

## Workload design (what makes results credible)
1. **Request mix** from production: top endpoints by volume and by cost, read/write ratio, from access logs or APM (e.g. 70% browse, 20% search, 8% cart, 2% checkout).
2. **Arrival rate** from production peak x growth factor x safety margin (e.g. peak 120 rps x 1.5 growth x 1.3 = ~235 rps target).
3. **Data**: production-sized dataset, realistic key distribution (hot keys, long tail), unique users/sessions so caches are not unrealistically warm. Parameterise IDs from CSV/feeders; never hammer one ID.
4. **Think time** and session flow for user journeys; none for API-to-API tests.
5. **Environment**: production-like topology and sizes, same config/flags, autoscaling as in prod (or disabled deliberately, stated). Load generators in the same region, never CPU-saturated (check generator CPU < 70%), distributed if needed.
6. **Pass/fail criteria** defined before the run: p95/p99 latency, error rate, throughput achieved vs requested, resource ceilings.
7. **Observability** on during the test: server metrics, DB, traces, so a failure can be explained.

## Open vs closed model
- Closed (fixed VUs looping): throughput falls when the system slows, hiding problems (coordinated omission). OK for user-journey realism.
- Open (fixed arrival rate): new requests arrive regardless; queueing shows up as latency. Use for SLO and capacity tests.
- VUs needed for open model ~= rate x (p99 response time + think time), plus headroom (Little's law).

## k6 example (open model, thresholds)
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

const ids = new SharedArray('ids', () => open('./ids.csv').split('\n').filter(Boolean));

export const options = {
  scenarios: {
    browse: {
      executor: 'ramping-arrival-rate',
      startRate: 10, timeUnit: '1s',
      preAllocatedVUs: 100, maxVUs: 500,
      stages: [
        { target: 200, duration: '5m' },   // ramp
        { target: 200, duration: '20m' },  // steady state at target
        { target: 0, duration: '2m' },
      ],
    },
  },
  thresholds: {
    'http_req_failed': ['rate<0.001'],
    'http_req_duration{name:product}': ['p(95)<300', 'p(99)<800'],
    'dropped_iterations': ['count<100'],     // generator could not keep up
  },
};

export default function () {
  const id = ids[Math.floor(Math.random() * ids.length)];
  const res = http.get(`${__ENV.BASE_URL}/products/${id}`, { tags: { name: 'product' } });
  check(res, { 'status 200': (r) => r.status === 200 });
}
```
Run: `k6 run -e BASE_URL=https://staging.example test.js --out json=results.json` (or to Prometheus/InfluxDB for dashboards). `abortOnFail` on thresholds acts as a kill switch.
- Tag URLs with IDs using `name` to keep metrics low-cardinality.

## Locust example
```python
from locust import HttpUser, task, between
class Shopper(HttpUser):
    wait_time = between(1, 3)
    @task(7)
    def browse(self): self.client.get("/products", name="/products")
    @task(2)
    def search(self): self.client.get("/search?q=shoe", name="/search")
    @task(1)
    def cart(self): self.client.post("/cart", json={"sku": "A1", "qty": 1})
```
`locust -f locustfile.py --headless -u 300 -r 20 -t 30m --host https://staging.example --csv out`. Locust is closed-model; distribute workers (`--master`/`--worker`) for high rates.

## JMeter / Gatling notes
- JMeter: run non-GUI `jmeter -n -t plan.jmx -l results.jtl -e -o report/`; use Throughput Shaping Timer / Arrival thread groups for open model; disable listeners in load runs.
- Gatling: `constantUsersPerSec(50).during(10.minutes)` (open) vs `constantConcurrentUsers` (closed); assertions in the simulation fail the build.

## Running the test
1. Announce window and owners; confirm approvals (shared/prod environments, third parties stubbed).
2. Smoke run; verify scripts return 2xx and realistic payloads (a fast 4xx/5xx is not a pass).
3. Warm up 2-5 min; then steady state; record start/end timestamps for correlating dashboards.
4. Watch generator health and errors live; abort on safety thresholds.
5. Repeat key runs at least twice; results within ~10% are comparable.

## Analysing results
- Throughput achieved vs requested: if lower, the system (or generator) saturated.
- Latency vs load curve: find the knee where p99 bends upward; capacity = load just before the knee with SLO met.
- Errors by type and endpoint; first error onset time vs resource graphs.
- Identify the saturated resource (CPU, DB connections, threads, I/O, downstream rate limit) and apply the profiling workflow.
- Soak: compare memory, open connections, disk, latency at start vs end; any upward trend is a finding.

## CI performance checks
- Per PR: microbenchmarks with `benchstat`/JMH comparison against main, fail on > 10% regression beyond noise; bundle size check; Lighthouse CI.
- Nightly or pre-release: k6 load test against a dedicated perf environment with thresholds; publish trend.
- Keep perf environment stable (fixed instance types, data snapshot) so trends mean something.

## Report table
| Scenario | Target rps | Achieved rps | p50 | p95 | p99 | Errors | Limiting resource |
|---|---|---|---|---|---|---|---|
| Baseline v1.8.2 | 200 | 200 | 80 ms | 260 ms | 610 ms | 0.02% | DB CPU 65% |
| After change v1.8.3 | 200 | 200 | 55 ms | 140 ms | 300 ms | 0.01% | app CPU 55% |
| Breakpoint | ramp | 410 | - | 300 ms at 410 | - | <0.1% | DB connections |
