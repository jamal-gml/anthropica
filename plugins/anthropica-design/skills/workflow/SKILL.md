---
name: workflow
description: Phase map of the Anthropica Design skills (Build > Review) and advice on which skill to use next. Use when starting a task, when unsure which skill applies, or after finishing a step to pick the next one.
---

# Anthropica Design workflow

Phases run in this order. Skip phases that do not apply; loop back when a later phase finds a problem.

**Build**: Write code, infrastructure and commits.
- `/anthropica-design:build-design-system`: Create, extend or audit a design system - tokens (tiers, W3C format, theming, dark mode), accessible component APIs, stories per state, visual regression, versioning, deprecation and design-to-code handoff.
- `/anthropica-design:build-frontend`: Build or change frontend UI to production standard (components, state ownership, data fetching, forms, accessibility, Core Web Vitals, XSS-safe rendering, token storage, Testing Library tests, i18n).

**Review**: Get the change or content reviewed before it merges or publishes.
- `/anthropica-design:review-ui-ux`: Review screens, flows, mockups or built UI for usability, WCAG 2.2 AA accessibility, UI states, microcopy, responsive layout and consistency, with severity-rated findings.

**Any phase**: Cross-cutting skills for documentation, prompts and token use.
- `/anthropica-design:core-documentation`: Write or update docs verified against the repo (README, setup and how-to guides, tutorials, reference, explanations, runbooks) using Diataxis and docs-as-code checks.
- `/anthropica-design:core-prompt-engineering`: Write and debug prompts, system prompts and task specs for Claude, coding agents and subagents (goal, context, constraints, output contract, XML, examples, reasoning) and test them.
- `/anthropica-design:core-token-efficiency`: Keep Claude Code sessions lean without losing correctness - what loads always vs on use, targeted search and reads, bounded output, subagent fan-out, /compact, /clear, handover notes.

## Choosing the next skill
1. Place the work in a phase from the request and the repository state: no agreed approach yet is Plan; uncommitted or unreviewed changes are Build, Test or Review; a merged but unreleased change is Ship; a production problem is Operate.
2. Recommend one skill and say why in one line. Name the one after it.
3. Use it once the user agrees, or straight away if they already asked you to proceed.
4. If no skill here fits, say so rather than stretching one.
