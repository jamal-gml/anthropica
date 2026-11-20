# Runbook template

One runbook per alert or operational task. Written for someone tired, at 3 a.m., who didn't build the system. For running a live incident and postmortems, use `operate-incident`.

````markdown
# Runbook: <alert name or task>
Owner: <team>   Last verified: <date, by whom>   Severity: <SEV level / priority>
Applies to: <service, environments>

## Summary
<What this alert/task means in one or two sentences, and user impact if ignored.>

## Triggers
- Alert: <name, condition, threshold, dashboard link>
- Or: <when to run this task>

## Prerequisites
- Access: <roles/groups needed, how to request>
- Tools: <CLIs with versions>, logged in to <context/cluster/account>
- Links: <dashboard>, <logs query>, <traces>, <architecture doc>

## Triage (5 minutes)
1. Confirm it's real: <command/dashboard>. Expected healthy value: <...>.
2. Scope: which environments/regions/customers? <command>
3. Recent changes: <deploy history command/link>, feature flags, config changes.
4. Decide: known cause below? -> go to Mitigation. Unknown and user impact? -> escalate now.

## Diagnosis
| Symptom | Check | Command | Healthy result |
|---|---|---|---|
| <...> | <...> | `<command>` | <...> |

## Mitigation
### Option A: <e.g. roll back last deploy>
When: <conditions>
```bash
<command>
```
Expected: <output/metric change within N minutes>
If it doesn't work: <next option / escalate>

### Option B: <e.g. scale out, fail over, disable flag>
...

## Verification
- <metric/dashboard returns to <value> for <duration>>
- <synthetic check / smoke test command>

## Rollback of the mitigation
<How to undo each option safely.>

## Escalation
| Condition | Contact | How |
|---|---|---|
| No progress in <N> min | <on-call role> | <pager/chat channel> |
| Data loss or security suspected | <role> | <...> |

## Communication
Status update template and where to post it; customer comms owner.

## Follow-up
- Ticket for root cause; update this runbook with anything that was wrong or missing.
````

## Rules
- Every command copy-pasteable, with environment/context explicit (no "the cluster"; name it or show how to select it).
- Destructive steps: warn in bold, state the blast radius, require a confirmation check before running, and give the undo.
- Mark read-only steps vs state-changing steps.
- "Last verified" date: re-test runbooks on a schedule (game days) and after architecture changes; stale runbooks get a banner.
- No secrets inline; reference where credentials are fetched.
- Link from the alert annotation directly to this runbook (see `operate-observability`).
