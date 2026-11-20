---
name: plan-threat-model
description: Threat-model a system or feature - scope, data-flow diagram with trust boundaries, STRIDE per element, abuse cases, LINDDUN privacy, likelihood x impact rating, mitigations, residual risk - as a threat table. Use when designing auth, payments, personal data, uploads, webhooks, multi-tenancy, LLM agents, integrations or new infrastructure, or for a security design review. Not for reviewing code (use `secure-code-review`).
---

# Threat Modeling

Answer four questions: what are we building, what can go wrong, what are we doing about it, did we do a good job. Output a threat table that drives design changes, tickets and tests.

## Procedure
1. **Scope.** State the change in one paragraph, what is in and out of scope, assumptions, and the security objectives (which data and actions must stay confidential, intact, available, attributable). For an incremental change, model only the new or changed flows plus every boundary they cross; link the prior model.
2. **Ground it in evidence.** Read the code, IaC, API specs, auth middleware and existing diagrams (`plan-codebase-analysis` for unfamiliar repos). Label anything you could not verify as `ASSUMED`.
3. **Inventory.** List external entities (users, admins, partners, third-party APIs, CI), processes (services, workers, functions), data stores (DBs, caches, buckets, queues, logs) and data flows. Classify data: secrets/credentials, payment, personal (special category if health, biometric, etc.), internal, public.
4. **Draw the DFD** (Mermaid `flowchart` is fine). Number every element and flow (`P1`, `DS1`, `F3`). Draw a trust boundary wherever privilege, ownership, network zone or tenant changes: internet to edge, edge to service, service to datastore, tenant to tenant, service to third party, CI to production, user code to sandbox, LLM output to tool execution.
5. **Enumerate threats with STRIDE per element.** Apply only the categories that fit each element type:

   | Element | S | T | R | I | D | E |
   |---|---|---|---|---|---|---|
   | External entity | x | | x | | | |
   | Process | x | x | x | x | x | x |
   | Data store | | x | (logs) | x | x | |
   | Data flow | | x | | x | x | |

   Give boundary-crossing flows the most attention. Write each threat as "actor can do X to asset Y via Z, causing W" - concrete, not "attacker may tamper with data". Read `references/stride-prompts.md` for per-component question lists (web, API, auth, storage, queues, CI/CD, cloud, mobile, LLM/agents).
6. **Abuse cases.** For each key user story, write the misuse story: enumeration, replay, race/double-spend, coupon or quota abuse, account takeover via recovery, privilege via invitations, scraping, cost exhaustion (paid APIs, LLM tokens), confused deputy. Business-logic threats are where scanners find nothing.
7. **Privacy.** If personal data is processed, run LINDDUN over the same DFD. Read `references/linddun-privacy.md` when personal data, tracking, analytics or cross-border transfer is in scope.
8. **Rate** each threat with the method below, taking existing controls into account (cite the file or config that proves the control exists; otherwise treat it as absent).
9. **Decide** per threat: mitigate, eliminate (drop the feature or the data), transfer (contract, insurance, managed service), or accept (named owner, reason, review date). Map every mitigation to the threat IDs it addresses and to a verifiable control: test case, config check, policy-as-code rule or alert.
10. **Residual risk.** Re-rate after planned mitigations. Anything still High or Critical needs an explicit accept decision from the owner, not silence.
11. **Hand off.** Turn mitigations into tickets; feed security test cases to `plan-test-strategy`; note which findings `secure-code-review` should verify once code exists.

## Rules
- **Rating.** Likelihood 1-4 (Rare: needs insider plus chained bugs; Unlikely: skilled attacker, non-default config; Likely: authenticated user with public tooling; Almost certain: anonymous, scripted, already seen in the wild). Impact 1-4 (Minor: no sensitive data, quick recovery; Moderate: limited data of one tenant, degraded service; Major: many users' personal data, integrity loss, hours of outage; Severe: secrets, payments, cross-tenant, regulatory breach, RCE). Score = L x I: 1-3 Low, 4-6 Medium, 8-9 High, 12-16 Critical. Use DREAD or CVSS only if the team already does; never mix scales in one table.
- Pre-auth and internet-facing raise likelihood by one step; a compensating control you verified lowers it by one. Never lower impact for "nobody would do that".
- Default trust assumption: anything crossing a boundary inbound is hostile until validated, including other internal services, webhooks, files, headers, and LLM output.
- Prefer eliminating data or capability over guarding it (do not store card numbers, do not give the agent a shell).
- Authorisation threats (IDOR, cross-tenant read, admin endpoints) almost always rate at least High when personal or financial data is involved.
- Every accepted risk has an owner and a review date; every mitigation names a control a reviewer can check. "Use best practices" is not a mitigation.
- Revisit the model when: a new trust boundary or external integration appears, auth or tenancy changes, a new data class is stored, the deployment topology changes, after a security incident or pentest finding in the area, or at least yearly for critical systems.
- Keep it proportional: a small feature gets a one-page model (DFD, 5-15 threats); a new platform gets the full treatment.

### Common failure modes
- Threat lists copied from a generic checklist with no link to this system's elements.
- Missing boundaries for CI/CD, admin tooling, support access, backups and observability pipelines (logs often hold the most sensitive data).
- Rating against intended controls instead of controls that exist today.
- Treating TLS as the answer to tampering inside the trust zone, or "it's internal" as authentication.

## Safety gates
- Threat modeling is analysis only: do not probe, scan or exploit any system to "confirm" a threat without written authorisation for that target.
- Never paste secrets, real customer data or production connection strings into the model; describe their class instead.
- The model itself is sensitive (it lists weaknesses): keep it in the team's private repo or doc space and ask before publishing or sharing it externally.

## Output
```markdown
## Threat model: <system/feature>  (version, date, author, status)
**Scope:** in / out / assumptions (ASSUMED items flagged)
**Security objectives:** ...
**Data classes:** ...
**DFD:** <mermaid diagram with numbered elements and trust boundaries>

| ID | Element/flow | STRIDE/LINDDUN | Threat (actor, action, asset, effect) | Existing controls (evidence) | L | I | Risk | Decision | Mitigation -> control/test | Owner | Residual |
|---|---|---|---|---|---|---|---|---|---|---|---|
| T1 | F3 browser->API | E | Authenticated user changes `orgId` in path to read another tenant's invoices | none found in `invoices.controller.ts` | 3 | 4 | Critical | Mitigate | Tenant-scoped query + authZ test `invoices.idor.spec` | @owner | Low |

**Abuse cases:** ...
**Top risks (priority order):** ...
**Accepted risks:** ID, owner, reason, review date
**Open questions / revisit triggers:** ...
```
