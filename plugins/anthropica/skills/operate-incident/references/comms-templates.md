# Incident communication templates

Rules: facts only, plain language, UTC times, no blame, no speculation on cause externally, always state the next update time and keep it. Get approval before posting anything outward-facing.

## Declaration (internal channel)
```
[INCIDENT] Declaring SEV<n>: <short title>
Impact: <who/what is affected, scope, since HH:MM UTC>
IC: <name>  Ops: <name>  Comms: <name>  Scribe: <name>
Channel/doc: <link>   Bridge: <link>
Current action: <what we are doing>
Next update: HH:MM UTC
```

## Internal update (every 30 min SEV1, 60 min SEV2)
```
SEV<n> <title> - update HH:MM UTC - status: investigating | identified | mitigating | monitoring
Impact now: <same | reduced to X | resolved>
Since last update: <actions and results>
Current hypothesis: <...> (unconfirmed)
Next steps: <action - owner>
Help needed: <teams/skills>
Next update: HH:MM UTC
```

## Status page / customer-facing
Investigating:
```
We are investigating reports of <symptom, e.g. errors when signing in> affecting <some|all> users <in region>.
Started around HH:MM UTC. We will provide an update by HH:MM UTC.
```
Identified:
```
We have identified the cause of <symptom> and are applying a fix. Some users may continue to see <symptom>.
Next update by HH:MM UTC.
```
Monitoring:
```
A fix has been applied and <service> is recovering. We are monitoring to confirm full recovery.
Next update by HH:MM UTC.
```
Resolved:
```
This incident has been resolved. Between HH:MM and HH:MM UTC, <who> experienced <symptom>.
<If applicable: no data was lost / actions customers need to take.>
We will share a summary of what happened and what we are doing to prevent it.
```
- Do not name vendors, people or internal systems. Do not promise times you cannot meet ("within the hour").
- Security/data incidents: external wording reviewed by security and legal before posting.

## Stakeholder / executive summary
```
What: <one sentence>
Impact: <users, revenue/transactions affected if known, SLO budget used>
Status: <mitigated at HH:MM / ongoing>, ETA: <if known, else "unknown">
Risk: <data, security, contractual>
Ask: <decisions needed from you>
```

## Handover
```
Handover HH:MM UTC: <from> -> <to> (IC|Ops|Comms)
State: <impact now, status>
Mitigations in place (to undo later): <flags off, extra capacity, silences with expiry, manual data fixes>
Ruled out: <hypotheses and evidence>
Open hypotheses: <...>
In progress: <action - owner - expected result>
Pending approvals/decisions: <...>
Next comms due: HH:MM UTC
Ack: "<to> is now IC"
```

## Timeline entry format
```
HH:MM UTC - <who> - <fact | action | decision | hypothesis> - <evidence link>
```
Examples:
```
14:02 - alerting - APIErrorBudgetFastBurn fired (5xx 6.1%)
14:05 - A. - ack, investigating
14:09 - A. - deploy api v1.8.3 at 13:58 correlates with onset (hypothesis)
14:12 - IC - decision: roll back to v1.8.2
14:15 - A. - rollback complete; 5xx 0.4% and falling
14:35 - IC - resolved; monitoring 30 min
```
