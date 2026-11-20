# Gathering evidence per checklist item

Run read-only commands. Record outputs or links; never paste secrets or tokens.

## 1. Scope
```bash
git fetch --tags origin
git describe --tags --abbrev=0 <sha>^            # previous release
git log <prod-tag>..<sha> --format='%h %an %s'   # every commit in the release
git diff --stat <prod-tag>..<sha>
git log <prod-tag>..<sha> --grep='BREAKING CHANGE' --format='%h %s'; git log <prod-tag>..<sha> --format='%s' | grep -E '^[a-z]+(\(.+\))?!:'
gh pr list --state merged --search "merged:>=<date-of-prod-tag>" --json number,title,url
```
Flag paths that raise risk: `migrations/`, `db/`, infra code, auth, payment, config files, dependency manifests.

## 2. CI on the exact SHA
```bash
gh run list --commit <sha> --json name,conclusion,url
gh api repos/{owner}/{repo}/commits/<sha>/check-runs -q '.check_runs[] | {name, conclusion}'
gh api repos/{owner}/{repo}/commits/<sha>/status -q '.state'
```
All required checks `success`. "Neutral"/"skipped" on a required check needs an explanation.

## 3. Staging is running the same artifact
- Compare digests: registry digest for the release tag vs what staging reports (`kubectl get deploy <name> -o jsonpath='{..image}'`, platform deployment record, `/version` endpoint).
- Smoke/e2e run link for that deployment.

## 4. Performance
- Load test result vs previous release baseline on the same environment and data volume (`test-performance`).
- Default budget: p95/p99 latency within +10%, error rate unchanged, resource use within +20%; use the team's SLOs if defined.

## 5. Security
- Dependency audit (`npm audit --omit=dev`, `pip-audit`, `govulncheck ./...`, `osv-scanner`), image scan (`trivy image <ref>`, `grype`), code scanning alerts: `gh api repos/{owner}/{repo}/code-scanning/alerts?state=open&ref=<branch>`, secret scanning alerts.
- Classify: critical/high + reachable/exploitable = blocking unless exception. Depth: `secure-code-review`, `secure-supply-chain`.

## 6. Migrations
- List new migration files in the range; read each for destructive ops (DROP, RENAME, type change, NOT NULL without default, non-concurrent index).
- Staging apply duration and lock waits on prod-sized data.
- Backup/PITR: latest successful backup time and restore test date.
- Ordering plan per `ship-ci-cd` (expand before deploy, contract later).

## 7. Config and secrets
- Diff config schema/defaults between `<prod-tag>` and `<sha>`; list new keys.
- Confirm presence (not values) in the target environment: secret manager listing by name, CI environment secret names (`gh secret list --env production`), config store keys.

## 8. Feature flags
- New flags in the diff (grep the flag SDK call names); default per environment in the flag service; who can flip them; kill-switch test in staging.

## 9. Compatibility
- API schema diff (OpenAPI/GraphQL/protobuf diff tools such as `oasdiff breaking`, `buf breaking`, GraphQL inspector).
- Event/message schema registry compatibility check.
- Client minimum versions (mobile apps cannot be rolled back).

## 10. Observability
- Dashboard and alert links covering new endpoints/jobs; test alert fired in staging if new.
- Current SLO status and remaining error budget; open incidents (`operate-incident`, `operate-observability`).

## 11. Rollback
- Previous production artifact ID/digest and the exact rollback command or pipeline job.
- Rehearsal record: date, environment, time to restore, issues found.

## 12. Dependencies and capacity
- Dependent services' versions deployed (contract first, consumer later).
- Quotas and rate limits for new third-party calls; autoscaling limits vs expected load.

## 13-15. Docs, on-call, approvals
- Release notes draft link; runbook PR merged.
- On-call schedule for the window (paging tool); owner confirmed by name.
- Required approvals recorded (environment approval, change ticket, sign-offs).
