---
name: workflow
description: Phase map of the Anthropica Software engineering skills (Plan > Build > Test > Review > Secure) and advice on which skill to use next. Use when starting a task, when unsure which skill applies, or after finishing a step to pick the next one.
---

# Anthropica Software engineering workflow

Phases run in this order. Skip phases that do not apply; loop back when a later phase finds a problem.

**Plan**: Understand the code, requirements, risks and approach before changing anything.
- `/anthropica-engineering:plan-architecture`: Design or evaluate software architecture (requirements and scale numbers, options with trade-offs, boundaries, data ownership, sync vs async, consistency, failure modes, migration, ADRs, phased plan).
- `/anthropica-engineering:plan-codebase-analysis`: Orient in an unfamiliar repository and trace a behaviour end to end (entry points, call paths, data flow, git history, tests as spec, conventions) with confidence-labelled findings.

**Build**: Write code, infrastructure and commits.
- `/anthropica-engineering:build-backend`: Implement backend APIs and services to production standard (status codes, RFC 9457 errors, pagination, versioning, idempotency keys, validation, authz, transactions, timeouts and retries, background jobs, tests).
- `/anthropica-engineering:build-commits`: Stage and write Git commits to standard (Conventional Commits incl. breaking changes, atomic hunks, verified identity and signing, secret scanning, autosquash on unpushed work).
- `/anthropica-engineering:build-database`: Design schemas and make production-safe database changes (constraints, indexes, EXPLAIN, N+1, expand/contract migrations, online index builds, batched backfills, isolation, pooling, backups) with approval gates for destructive SQL.
- `/anthropica-engineering:build-frontend`: Build or change frontend UI to production standard (components, state ownership, data fetching, forms, accessibility, Core Web Vitals, XSS-safe rendering, token storage, Testing Library tests, i18n).
- `/anthropica-engineering:build-project-setup`: Onboard a repo for AI coding assistants (lean CLAUDE.md/AGENTS.md with verified commands, conventions, gotchas) and scaffold new projects (lockfile, formatter, linter, tests, pre-commit, CI, .gitignore, .editorconfig, README) without overwriting files.

**Test**: Prove the change works, including under load.
- `/anthropica-engineering:test-performance`: Measure and fix latency, throughput, memory, database and frontend performance with baselines, percentiles, profilers, EXPLAIN ANALYZE, Core Web Vitals and load tests (k6, Locust, JMeter), reporting before/after numbers.
- `/anthropica-engineering:test-validation`: Choose, write and run the smallest tests that prove a change - regression test first for bugs, repo run commands, flaky-test diagnosis and quarantine, mocking boundaries, coverage and mutation signals.

**Review**: Get the change or content reviewed before it merges or publishes.
- `/anthropica-engineering:review-pull-request`: Open, review and merge pull requests (branch naming, draft until green, Conventional Commit titles, body template, severity-labelled comments, CI triage, squash vs merge vs rebase, stacked PRs).

**Secure**: Check code, dependencies and repository settings for security gaps.
- `/anthropica-engineering:secure-code-review`: Security-review code or a diff by tracing untrusted input source to sink (injection, IDOR/authZ, authN/session, SSRF, deserialization, path traversal, XSS, CSRF, crypto, secrets, logging, races, mass assignment), proving reachability, rating severity and giving fixes.

**Any phase**: Cross-cutting skills for documentation, prompts and token use.
- `/anthropica-engineering:core-documentation`: Write or update docs verified against the repo (README, setup and how-to guides, tutorials, reference, explanations, runbooks) using Diataxis and docs-as-code checks.
- `/anthropica-engineering:core-prompt-engineering`: Write and debug prompts, system prompts and task specs for Claude, coding agents and subagents (goal, context, constraints, output contract, XML, examples, reasoning) and test them.
- `/anthropica-engineering:core-token-efficiency`: Keep Claude Code sessions lean without losing correctness - what loads always vs on use, targeted search and reads, bounded output, subagent fan-out, /compact, /clear, handover notes.

## Choosing the next skill
1. Place the work in a phase from the request and the repository state: no agreed approach yet is Plan; uncommitted or unreviewed changes are Build, Test or Review; a merged but unreleased change is Ship; a production problem is Operate.
2. Recommend one skill and say why in one line. Name the one after it.
3. Use it once the user agrees, or straight away if they already asked you to proceed.
4. If no skill here fits, say so rather than stretching one.
