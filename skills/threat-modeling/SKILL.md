---
name: threat-modeling
description: Threat-model a system, feature or architecture change (data flows, trust boundaries, STRIDE threats, risk rating, mitigations) before or during design. Use when designing a feature that handles authentication, payments, personal data, file uploads, external integrations or new infrastructure, or when a security design review is requested.
---

# Threat Modeling

Answer four questions: What are we building? What can go wrong? What are we doing about it? Did we do a good job?

## 1. Model the system
- Inventory components, actors (users, admins, services, third parties), data stores and external dependencies.
- Draw the data flows. Mark **trust boundaries** where privilege or ownership changes: internet to edge, edge to service, service to database, tenant to tenant, CI to production.
- Classify data: credentials, personal data, payment data, secrets, internal-only, public.
- Base this on the actual code, infrastructure and diagrams. Label anything assumed.

## 2. Identify threats (STRIDE per element and boundary crossing)
| Threat | Question | Typical controls |
|---|---|---|
| **S**poofing | Can someone pretend to be another user or service? | Strong auth, MFA, mTLS, signed tokens |
| **T**ampering | Can data or code be modified in transit or at rest? | TLS, integrity checks, signed artifacts, least privilege |
| **R**epudiation | Can actions be denied later? | Audit logs, immutable logging, timestamps |
| **I**nformation disclosure | Can data leak? | Authorisation checks, encryption, minimal logging, tenant isolation |
| **D**enial of service | Can it be made unavailable? | Rate limits, quotas, timeouts, autoscaling |
| **E**levation of privilege | Can someone gain more rights? | Resource-level authZ, input validation, sandboxing |

Also consider abuse cases (business-logic misuse, enumeration, fraud) and supply-chain paths (dependencies, CI, build pipeline).

## 3. Rate and decide
- Rate each threat by likelihood × impact (Low, Medium, High, Critical), considering existing controls.
- Decide per threat: mitigate, eliminate (remove the feature or data), transfer, or accept (with a named owner and a reason).

## 4. Validate
- Each mitigation maps to a testable control: a test case, configuration check or monitoring alert.
- Revisit the model when the design changes.

## Output
System summary and data-flow description (with trust boundaries), a threat table (ID, element, STRIDE category, threat, likelihood, impact, existing controls, mitigation, owner, status), the top risks in priority order, and open questions.
