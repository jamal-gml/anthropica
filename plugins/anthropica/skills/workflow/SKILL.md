---
name: workflow
description: Phase map of the Anthropica skills (Plan > Build > Create > Test > Review > Secure > Ship > Operate) and advice on which skill to use next. Use when starting a task, when unsure which skill applies, or after finishing a step to pick the next one.
---

# Anthropica workflow

Phases run in this order. Skip phases that do not apply; loop back when a later phase finds a problem.

**Plan**: Understand the code, requirements, risks and approach before changing anything.
- `/anthropica:plan-architecture`: Design or evaluate software architecture (requirements and scale numbers, options with trade-offs, boundaries, data ownership, sync vs async, consistency, failure modes, migration, ADRs, phased plan).
- `/anthropica:plan-campaign`: Plan marketing campaigns and launches - objective and KPI tree, segments, message house, channel mix by funnel stage, budget and pacing, timeline, UTMs, measurement, experiments and post-campaign review.
- `/anthropica:plan-codebase-analysis`: Orient in an unfamiliar repository and trace a behaviour end to end (entry points, call paths, data flow, git history, tests as spec, conventions) with confidence-labelled findings.
- `/anthropica:plan-test-strategy`: Plan risk-based testing - risk scoring, test levels, test-design techniques, exploratory charters, entry/exit criteria, test data, traceability and bug reports.
- `/anthropica:plan-threat-model`: Threat-model a system or feature - scope, data-flow diagram with trust boundaries, STRIDE per element, abuse cases, LINDDUN privacy, likelihood x impact rating, mitigations, residual risk - as a threat table.

**Build**: Write code, infrastructure and commits.
- `/anthropica:build-backend`: Implement backend APIs and services to production standard (status codes, RFC 9457 errors, pagination, versioning, idempotency keys, validation, authz, transactions, timeouts and retries, background jobs, tests).
- `/anthropica:build-commits`: Stage and write Git commits to standard (Conventional Commits incl. breaking changes, atomic hunks, verified identity and signing, secret scanning, autosquash on unpushed work).
- `/anthropica:build-containers`: Build and review container images and Kubernetes workloads (pinned non-root images, scanning, resources, probes, securityContext, PDB/HPA, graceful shutdown, Helm/Kustomize, rollout/rollback) and debug pods.
- `/anthropica:build-database`: Design schemas and make production-safe database changes (constraints, indexes, EXPLAIN, N+1, expand/contract migrations, online index builds, batched backfills, isolation, pooling, backups) with approval gates for destructive SQL.
- `/anthropica:build-design-system`: Create, extend or audit a design system - tokens (tiers, W3C format, theming, dark mode), accessible component APIs, stories per state, visual regression, versioning, deprecation and design-to-code handoff.
- `/anthropica:build-frontend`: Build or change frontend UI to production standard (components, state ownership, data fetching, forms, accessibility, Core Web Vitals, XSS-safe rendering, token storage, Testing Library tests, i18n).
- `/anthropica:build-iac`: Write, review and apply IaC (Terraform/OpenTofu, Bicep, Pulumi, CloudFormation/CDK) with saved plans, locked remote state, pinned versions, policy checks, safe state surgery and drift control.
- `/anthropica:build-project-setup`: Onboard a repo for AI coding assistants (lean CLAUDE.md/AGENTS.md with verified commands, conventions, gotchas) and scaffold new projects (lockfile, formatter, linter, tests, pre-commit, CI, .gitignore, .editorconfig, README) without overwriting files.

**Create**: Write content and copy.
- `/anthropica:create-marketing-copy`: Write and edit conversion-focused marketing copy (landing pages, emails, ads, social, announcements, push/SMS) that is on-brand, accurate, accessible and testable.
- `/anthropica:create-seo-content`: Plan, write and optimise search-led content - intent, SERP analysis, keyword clusters, briefs, E-E-A-T, on-page and schema, refreshing or pruning pages.

**Test**: Prove the change works, including under load.
- `/anthropica:test-performance`: Measure and fix latency, throughput, memory, database and frontend performance with baselines, percentiles, profilers, EXPLAIN ANALYZE, Core Web Vitals and load tests (k6, Locust, JMeter), reporting before/after numbers.
- `/anthropica:test-validation`: Choose, write and run the smallest tests that prove a change - regression test first for bugs, repo run commands, flaky-test diagnosis and quarantine, mocking boundaries, coverage and mutation signals.

**Review**: Get the change or content reviewed before it merges or publishes.
- `/anthropica:review-pull-request`: Open, review and merge pull requests (branch naming, draft until green, Conventional Commit titles, body template, severity-labelled comments, CI triage, squash vs merge vs rebase, stacked PRs).
- `/anthropica:review-seo`: Audit technical SEO from evidence (headers, robots.txt, sitemaps, rendered HTML, crawl and search-console exports) - crawlability, indexation, canonicals, redirects, rendering, Core Web Vitals, structured data, hreflang, migrations.
- `/anthropica:review-ui-ux`: Review screens, flows, mockups or built UI for usability, WCAG 2.2 AA accessibility, UI states, microcopy, responsive layout and consistency, with severity-rated findings.

**Secure**: Check code, dependencies and repository settings for security gaps.
- `/anthropica:secure-code-review`: Security-review code or a diff by tracing untrusted input source to sink (injection, IDOR/authZ, authN/session, SSRF, deserialization, path traversal, XSS, CSRF, crypto, secrets, logging, races, mass assignment), proving reachability, rating severity and giving fixes.
- `/anthropica:secure-repo-governance`: Audit or set up GitHub repo and org governance with `gh api` - rulesets vs branch protection, required checks, reviews, CODEOWNERS, signed commits, tag rules, environments, secret and code scanning, Dependabot, vulnerability reporting, SECURITY.md, Actions permissions - as a gap table.
- `/anthropica:secure-supply-chain`: Harden the software supply chain - lockfiles, dependency audits with reachability triage, pinning actions by SHA and images by digest, update bots, typosquatting and dependency confusion, SBOM, signing and provenance, CI hardening, leaked-secret response.

**Ship**: Run pipelines, decide go/no-go and release.
- `/anthropica:ship-ci-cd`: Design CI/CD pipelines on any platform (stages, build once and promote, environment gates, rolling/blue-green/canary/flags, expand-contract migrations, rollback vs roll-forward, monorepos, speed, OIDC).
- `/anthropica:ship-github-actions`: Author, harden and debug GitHub Actions workflows (least-privilege token, SHA-pinned actions, OIDC to clouds, script-injection safety, caching, matrices, concurrency, environments, reusable workflows).
- `/anthropica:ship-readiness`: Run an evidence-based go/no-go review (CI on the exact SHA, security, migrations, config, flags, observability, rehearsed rollback, on-call, approvals) with blocking vs non-blocking items and a verdict.
- `/anthropica:ship-release`: Version, tag, publish and promote releases (semver, changelogs, annotated tags, release-please/changesets/semantic-release, registries, hotfixes, rollback, yanking a bad version).

**Operate**: Monitor production and respond to incidents; feed what you learn back into Plan.
- `/anthropica:operate-incident`: Run a production incident through triage, mitigation (rollback, flag off, scale, failover), comms, handover and a blameless postmortem; write runbooks.
- `/anthropica:operate-observability`: Design logs, metrics, OpenTelemetry traces, dashboards, SLOs with multiwindow burn-rate alerts, and telemetry cost control.

**Any phase**: Cross-cutting skills for documentation, prompts and token use.
- `/anthropica:core-documentation`: Write or update docs verified against the repo (README, setup and how-to guides, tutorials, reference, explanations, runbooks) using Diataxis and docs-as-code checks.
- `/anthropica:core-prompt-engineering`: Write and debug prompts, system prompts and task specs for Claude, coding agents and subagents (goal, context, constraints, output contract, XML, examples, reasoning) and test them.
- `/anthropica:core-token-efficiency`: Keep Claude Code sessions lean without losing correctness - what loads always vs on use, targeted search and reads, bounded output, subagent fan-out, /compact, /clear, handover notes.

## Choosing the next skill
1. Place the work in a phase from the request and the repository state: no agreed approach yet is Plan; uncommitted or unreviewed changes are Build, Test or Review; a merged but unreleased change is Ship; a production problem is Operate.
2. Recommend one skill and say why in one line. Name the one after it.
3. Use it once the user agrees, or straight away if they already asked you to proceed.
4. If no skill here fits, say so rather than stretching one.
