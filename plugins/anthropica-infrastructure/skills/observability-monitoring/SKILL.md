---
name: observability-monitoring
description: Design logging, metrics, tracing, dashboards, alerts and SLOs that catch real user impact without alert fatigue. Use when instrumenting a service, creating or tuning alerts and dashboards, defining SLOs, or reviewing monitoring coverage before a launch.
---

# Observability and Monitoring

## Signals
- **Logs:** structured (JSON), with timestamp, level, service, environment and a request or trace ID. No secrets or personal data. Log at the right level; errors must be actionable.
- **Metrics:** RED for services (Rate, Errors, Duration) and USE for resources (Utilisation, Saturation, Errors). Watch label cardinality; never use user IDs or URLs with IDs as labels.
- **Traces:** OpenTelemetry instrumentation, propagating context across service and queue boundaries.

## SLOs
1. Pick a user-facing SLI (successful requests ÷ valid requests, or the share of requests under a latency threshold).
2. Set a target from real baseline data (for example 99.9% over 28 days). Derive the error budget.
3. Alert on error-budget burn rate (multi-window: fast burn pages, slow burn files a ticket), not on raw thresholds.

## Alerts
- Every page must be urgent, actionable and user-impacting. Everything else is a ticket or a dashboard.
- Each alert has an owner, a severity, and a runbook link in its annotations.
- Alert on symptoms (errors, latency, saturation nearing limits) before causes (CPU spikes).
- Add `for:` durations to avoid flapping. Review noisy alerts and delete or tune them.

## Dashboards
- One overview per service: traffic, errors, latency percentiles (p50, p95, p99), saturation, and recent deploys as annotations.
- Put the user-impact view at the top and drill-downs below. Avoid walls of unlabelled graphs.

## Launch checklist
- [ ] Health endpoint and readiness signal
- [ ] RED metrics and an SLO defined
- [ ] Burn-rate alerts routed to the owning on-call
- [ ] Runbook for each alert
- [ ] Logs searchable with retention set
- [ ] Dashboards linked from the service README

## Output
Signals added or changed, SLIs and SLOs, alert rules with severity and routing, dashboards, and any gaps that remain.
