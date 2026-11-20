---
name: ship-readiness
description: Run an evidence-based go/no-go review (CI on the exact SHA, security, migrations, config, flags, observability, rehearsed rollback, on-call, approvals) with blocking vs non-blocking items and a verdict. Use before tagging, promoting or deploying to production, before a launch, or when asked "are we ready to ship?". Not for cutting the version (use `ship-release`).
---

# Release Readiness

Decide GO / NO-GO from evidence, not assurances. A readiness review never authorises a deploy by itself.

## Procedure
1. **Pin the candidate:** version, commit SHA, artifact digest/checksum, target environment, planned window. Everything below is checked against this exact SHA/artifact.
2. **Scope:** `git log <prod-tag>..<sha> --format='%h %s'`, linked issues, PRs merged. Flag breaking changes, migrations, config changes, new dependencies, infra changes.
3. **Collect evidence** for every row of the checklist. Read `references/evidence.md` for commands per item. "Someone said it is fine" is not evidence; a link, run, log or output is.
4. **Classify gaps** as blocking or non-blocking using the rules below. Missing evidence for a blocking item = blocking.
5. **Rollback check:** confirm the mechanism, its time to restore, and that it was rehearsed (Rules).
6. **Operational check:** owner and on-call present for the window and after, comms drafted, dashboards and alerts ready. Read `references/launch-plan.md` for launches, high-risk changes or customer-facing announcements.
7. **Verdict:** GO, GO WITH CONDITIONS (named conditions, owners, deadline), or NO-GO (blockers). Present it; the user decides.

## Checklist
| # | Item | Evidence required | Blocking if |
|---|---|---|---|
| 1 | Scope known | change list, issues, breaking changes noted | unknown commits in range |
| 2 | CI on exact SHA | link to green run(s) for `<sha>`: build, lint, types, unit, integration | any required check failing or not run on this SHA |
| 3 | Staging verified | same artifact deployed to staging; e2e/smoke results | artifact differs from tested one, or critical flow fails |
| 4 | Performance | load/perf comparison vs baseline for changed hot paths (`test-performance`) | p95/p99 regression beyond agreed budget (default >10%) |
| 5 | Security | dependency audit, SAST, image scan, secret scan results | unresolved critical/high exploitable finding without accepted exception |
| 6 | Migrations | migration plan (expand/contract), duration on prod-sized data, lock behaviour, backup/PITR confirmed | destructive step in same release as code change; untested on realistic data; no backup |
| 7 | Config and secrets | new config keys and secrets exist in target env (names only); defaults safe | missing in prod; secret values in code |
| 8 | Feature flags | flag names, default state off in prod, kill switch tested, rollout plan | risky feature without a flag or kill switch where one was planned |
| 9 | Compatibility | API/schema/event backward compatible with N-1 and clients; deprecations announced | unannounced breaking change for external consumers |
| 10 | Observability | dashboards, alerts and logs for new paths; deploy marker; SLO/error budget status | no way to detect failure of the new path; error budget exhausted |
| 11 | Rollback | mechanism, previous artifact ID, time to restore, rehearsal date | no tested rollback, or rollback needs schema rollback |
| 12 | Dependencies and capacity | downstream teams/vendors informed; quotas, rate limits, scaling headroom | dependent change not deployed; known capacity shortfall |
| 13 | Docs and support | release notes, runbook updates, support briefed (`core-documentation`) | user-facing change with no notes (non-blocking unless contractual) |
| 14 | On-call and comms | named owner during and 2 h after; on-call aware; status/comms drafts | nobody available to respond; window conflicts with freeze |
| 15 | Approvals | required approvers per policy (change board, product, security) | required approval missing |

## Rules
- **Blocking by default:** failing required checks, untested artifact, critical/high exploitable vulnerability, missing secrets/config in prod, untested rollback, destructive migration coupled to code, no owner on call, active incident or freeze, exhausted error budget (unless the release is the fix).
- **Non-blocking (track with owner):** docs polish, nits, low/medium findings with tickets, non-critical dashboards, flaky non-required tests with an owner.
- **Exceptions** need a named approver with authority, a reason, and an expiry; record them in the output.
- **Freshness:** evidence older than the candidate SHA does not count. If the SHA changes, re-check items 2, 3, 5, 6.
- **Rollback rehearsal:** rolled back in staging (or prod for blue-green/canary) within the last release cycle, using the same mechanism; record time to restore. For flags, flip off in staging and verify.
- **Window:** avoid deploying with no one around after (late Friday, before holidays) unless urgent and staffed.
- **Partial readiness:** split scope - ship ready parts behind flags, defer the rest - rather than waiving blockers.
- **Hotfixes:** same checklist, shortened evidence: CI on SHA, targeted test, rollback, owner. Never skip rollback and approval.

## Safety gates
- This review does not deploy, tag, publish or approve anything. Hand over to `ship-release` / the pipeline only after the user's explicit go.
- Do not mark an item passed without evidence; say "not verified" instead.
- Do not waive a blocking item on the user's behalf; present it and let an authorised person decide.

## Output
```
Release readiness: <service> <version>
Candidate: SHA <sha>  artifact <digest/checksum>  target <env>  window <date/time TZ>
Verdict: GO | GO WITH CONDITIONS | NO-GO

Blocking:
  - <item #> <issue> - evidence: <link/output> - owner: <name> - fix: <action>
Conditions (if GO WITH CONDITIONS):
  - <condition> - owner - due before <time>
Non-blocking:
  - <item #> <issue> - ticket/owner

Checklist:
  | # | Item | Status (pass/fail/not verified/n.a.) | Evidence |
Rollback: <mechanism>, previous <version/digest>, TTR <min>, rehearsed <date>
Migrations: <none | plan summary>
Flags: <name: default, rollout plan>
On-call: <owner>, <on-call>, comms <channel/status page draft>
Exceptions: <item, approver, reason, expiry>
```
