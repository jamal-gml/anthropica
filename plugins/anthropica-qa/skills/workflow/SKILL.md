---
name: workflow
description: Phase map of the Anthropica Quality assurance skills (Plan > Test > Review > Ship) and advice on which skill to use next. Use when starting a task, when unsure which skill applies, or after finishing a step to pick the next one.
---

# Anthropica Quality assurance workflow

Phases run in this order. Skip phases that do not apply; loop back when a later phase finds a problem.

**Plan**: Understand the code, requirements, risks and approach before changing anything.
- `/anthropica-qa:plan-codebase-analysis`: Orient in an unfamiliar repository and trace a behaviour end to end (entry points, call paths, data flow, git history, tests as spec, conventions) with confidence-labelled findings.
- `/anthropica-qa:plan-test-strategy`: Plan risk-based testing - risk scoring, test levels, test-design techniques, exploratory charters, entry/exit criteria, test data, traceability and bug reports.

**Test**: Prove the change works, including under load.
- `/anthropica-qa:test-performance`: Measure and fix latency, throughput, memory, database and frontend performance with baselines, percentiles, profilers, EXPLAIN ANALYZE, Core Web Vitals and load tests (k6, Locust, JMeter), reporting before/after numbers.
- `/anthropica-qa:test-validation`: Choose, write and run the smallest tests that prove a change - regression test first for bugs, repo run commands, flaky-test diagnosis and quarantine, mocking boundaries, coverage and mutation signals.

**Review**: Get the change or content reviewed before it merges or publishes.
- `/anthropica-qa:review-pull-request`: Open, review and merge pull requests (branch naming, draft until green, Conventional Commit titles, body template, severity-labelled comments, CI triage, squash vs merge vs rebase, stacked PRs).
- `/anthropica-qa:review-ui-ux`: Review screens, flows, mockups or built UI for usability, WCAG 2.2 AA accessibility, UI states, microcopy, responsive layout and consistency, with severity-rated findings.

**Ship**: Run pipelines, decide go/no-go and release.
- `/anthropica-qa:ship-readiness`: Run an evidence-based go/no-go review (CI on the exact SHA, security, migrations, config, flags, observability, rehearsed rollback, on-call, approvals) with blocking vs non-blocking items and a verdict.

**Any phase**: Cross-cutting skills for documentation, prompts and token use.
- `/anthropica-qa:core-documentation`: Write or update docs verified against the repo (README, setup and how-to guides, tutorials, reference, explanations, runbooks) using Diataxis and docs-as-code checks.
- `/anthropica-qa:core-prompt-engineering`: Write and debug prompts, system prompts and task specs for Claude, coding agents and subagents (goal, context, constraints, output contract, XML, examples, reasoning) and test them.
- `/anthropica-qa:core-token-efficiency`: Keep Claude Code sessions lean without losing correctness - what loads always vs on use, targeted search and reads, bounded output, subagent fan-out, /compact, /clear, handover notes.

## Choosing the next skill
1. Place the work in a phase from the request and the repository state: no agreed approach yet is Plan; uncommitted or unreviewed changes are Build, Test or Review; a merged but unreleased change is Ship; a production problem is Operate.
2. Recommend one skill and say why in one line. Name the one after it.
3. Use it once the user agrees, or straight away if they already asked you to proceed.
4. If no skill here fits, say so rather than stretching one.
