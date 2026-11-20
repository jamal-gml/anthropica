# ADR template

Architecture Decision Records capture one hard-to-reverse decision each. Store them in the repo (commonly `docs/adr/NNNN-title.md`, numbered, never renumbered). Check for an existing ADR folder and format first and follow it.

## Rules
- One decision per ADR. Short: one page is the target.
- Immutable once accepted. To change a decision, write a new ADR that supersedes the old one and update the old one's status line only.
- Write the context so a reader in two years understands the forces without the meeting.
- List the options that were genuinely considered, including "do nothing".
- Consequences include the negative ones and the follow-up work.
- Status values: Proposed, Accepted, Rejected, Deprecated, Superseded by ADR-NNNN.

## Template
```markdown
# ADR-NNNN: <decision as a short imperative, e.g. "Use a transactional outbox for order events">

- Status: Proposed | Accepted | Superseded by ADR-NNNN
- Date: YYYY-MM-DD
- Deciders: <roles or names as the team prefers>
- Related: <issues, ADRs, design doc>

## Context
<Problem, forces and constraints. Include the numbers that matter (load, data size, latency target, deadline).
What happens if we do nothing.>

## Decision drivers
- <requirement or quality attribute, ideally measurable>
- ...

## Options considered
### Option 1: <name>
- Pros: ...
- Cons: ...
### Option 2: <name>
- Pros: ...
- Cons: ...
### Option 3: Do nothing / extend current
- Pros: ...
- Cons: ...

## Decision
<The chosen option and the main reasons, tied to the drivers.>

## Consequences
- Positive: ...
- Negative / accepted trade-offs: ...
- Follow-up work: ...
- Revisit when: <trigger, e.g. "write volume exceeds 2k/s" or "second team needs to deploy independently">
```

## Lightweight variant (for smaller decisions)
```markdown
# ADR-NNNN: <title>
Status: Accepted (YYYY-MM-DD)
In the context of <use case>, facing <concern>,
we decided for <option> and against <other options>,
to achieve <quality>, accepting <downside>.
```
