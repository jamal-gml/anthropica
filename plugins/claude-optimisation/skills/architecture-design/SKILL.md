---
name: architecture-design
description: Design or evaluate software architecture, system boundaries, technology choices and implementation plans. Use when planning a new application, major feature, migration or architectural change.
disable-model-invocation: true
---

# Architecture Design

## Method
1. Establish functional requirements, non-functional requirements, constraints, scale and operational environment.
2. Inspect existing architecture before proposing a replacement.
3. Define system boundaries, components, data ownership, trust boundaries and integration contracts.
4. Compare viable options against explicit criteria: complexity, maintainability, security, cost, performance and team capability.
5. Identify failure modes, observability, backup/recovery and deployment concerns.
6. Propose the smallest architecture that meets the requirements; avoid speculative infrastructure.
7. Break implementation into independently testable phases with acceptance criteria.

## Output
- Requirements and assumptions
- Proposed architecture and component responsibilities
- Data flow and integration contracts
- Security and operational considerations
- Alternatives and trade-offs
- Phased implementation plan

Label assumptions clearly. Do not present a diagram or technology choice as validated unless it has been checked against the actual project and constraints.
