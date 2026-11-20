# Blameless postmortem and runbooks

## Principles
- Blameless: people acted reasonably on the information they had; ask "why did the system allow this?", not "who did it?". Use roles, not names, in the analysis.
- Draft within 2 working days, review meeting within 5, publish to the wider org. Track action items to closure in the normal backlog.
- Required for SEV1/SEV2, any data loss, security incident, or error budget > 20% consumed by one event. Optional for near misses (often the most valuable).

## Template
```markdown
# Postmortem: <title> (<YYYY-MM-DD>, SEV<n>)
Status: draft | reviewed | published     Owner: <role/name>    Reviewers: <...>

## Summary
<3-5 sentences: what happened, impact, duration, how it was mitigated.>

## Impact
- Users/customers affected: <number or %, segments, regions>
- Duration: start HH:MM UTC -> mitigated HH:MM -> resolved HH:MM (total <n> min)
- Symptoms: <errors, latency, data issues>
- Data loss/corruption: <none | details>
- SLO / error budget: <budget consumed, % of window>
- Business: <failed transactions, support tickets, contractual/SLA exposure>

## Detection
- How detected: <alert name | customer report | employee>; time to detect: <onset -> first alert/ack>
- Did the right alert fire? Was it routed and actionable?

## Response timeline (UTC)
| Time | Event |
|---|---|
| HH:MM | <change/trigger> |
| HH:MM | <detected> |
| HH:MM | <escalation, decisions, mitigation attempts> |
| HH:MM | <mitigated> |
| HH:MM | <resolved> |
Key intervals: time to detect, time to acknowledge, time to mitigate, time to resolve.

## Root cause and contributing factors
Trigger: <the change/event that started it>
Root cause(s): <the underlying condition(s)>
Contributing factors:
- Technical: <missing test, unsafe default, capacity limit, retry storm, single point of failure>
- Process: <review gap, no canary, runbook outdated, unclear ownership>
- Detection/response: <alert missing or noisy, dashboard gap, access delays>

### 5 whys
1. Why did users see errors? <...>
2. Why did <1> happen? <...>
3. Why ...? <...>
4. Why ...? <...>
5. Why ...? <systemic cause>
(Branch when there are several causes; stop at something the team can change.)

## What went well
## What went poorly
## Where we got lucky

## Action items
| # | Action | Type | Owner | Priority | Due | Ticket |
|---|---|---|---|---|---|---|
| 1 | <specific, verifiable change> | prevent / detect / mitigate / process | <name> | P1 | <date> | <link> |

## Supporting evidence
<graphs, logs, queries, related incidents>
```

## Action item rules
- Specific and verifiable ("add canary step with automatic abort on 5xx > 1% to api pipeline"), not "be more careful" or "improve monitoring".
- Mix: at least one prevent, one detect (faster detection), one mitigate (faster recovery). Prefer automation and guardrails over training and docs.
- Every item has one owner and a due date; P1 items due within 2 sprints. Review open items in the weekly ops review.
- Recurring incident with the same contributing factor = escalate the action item priority.

## Analysis techniques
- 5 whys for linear chains; a contributing-factor tree or fishbone (people/process/technology/environment) when several causes interact.
- Ask of each factor: what made this reasonable at the time? what signal was missing? what would have contained the blast radius?
- Check: change management (was it reviewed, tested, canaried?), observability (detected by alert or customer?), resilience (timeouts, retries with backoff, bulkheads, limits), capacity, dependency failure handling, runbook accuracy.

## Metrics to report over time
Count by severity, MTTD (detect), MTTA (acknowledge), MTTM (mitigate), MTTR (resolve), % detected by alert vs customer, % caused by change, action item closure rate.

## Runbook template
```markdown
# Runbook: <AlertName>
Severity: page|ticket   Owner: <team>   Last verified: <date>
## What it means
<condition in plain words, threshold, and user impact>
## Verify
- Dashboard: <link>; query: <promql/log query>
- Confirm it is real (not a scrape/telemetry gap): <check>
## Mitigate (in order, each needs IC/owner approval in production)
1. <if recent deploy: rollback command>
2. <if overload: scale command / rate-limit flag>
3. <if dependency: failover / circuit breaker>
## Diagnose
<common causes and how to tell them apart>
## Escalate
<team / rota / when to escalate>
## Undo
<how to revert each mitigation and clean up silences/flags>
```
- Runbook steps contain exact commands with placeholders, expected output, and how to undo. Test runbooks in game days; mark the date verified.
