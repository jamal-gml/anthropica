# Launch plan for high-risk or customer-facing releases

Use for: launches with announcements, migrations of user data, pricing/billing/auth changes, infra cutovers, anything with a large blast radius.

## Risk rating
| Factor | Low | Medium | High |
|---|---|---|---|
| Blast radius | internal or single feature | one product area | all users, money, auth, data |
| Reversibility | flag flip / instant rollback | redeploy in < 15 min | irreversible (data, external contract, mobile app) |
| Novelty | routine change | new component | new infra, new vendor, first-of-kind |
| Detectability | covered by alerts | partially | silent failure possible |

Any High -> full launch plan, rehearsal, named incident commander for the window.

## Timeline
| When | Do |
|---|---|
| T-5 to T-2 days | readiness review draft; migration rehearsal on prod-sized data; rollback rehearsal; comms drafts; freeze conflicts checked |
| T-1 day | final readiness review on the candidate SHA; go/no-go meeting; on-call briefed; stakeholders get the window |
| T-1 hour | confirm no active incidents, error budget, staff present; dashboards open; status page draft ready |
| T0 | deploy dark / canary; verify; ramp flag in steps |
| T+15 min to T+2 h | watch golden signals and business metrics; decision checkpoints at each ramp step |
| T+1 day | close out: flags cleaned up or scheduled, notes published, follow-ups ticketed |

## Go/no-go meeting
- Attendees: release owner, on-call engineer, product owner, plus security/data/support for relevant changes.
- Walk the blocking items only; non-blocking go to tickets.
- Each attendee states GO or NO-GO; any credible NO-GO pauses for resolution.
- Record decision, conditions and approver.

## Abort criteria (decide before starting)
Write them down with numbers, for example:
- Error rate > baseline + 1 pp for 5 minutes.
- p99 latency > 1.5x baseline for 10 minutes.
- Any data integrity check failing.
- Business metric (sign-ups, checkouts) below lower normal bound for 15 minutes.
- Unknown behaviour the owner cannot explain within 15 minutes.
Abort = execute rollback immediately; no debate during the window.

## Rollback rehearsal
1. Deploy candidate to staging using the production mechanism.
2. Execute the documented rollback exactly as written; time it.
3. Verify data written by the new version is readable by the old one.
4. Update the runbook with anything that differed.
5. Record date, TTR and result in the readiness output.

## Communications
| Audience | Before | During | After |
|---|---|---|---|
| Engineering/on-call | window, owner, abort criteria, rollback | progress in the release channel | outcome, follow-ups |
| Support | what changes, FAQs, known issues, escalation path | heads-up at start | confirmation, known issues |
| Customers (if impacted) | maintenance notice per contract (often 72 h+) | status page updates for degradations | release notes / announcement |
| Stakeholders | go/no-go outcome | major checkpoints | summary with metrics |

Ask before sending any external or customer-facing communication.

## Feature flag ramp
- Internal users -> 1% -> 10% -> 50% -> 100%, each step held long enough for meaningful traffic (e.g. 30-60 min or a peak period).
- Owner watches defined metrics at each step; advance only on explicit check.
- Schedule flag removal (release flags) within 1-2 sprints.
