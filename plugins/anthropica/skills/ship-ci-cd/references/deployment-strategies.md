# Deployment strategies in detail

## Preconditions for any progressive strategy
- Two versions can run side by side: API, message and schema compatibility in both directions (N and N-1).
- Readiness and liveness checks that reflect real ability to serve (dependencies reachable), not just "process up".
- Per-version metrics: request rate, error rate, latency percentiles, saturation, plus one business metric (checkouts, sign-ins). Label by version.
- Sessions and caches tolerate mixed versions (sticky sessions only if unavoidable).

## Rolling
- Kubernetes: `maxUnavailable: 0`, `maxSurge: 25%` for zero-downtime; `minReadySeconds` 10-30 to catch crash loops; `progressDeadlineSeconds` to fail stuck rollouts.
- Rollback: `kubectl rollout undo deployment/<name>` (or redeploy previous digest via GitOps revert). Check `kubectl rollout status` exit code in the pipeline.
- PodDisruptionBudgets so node maintenance cannot take the service down mid-rollout.

## Blue-green
1. Deploy green alongside blue with production config; no user traffic.
2. Run smoke and synthetic tests against green directly.
3. Switch traffic (LB target group, service selector, DNS with low TTL - DNS is slowest and least reliable).
4. Keep blue running for the bake period (at least one peak cycle or 1-24 h by risk) for instant switch-back.
5. Decommission blue after the bake period.
- Shared database: both colours hit it, so schema changes must be expand/contract.
- Long-lived connections (websockets, gRPC streams) need connection draining.

## Canary
Typical steps (adjust to traffic volume so each step sees enough requests to be statistically meaningful, e.g. 1000+ requests):

| Step | Traffic | Bake | Promote if |
|---|---|---|---|
| 1 | 1-5% | 10-15 min | error rate <= baseline + 0.5 pp, p99 latency <= baseline x 1.2, no new error types |
| 2 | 25% | 15-30 min | same, plus business metric within normal range |
| 3 | 50% | 15-30 min | same |
| 4 | 100% | - | keep previous version warm for fast rollback until bake time ends |

- Compare canary against the baseline running at the same time, not against yesterday.
- Automate analysis and abort (Argo Rollouts, Flagger, cloud-native deploy services, service mesh); manual watching misses regressions.
- Abort = shift traffic to 0% immediately, then investigate.
- Low-traffic services: canary by time and synthetic load, or use flags instead; percentages of tiny traffic prove nothing.

## Feature flags
- Separate deploy (code present, flag off) from release (flag on). Deploy dark, then ramp: internal users -> 1% -> 10% -> 50% -> 100%.
- Flag types: release (short-lived, remove within 1-2 sprints), ops/kill switch (long-lived), experiment, permission. Track owner and expiry for each.
- Default to the safe state when the flag service is unreachable.
- Test both flag states in CI for flags that are live in production.
- Remove flags after full rollout; stale flags are dead code and combinatorial risk.

## Recreate
Only when two versions cannot coexist (exclusive locks, singleton consumers, incompatible schema without a migration window). Announce downtime; have the previous artifact ready.

## Choosing quickly
- Need seconds-level rollback and can afford 2x capacity -> blue-green.
- High traffic and good per-version metrics -> canary.
- Want business control over exposure, or partial features -> flags (combine with rolling).
- Simple internal service, compatible change -> rolling.
- Stateful singleton -> recreate with a window, or redesign.

## Post-deploy verification
- Smoke tests on critical paths immediately after each step.
- Watch golden signals for the bake time; alert thresholds tied to SLOs (`operate-observability`).
- Record deploy markers on dashboards so regressions correlate with versions.
- If anything regresses: roll back first, investigate second (`operate-incident` when users are affected).
