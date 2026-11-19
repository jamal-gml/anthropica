---
name: incident-response
description: Run or support an incident from detection through mitigation, communication and a blameless postmortem. Use when production is degraded or down, an alert is firing, a security event is suspected, or a postmortem or runbook needs writing.
---

# Incident Response

## Priorities, in order
1. Stop user impact (mitigate). 2. Communicate. 3. Find the root cause. Do not debug at length while users are down if a rollback is available.

## During an incident
1. **Declare:** severity (SEV1 = full outage or data risk, SEV2 = major degradation, SEV3 = minor), incident lead, channel, start time.
2. **Assess:** what is broken, for whom, since when. Check recent deploys, config changes, dependency status pages and dashboards.
3. **Mitigate:** the safest reversible action first: roll back the last deploy, toggle the feature flag off, fail over, scale out, rate-limit. Every production action needs the incident lead's or user's go-ahead.
4. **Keep a timeline:** timestamped facts, actions taken and by whom. Mark hypotheses as hypotheses.
5. **Communicate** at a regular cadence (for example every 30 minutes for SEV1): impact, current status, next update time. No speculation in external updates.
6. **Resolve:** confirm recovery with metrics, not assumptions. Keep monitoring for regression.

For suspected security incidents, preserve evidence (logs, snapshots) before changing systems, rotate exposed credentials, and involve the security owner immediately.

## Diagnosis
- Read-only investigation first: dashboards, logs, traces, `kubectl get/describe`, cloud console views.
- Correlate the start time with deploys and changes. Compare against a healthy period.
- Change one thing at a time and observe.

## Postmortem (blameless)
```
Summary         — one paragraph: what happened and the impact
Impact          — users and regions affected, duration, data loss, SLO budget consumed
Timeline        — detection, escalation, mitigation, resolution (with UTC times)
Root cause      — technical and contributing factors (5 whys); systems, not people
What went well / what went poorly / where we got lucky
Action items    — owner, due date, ticket; prioritise prevention and faster detection
```
Publish it within the team's deadline (typically 5 working days) and track the action items to closure.

## Runbooks
For each alert: what it means, the user impact, how to verify it, step-by-step mitigation with exact commands, escalation contacts, and links to dashboards.

## Output
During an incident: status, impact, actions taken, next step. Afterwards: the postmortem or runbook document.
