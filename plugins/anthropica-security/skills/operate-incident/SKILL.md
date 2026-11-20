---
name: operate-incident
description: Run a production incident through triage, mitigation (rollback, flag off, scale, failover), comms, handover and a blameless postmortem; write runbooks. Use when production is degraded or down, an alert is paging, customers report errors, a security event is suspected, or a postmortem, timeline, status update or runbook is needed. Not for designing alerts/SLOs (use `operate-observability`).
---

# Incident Response

Stop user impact first, communicate on a clock, find root cause after, learn without blame.

## Procedure
**First 15 minutes**
1. **Acknowledge** the page and open an incident record (channel/doc) with UTC start time. Say "I'm investigating" within 5 minutes.
2. **Scope the impact:** what is broken, for whom (all users, a region, a tenant, a feature), since when, and is data at risk. Use the SLO dashboard and error/latency graphs, not guesses.
3. **Set severity** from the matrix below. SEV1/SEV2: page an Incident Commander (IC) and assign roles. Do not wait for certainty to escalate; downgrading later is cheap.
4. **Check what changed:** deploys, config/flag changes, infra changes, cert/DNS/secret expiry, traffic spikes, dependency and cloud status pages, in the 2 hours before onset. `kubectl rollout history`, CI deploy log, flag audit log, IaC applies.
5. **Pick a mitigation** from the ladder below and propose it to the IC/user with expected effect and risk. Execute only after go-ahead.
6. **First comms** (internal; external if customer-facing SEV1/2) using the template in `references/comms-templates.md`.

**Stabilise and resolve**
7. **One change at a time.** Announce, apply, observe for a defined period (5-10 min), record in the timeline.
8. **Keep the timeline** as you go: `HH:MM UTC - who - fact/action/decision`. Mark hypotheses as "hypothesis". Save graph links and command outputs.
9. **Update on cadence** (SEV1 every 30 min, SEV2 every 60 min, or as promised). Each update: impact, what we did, what is next, next update time.
10. **Resolve** when SLIs are back within SLO for a sustained period (15-30 min) and confirmed by metrics, not by "the deploy finished". Keep monitoring; note any manual state to undo (feature disabled, extra capacity, silenced alerts).
11. **Handover** if the incident outlasts a shift (~4 h for IC): written state, open hypotheses, pending actions, owner of next update, explicit "you are now IC" acknowledgement.
12. **Close out:** final comms, schedule the postmortem (SEV1/SEV2 mandatory, SEV3 optional), file follow-up tickets. Read `references/postmortem.md` to write it.

## Rules

**Severity matrix** (adapt names to the team's scale; decide in < 5 min)
| Sev | Impact | Examples | Response |
|---|---|---|---|
| SEV1 | Critical: full outage, data loss/corruption, security breach, legal/financial exposure | site down, payments failing, data leak | page IC + on-call now, exec + external comms, 24/7 until mitigated |
| SEV2 | Major: core feature degraded for many users, or SLO fast-burn | checkout 20% errors, one region down | page on-call + IC, comms to stakeholders |
| SEV3 | Minor: partial degradation with workaround, small user subset | slow reports, one integration failing | business hours, ticket |
| SEV4 | Cosmetic / no user impact yet | non-prod, redundant component lost | ticket |
- When unsure between two levels, pick the higher. Any suspected data loss or security exposure is at least SEV2.

**Roles** (one person may hold several in small incidents; IC never also does deep debugging in SEV1)
| Role | Owns |
|---|---|
| Incident Commander | decisions, priorities, role assignment, severity, when to escalate or resolve |
| Ops / Tech lead | investigation and executing approved changes; pulls in SMEs |
| Comms lead | internal + status page updates on cadence, stakeholder questions |
| Scribe | timeline, decisions, links |

**Mitigation ladder** (prefer reversible, fast, well-understood)
| Situation | Action |
|---|---|
| Started after a deploy | roll back the deploy (do not fix forward under pressure unless rollback is impossible, e.g. irreversible migration) |
| Started after flag/config change | revert the flag/config |
| Overload / traffic spike | scale out, shed load, rate-limit, disable expensive non-critical features |
| Single zone/region/node bad | drain or fail over to healthy capacity |
| Bad dependency | circuit-break, serve cached/degraded, switch provider, queue writes |
| Poison message / bad job | pause consumer, move message to DLQ |
| Data corruption ongoing | stop the writer first, then assess restore; never restore over live data without a snapshot |
- Mitigation is not root cause. Once stable, keep evidence for the postmortem.
- If no hypothesis after 30 min on SEV1/2, widen: escalate to more teams, check dependencies and infra layer, compare with a healthy period/region.

**Diagnosis**
- Read-only first: dashboards, logs by trace ID, traces, `kubectl get/describe/logs`, cloud console views, DB slow-query and lock views.
- Correlate onset with changes; compare failing vs healthy instance/region/tenant. Check the obvious early: cert expiry, DNS, disk full, quota, credentials rotated, clock skew.
- Use `build-containers` troubleshooting for pod states and `operate-observability` for queries.

**Security incidents**
- Preserve evidence before changing systems (snapshot disks, export logs, record IPs/accounts). Do not tip off an attacker in public channels.
- Contain (revoke/rotate credentials, isolate hosts, block IPs), involve the security owner immediately, and keep a restricted channel. Legal/regulatory notification clocks may start at discovery (e.g. 72 h under GDPR); flag it to the IC.

## Safety gates
- Every production action (rollback, restart, scale down, failover, flag change, data fix, firewall change) needs the IC's or user's explicit go-ahead in this conversation, stating the action, target, expected effect and how to undo it.
- Ask before paging or escalating people, posting to status pages, customers or public channels, and before silencing alerts.
- Never delete data, drop tables, purge queues, force-push, or restore backups over live data during an incident without explicit approval and a fresh snapshot.
- Do not speculate on root cause or blame in external communications.

## References
- Read `references/comms-templates.md` when writing internal updates, status page posts, customer notices or a handover.
- Read `references/postmortem.md` when writing a postmortem, analysing contributing factors, or writing a runbook.

## Output
During the incident:
```
[SEV<n>] <title> - <investigating|identified|monitoring|resolved> - <HH:MM UTC>
Impact:   <who/what, scope, since when>
Done:     <actions taken, results>
Next:     <next action, owner>        Next update: <HH:MM UTC>
```
After: postmortem (template in `references/postmortem.md`) and follow-up tickets with owners and due dates.

---

**Phase:** Operate. **Next:** last phase; feed findings back into Plan: `/anthropica-security:plan-codebase-analysis`, `/anthropica-security:plan-threat-model`. Phase map: `/anthropica-security:workflow`.
