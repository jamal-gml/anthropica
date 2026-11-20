---
name: ship-ci-cd
description: Design CI/CD pipelines on any platform (stages, build once and promote, environment gates, rolling/blue-green/canary/flags, expand-contract migrations, rollback vs roll-forward, monorepos, speed, OIDC). Use when planning or reviewing a pipeline or deploy process, choosing a deployment strategy, or ordering migrations with deploys. Not for workflow YAML (use `ship-github-actions`).
---

# CI/CD Pipeline Design

Design a pipeline that gives fast feedback, deploys one verified artifact through gated environments, and can always go back.

## Procedure
1. **Inventory:** CI config files (`.github/workflows`, `.gitlab-ci.yml`, `Jenkinsfile`, `azure-pipelines.yml`, `.circleci/`, `buildspec.yml`), branch model, environments, deploy targets, artifact registries, who approves production.
2. **Map current flow** as stages with inputs/outputs: trigger -> build -> test -> scan -> package -> publish -> deploy(env) -> verify. Note where artifacts are rebuilt, where secrets enter, and manual steps.
3. **Measure:** median and p90 duration per stage, failure and flake rate, lead time commit -> production, deploy frequency, change failure rate, time to restore (DORA). Use the platform's run history.
4. **Design** the target using the Rules below. Choose a deployment strategy (table) and migration ordering per service.
5. **Plan gates** per environment (table). Decide rollback mechanism before the first deploy.
6. **Implement incrementally** behind the existing pipeline: add stages, run in parallel, switch over, then delete old paths. Platform specifics: `ship-github-actions`, `build-iac`, `build-containers`.
7. **Verify** with a dry run to a non-production environment and a rehearsed rollback.

## Stage order
| Stage | Contents | Target time | Fails fast on |
|---|---|---|---|
| Validate | format, lint, type check, config/IaC validate, commit/PR title lint | < 2 min | style, syntax |
| Build | compile, bundle, container image; record digest/checksum | < 5 min | build errors |
| Test | unit (parallel), then integration with real deps in containers | < 10 min | regressions |
| Scan | SAST, dependency audit, secret scan, image scan, IaC policy, licence | parallel with Test | known vulns above threshold |
| Package and publish | push immutable artifact once; SBOM and provenance (`secure-supply-chain`) | < 3 min | registry/auth |
| Deploy non-prod | same artifact to dev/staging; run migrations per plan | - | deploy errors |
| Verify | smoke tests, e2e on critical paths, synthetic checks | < 10 min | broken flows |
| Promote to prod | gate, deploy with strategy, automated health analysis | - | SLO breach -> rollback |

Total PR feedback target: under 10-15 minutes. Deploy pipeline commit -> production: under an hour excluding approval waits.

## Rules
- **Build once, promote the same artifact.** Identify it by digest or checksum, never by a mutable tag. Rebuilding per environment invalidates testing. Environment differences come from runtime config/secrets, not build flags.
- **Config separate from artifact:** per-environment config in a config store or deploy manifests; secrets from a secret manager at runtime.
- **Secrets:** OIDC/workload identity from CI to cloud, short-lived; no long-lived keys in CI variables. Scope prod credentials to the prod deploy job/environment only. Never print them.
- **Pipeline as code**, reviewed like application code; reusable templates for shared stages.
- **Gates are automated by default;** manual approval only for production (or regulated changes), by someone other than the author.
- **Main is always deployable:** trunk-based or short-lived branches, unfinished work behind feature flags.
- **Every deploy is recorded** (version, SHA, environment, who, when) and every production deploy has a tested rollback path.
- **Migrations decouple from code deploys** via expand/contract; never ship a destructive schema change in the same release as the code that stops using it. Read `references/database-migrations.md` for any schema change.
- **Rollback vs roll-forward:** roll back (redeploy previous artifact, flip flag, shift traffic) when users are impacted and the previous version is compatible with current data. Roll forward only when the fix is small, tested and faster than rollback, or rollback is unsafe (data already migrated irreversibly).
- **Deploy freezes and change windows** come from the team's policy; encode them as gates, not memory.
- **Monorepos:** run only affected projects (path filters or build-graph tools), but keep one always-reporting status for required checks. Read `references/pipeline-speed.md` for monorepos and slow pipelines.
- **Flaky tests** are bugs: quarantine with an owner and deadline, never auto-retry silently forever.

## Environment gates
| Environment | Trigger | Gate to enter | Exit evidence |
|---|---|---|---|
| Preview/ephemeral | PR | CI green | reviewer sign-off |
| Dev/integration | merge to main | CI green | smoke pass |
| Staging/pre-prod | auto after dev, or release candidate | smoke on dev, migrations applied | e2e + perf baseline pass |
| Production | tag or manual promote of a staged version | approval, readiness (`ship-readiness`), no freeze, error budget available | health analysis pass for bake time |

## Deployment strategy
| Strategy | Use when | Needs | Rollback |
|---|---|---|---|
| Rolling | stateless services, backward-compatible change, low risk | health checks, `maxUnavailable/maxSurge` | roll back deployment (minutes) |
| Blue-green | need instant switch and full pre-prod verification on prod infra | 2x capacity, traffic switch (LB/DNS), shared-DB compatibility | switch back (seconds) |
| Canary | high traffic, risky change, good metrics | traffic splitting, automated analysis on error/latency SLIs | shift traffic to 0% |
| Feature flag | decouple deploy from release, gradual exposure by user/segment, kill switch | flag service, flag cleanup discipline | flip flag (seconds) |
| Recreate | dev/test, or app cannot run two versions at once | downtime window | redeploy previous |

Read `references/deployment-strategies.md` when configuring canary steps, analysis thresholds, blue-green cutover or flag rollout.

## Safety gates
- Ask before: changing pipelines that deploy to production, triggering any deploy, editing environment protection, secrets, variables or approvers, running migrations against shared environments.
- Never remove or weaken a gate, required check or approval to unblock a release; escalate instead.
- Production deploys require the user's explicit approval after stating version, SHA, environment and rollback plan.

## Output
```
Current: <stages, durations, gaps, rebuilds, manual steps>
Target pipeline: <stage list with tools and parallelism>
Artifact: <type, identifier (digest/checksum), registry>
Environments and gates: <table>
Deployment strategy: <per service, with reason>
Migrations: <expand/contract plan or "none">
Rollback: <mechanism, time to restore, last rehearsed>
Secrets: <names only, delivery method (OIDC/secret store)>
Metrics to track: <lead time, deploy freq, change failure rate, MTTR, CI p90>
Changes required: <ordered list, each reversible>
```
