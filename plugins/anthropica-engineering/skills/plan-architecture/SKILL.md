---
name: plan-architecture
description: Design or evaluate software architecture (requirements and scale numbers, options with trade-offs, boundaries, data ownership, sync vs async, consistency, failure modes, migration, ADRs, phased plan). Use when planning a new application, major feature, service split, technology choice, migration or rewrite. Not for mapping existing code (use `plan-codebase-analysis`) or threats (use `plan-threat-model`).
---

# Architecture Design

Produce the smallest design that meets stated requirements, with its trade-offs and a plan that can be delivered and verified in phases.

## Procedure
1. **Requirements first.** Capture, and ask for any that are missing and decision-relevant:
   - Functional: the 3-7 core use cases, actors, and what "done" means.
   - Numbers: users, requests/s (average and peak), data volume and growth per year, payload sizes, read/write ratio.
   - Quality targets: latency (p95/p99), availability/SLO, RPO/RTO, consistency needs, retention, compliance (personal data, audit, residency).
   - Constraints: team size and skills, existing platform and languages, budget, deadline, what must not change.
   Write assumptions explicitly when the user cannot answer; mark them as assumptions.
2. **Inspect what exists** (use `plan-codebase-analysis`): current components, data stores, integration points, deployment model, pain points with evidence. A replacement must beat the current system on stated criteria, not on taste.
3. **Identify the hard parts.** The 1-3 decisions that are expensive to reverse (data model and ownership, sync vs async boundaries, tenancy, identity, consistency model). Spend effort there; keep everything else conventional.
4. **Generate 2-3 real options**, always including "extend the current system" or the simplest viable option. Compare in a table against the requirements (see Output). Score with evidence or reasoning, not adjectives.
5. **Define boundaries.** For each component: responsibility (one sentence), data it owns (single writer), interfaces it exposes, dependencies it calls, and its trust boundary. Draw the main flow (C4 context + container level is enough; Mermaid is fine).
6. **Decide communication and consistency** per interaction using the Rules tables.
7. **Walk the failure modes.** For every external call and data store: what happens when it is slow, down, returns errors, or delivers twice? Define timeouts, retries, fallbacks, idempotency and backpressure. Then cover observability, backup/restore and deployment/rollback.
8. **Check security and operations** at design level: authn/authz placement, secrets, data classification, tenant isolation, audit. Hand deep analysis to `plan-threat-model`.
9. **Plan migration** if replacing anything: strangler fig, parallel run, or branch by abstraction; data migration path; cutover and rollback per phase. Read `references/patterns.md` when choosing decomposition, integration, consistency, resilience or migration patterns.
10. **Phase the delivery.** Each phase is independently deployable, delivers observable value or risk reduction, and has acceptance criteria that can be tested. Riskiest assumption first (spike or prototype if it is unproven).
11. **Record decisions** as ADRs for each hard-to-reverse choice. Read `references/adr-template.md` when writing an ADR and the repo has no existing format.

## Rules
- Smallest design that meets requirements: one deployable and one database until a requirement (independent scaling, team autonomy, isolation, different availability targets) forces a split. Name the requirement that justifies each extra component.
- Back-of-envelope before choosing technology: requests/s x payload, storage x growth x retention, fan-out x latency. Read `references/patterns.md` when you need reference latency, throughput or availability numbers.
- Each piece of data has exactly one owning component that writes it. Others read through its API, events or a replicated read model - never shared tables across service boundaries.
- Synchronous vs asynchronous:
  | Use sync (request/response) when | Use async (queue/event) when |
  |---|---|
  | Caller needs the result to continue | Work can finish later (email, reports, indexing) |
  | Latency budget is small and the callee is reliable | Callee is slow, flaky or rate-limited |
  | Chain depth is 1-2 hops | Spikes must be absorbed (load levelling) |
  | Strong consistency is needed now | Several consumers need the same fact |
- Every sync call chain multiplies failure: availability of a chain is roughly the product of its parts (three 99.9% services in series is about 99.7%). Shorten chains or decouple them.
- Consistency: default to strong consistency inside one database transaction. Across boundaries use the outbox pattern plus idempotent consumers; use sagas with compensation for multi-step business transactions. Never use distributed two-phase commit by default.
- Async messaging is at-least-once in practice: consumers must be idempotent, handle out-of-order delivery and have a dead-letter path.
- Caches need an owner, a TTL, an invalidation rule and a stampede strategy; never the only copy of data.
- Prefer managed, boring, team-known technology. New technology needs a stated capability the current stack lacks.
- Coupling check: if two components must deploy together or change together for most features, they are one component.
- Design for the next order of magnitude, not three. State the scaling trigger that would require the next step.
- Separate fact from proposal: label checked items "Verified against <source>" and everything else "Assumption".

## Safety gates
- Planning only: do not create infrastructure, repositories, or files beyond the design doc/ADRs without asking.
- Ask before recommending deletion of data, decommissioning a system, or a big-bang cutover; present the incremental alternative.
- Flag decisions that need sign-off (cost increase, vendor lock-in, compliance, data residency) instead of deciding them.

## Output
```
## Context and requirements
Functional: ...
Scale: <users, RPS avg/peak, data size, growth>
Quality: <latency, availability, RPO/RTO, consistency>
Constraints: ...
Assumptions: ...

## Options
| Criterion | Option A (extend current) | Option B | Option C |
|---|---|---|---|
| Meets requirements | | | |
| Complexity / ops load | | | |
| Cost (rough) | | | |
| Risk / reversibility | | | |
| Team fit | | | |
Recommendation: <option> because <requirement-linked reasons>

## Design
Components: <name - responsibility - owns data - interfaces>
Flow: <diagram or numbered steps>
Communication and consistency: <per interaction>
Failure modes: | Dependency | Failure | Behaviour | Mitigation |
Security and ops: <authz, secrets, observability, backup, deploy/rollback>

## Migration (if any)
<strategy, data path, cutover, rollback>

## Phases
| Phase | Scope | Acceptance criteria | Rollback |

## Decisions (ADRs) and open questions
```

---

**Phase:** Plan. **Next:** when this step is done, offer Build: `/anthropica-engineering:build-backend`, `/anthropica-engineering:build-commits`, `/anthropica-engineering:build-database`, `/anthropica-engineering:build-frontend`, `/anthropica-engineering:build-project-setup`. Phase map: `/anthropica-engineering:workflow`.
