---
name: operate-observability
description: Design logs, metrics, OpenTelemetry traces, dashboards, SLOs with multiwindow burn-rate alerts, and telemetry cost control. Use when instrumenting a service, adding or tuning alerts or dashboards, defining SLOs or error budgets, fixing noisy alerts, high cardinality or log cost, or reviewing monitoring before launch. Not for a live incident (use `operate-incident`) or load tests (use `test-performance`).
---

# Observability and Monitoring

Instrument for the questions on-call will ask, alert only on user-visible symptoms, and keep telemetry affordable.

## Procedure
1. **Inventory.** Service type (request/response, queue worker, batch, frontend), stack and telemetry backend (Prometheus, OTel Collector, vendor), what already exists: metrics names, log format, trace coverage, alerts, dashboards, runbooks.
2. **Map critical user journeys** (login, checkout, API X). Each journey gets an SLI; everything else is supporting telemetry.
3. **Instrument the golden signals.** RED per endpoint/operation for services, USE per resource (CPU, memory, disk, pool, queue), plus saturation of limits (connection pools, thread pools, queue lag, quota).
4. **Logs and traces.** Structured logs with the required fields; OpenTelemetry tracing with context propagated across HTTP, gRPC and messaging; trace/span IDs in every log line.
5. **SLOs.** Define SLI, target, window and error budget from baseline data; read `references/slo-alerting.md` for SLI formulas, burn-rate maths and alert rules.
6. **Alerts.** Burn-rate pages and tickets on SLOs; a small set of cause-based tickets for impending saturation (disk full in < 24h, cert expiry < 14d). Each with severity, owner, routing and runbook link.
7. **Dashboards.** One service overview, then drill-downs. Annotate deploys and flag changes.
8. **Verify.** Fire each alert in a test (unit-test rules with `promtool test rules`, or inject failure in staging) and confirm it routes to the right receiver with a working runbook link.
9. **Cost review.** Cardinality, log volume and retention against budget (see Rules).

## Rules

**Logs**
- JSON, one event per line, UTC ISO-8601 timestamps. Required fields: `timestamp`, `level`, `service`, `version`, `environment`, `message`, `trace_id`, `span_id`, `request_id`; add `user_id_hash`/`tenant_id` only if policy allows.
- Levels: ERROR = needs a human (actionable, includes error type and context); WARN = degraded but handled; INFO = business events and request summaries; DEBUG off in production by default.
- Never log secrets, tokens, passwords, full card/bank numbers, auth headers, cookies, session IDs or raw request bodies. PII (email, phone, address, IP where regulated) masked or hashed; redact in the app and again in the collector.
- Log once per error at the boundary that handles it; no stack trace per retry. Sample noisy success logs; keep all errors.

**Metrics**
- Naming: `<namespace>_<what>_<unit>` with base units (`_seconds`, `_bytes`, `_total` for counters). Latency as histograms (not averages, not summaries if you need aggregation across instances).
- Cardinality: labels must have bounded values. Never label by user ID, email, request ID, full URL/path with IDs, raw error message or timestamp. Use route templates (`/users/{id}`). Budget: series per metric = product of label cardinalities; investigate any metric over ~10k series.
- Histogram buckets cover the SLO threshold exactly (if SLO is 300 ms, have a 0.3 bucket).
- Counters for rates (`rate()` over >= 4x scrape interval), gauges for levels. Do not compute percentiles by averaging percentiles.

**Traces**
- OTel SDK + auto-instrumentation first, manual spans for business operations. W3C `traceparent`/`tracestate` propagation; inject/extract on message headers for queues; link spans for fan-in/batch.
- Resource attributes: `service.name`, `service.version`, `deployment.environment`. Follow semantic conventions for HTTP/DB/messaging attributes.
- Sampling: head-based 1-10% for high volume, or tail-based in the collector to keep 100% of errors and slow traces. Never put PII or secrets in span attributes.

**SLOs and alerts**
- SLI = good events / valid events, measured as close to the user as possible (load balancer or client, not the DB). Exclude invalid traffic (4xx caused by clients) deliberately and document it.
- Target from baseline, below current performance and below dependencies' SLOs. 99.9% over 30 days = 43.2 min/month budget; 99.95% = 21.6 min; 99.99% = 4.3 min (rarely realistic).
- Page on fast burn, ticket on slow burn (multiwindow; thresholds in the reference). Page only if urgent, user-impacting and actionable now; everything else is a ticket or dashboard.
- Alert on symptoms, not causes: no paging on CPU, memory or a single pod restart. Every alert has `severity`, `team`/owner, `runbook_url`, a summary with the current value, and a dashboard link.
- Low-traffic services: burn-rate alerts misfire; use longer windows, synthetic traffic, or minimum-event-count guards (`and sum(rate(requests[1h])) > N`).
- Missing data is a failure mode: alert on `absent()`/scrape `up == 0` for critical targets and on a heartbeat (dead-man's switch) for the alerting pipeline.
- Review alerts monthly: delete or fix any that paged without action needed; target < 2 pages per on-call shift.

**Dashboards**
- Top row: SLI vs SLO, error budget remaining, request rate, error rate, p50/p95/p99 latency. Then saturation, dependencies, then per-instance drill-down. Deploy annotations on every graph.
- Units and thresholds on every panel; same time range linked; no pies for time series. Dashboards as code (JSON/Jsonnet/Terraform) in Git.

**Cost control**
- Drop unused metrics and labels at the collector (`metric_relabel_configs`, OTel `filter`/`transform` processors). Find top series with `topk(20, count by (__name__)({__name__=~".+"}))`.
- Logs: set retention by tier (hot 7-30 d, archive longer if compliance requires), sample DEBUG/INFO success logs, avoid logging large payloads. Traces: tail sampling.

## Safety gates
- Ask before changing alert routing or silencing/muting alerts in production, deleting dashboards, metrics or log indexes, lowering retention (irreversible data loss), or enabling verbose/debug logging in production (cost and PII risk).
- Never send test pages to real on-call receivers without warning the on-call first. Do not add PII to telemetry to make debugging easier.

## References
- Read `references/slo-alerting.md` when defining SLIs/SLOs, writing burn-rate alert rules, or tuning noisy alerts.
- Read `references/instrumentation.md` when adding OpenTelemetry, structured logging or metrics to code, or wiring a collector pipeline.

## Output
```
Service / journeys:  <service>, <journeys covered>
SLIs and SLOs:       <SLI formula> - <target>/<window> - budget <minutes|events>
Alerts:              <name> - page|ticket - condition - owner - runbook
Signals added:       metrics (names, labels, est. series), log fields, spans
Dashboards:          <links / files>
Verification:        rule tests / staging fire drill results
Cost impact:         series delta, log GB/day delta, retention
Gaps and risks:      <...>
```

---

**Phase:** Operate. **Next:** last phase; feed findings back into Plan: `/anthropica:plan-architecture`, `/anthropica:plan-campaign`, `/anthropica:plan-codebase-analysis`, `/anthropica:plan-test-strategy`, `/anthropica:plan-threat-model`. Phase map: `/anthropica:workflow`.
