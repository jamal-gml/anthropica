# SLOs and burn-rate alerting

## SLI menu
| Service type | SLI | Formula (good / valid) |
|---|---|---|
| Request/response | availability | non-5xx (and non-timeout) responses / all requests excluding health checks |
| Request/response | latency | requests faster than threshold (e.g. 300 ms) / all requests |
| Queue/stream worker | freshness | messages processed within N s of enqueue / all messages |
| Batch/pipeline | completeness/timeliness | runs finishing successfully before deadline / scheduled runs |
| Storage | durability | objects readable / objects written |
| Frontend | real-user | page loads with LCP < 2.5 s / page loads |

- Count 429 from your own throttling as bad unless intended; count client 4xx as valid-but-good or exclude, and write down which.
- Latency SLI uses a threshold, not a percentile: "99% of requests < 300 ms" is computed as a ratio from histogram buckets.

## Error budget
- Budget = (1 - SLO) x valid events in window. Time view: 99.9% / 30 d = 43.2 min; 99.5% = 3.6 h; 99% = 7.2 h.
- Budget policy (agree with product owner): budget exhausted -> freeze risky releases, prioritise reliability work until it recovers; budget mostly unused -> room for faster shipping or the SLO is too loose.

## Burn rate
- Burn rate = observed error ratio / (1 - SLO). Burn rate 1 spends exactly the budget over the full window.
- Budget consumed in a window = burn rate x (alert window / SLO window).
- Error ratio threshold to alert = burn rate x (1 - SLO). For 99.9%: 14.4x -> 1.44% errors.

## Multiwindow, multi-burn-rate (30-day SLO)
| Severity | Long window | Short window | Burn rate | Budget spent when firing | Action |
|---|---|---|---|---|---|
| Page | 1 h | 5 min | 14.4 | 2% | page now |
| Page | 6 h | 30 min | 6 | 5% | page |
| Ticket | 3 d | 6 h | 1 | 10% | ticket, next business day |
- Long window gives significance; short window (1/12 of long) makes the alert reset quickly once fixed.
- Both windows must exceed the threshold (AND).
- For 28-day windows the same burn rates are fine; for very short windows recompute consumption.

## Prometheus rules (availability, SLO 99.9%)
```yaml
groups:
- name: slo-api-availability
  rules:
  - record: sli:api_errors:ratio_rate5m
    expr: |
      sum(rate(http_requests_total{job="api",code=~"5.."}[5m]))
      / sum(rate(http_requests_total{job="api"}[5m]))
  # repeat record for 30m, 1h, 6h, 3d windows: ratio_rate30m, ratio_rate1h, ratio_rate6h, ratio_rate3d

  - alert: APIErrorBudgetFastBurn
    expr: |
      (sli:api_errors:ratio_rate1h > (14.4 * 0.001) and sli:api_errors:ratio_rate5m > (14.4 * 0.001))
      or
      (sli:api_errors:ratio_rate6h > (6 * 0.001) and sli:api_errors:ratio_rate30m > (6 * 0.001))
    labels: { severity: page, team: api }
    annotations:
      summary: "API burning error budget fast ({{ $value | humanizePercentage }} errors)"
      runbook_url: "<runbook link>"
      dashboard: "<dashboard link>"

  - alert: APIErrorBudgetSlowBurn
    expr: sli:api_errors:ratio_rate3d > (1 * 0.001) and sli:api_errors:ratio_rate6h > (1 * 0.001)
    labels: { severity: ticket, team: api }
    annotations: { summary: "API error budget slow burn", runbook_url: "<runbook link>" }
```
Latency SLI ratio (bad = slower than 0.3 s):
```promql
1 - (
  sum(rate(http_request_duration_seconds_bucket{job="api",le="0.3"}[1h]))
  / sum(rate(http_request_duration_seconds_count{job="api"}[1h]))
)
```
- Tools that generate these rules from an SLO spec: Sloth, Pyrra, OpenSLO. Prefer them over hand-copied rules.
- Test with `promtool test rules tests.yaml` using synthetic series.

## Low-traffic guard
```promql
(... burn condition ...) and sum(rate(http_requests_total{job="api"}[1h])) > 1
```
With < ~1 req/s, one error can trip 14.4x. Options: guard on event count, longer windows, aggregate several services into one SLO, or synthetic probes.

## Cause-based alerts worth keeping (tickets unless imminent)
| Alert | Condition |
|---|---|
| Disk filling | `predict_linear(node_filesystem_avail_bytes[6h], 24*3600) < 0` |
| Certificate expiry | expires in < 14 d (ticket), < 3 d (page) |
| Queue backlog age | oldest message age > freshness SLO |
| Pool saturation | connection pool in use > 90% for 10 min |
| Scrape/target down | `up == 0` for 5 min on critical targets |
| Dead-man's switch | always-firing `Watchdog` alert must reach the heartbeat receiver |

## Alert quality checklist
- [ ] Fires only when a human must act now (page) or soon (ticket)
- [ ] `for:` set (typically 2-5 min for pages on non-burn-rate alerts) to avoid flapping
- [ ] Summary says what is wrong, for whom, with the current value
- [ ] `runbook_url` resolves and has verify + mitigate steps
- [ ] Routed to the owning team; grouping (`group_by: [alertname, service]`) and inhibition (cluster down inhibits per-pod) configured
- [ ] Tested: fired at least once in staging or by rule test

## Tuning noisy alerts
1. Export last 30 days of firings; for each alert, count pages, acknowledged, and pages that led to action.
2. Action rate < 50%: raise threshold/window, convert to ticket, or delete.
3. Flapping: add `for:`, hysteresis, or longer window.
4. Duplicate pages for one incident: group or inhibit.
