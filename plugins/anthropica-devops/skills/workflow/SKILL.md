---
name: workflow
description: Phase map of the Anthropica DevOps skills (Plan > Build > Review > Secure > Ship > Operate) and advice on which skill to use next. Use when starting a task, when unsure which skill applies, or after finishing a step to pick the next one.
---

# Anthropica DevOps workflow

Phases run in this order. Skip phases that do not apply; loop back when a later phase finds a problem.

**Plan**: Understand the code, requirements, risks and approach before changing anything.
- `/anthropica-devops:plan-codebase-analysis`: Orient in an unfamiliar repository and trace a behaviour end to end (entry points, call paths, data flow, git history, tests as spec, conventions) with confidence-labelled findings.

**Build**: Write code, infrastructure and commits.
- `/anthropica-devops:build-commits`: Stage and write Git commits to standard (Conventional Commits incl. breaking changes, atomic hunks, verified identity and signing, secret scanning, autosquash on unpushed work).
- `/anthropica-devops:build-containers`: Build and review container images and Kubernetes workloads (pinned non-root images, scanning, resources, probes, securityContext, PDB/HPA, graceful shutdown, Helm/Kustomize, rollout/rollback) and debug pods.
- `/anthropica-devops:build-iac`: Write, review and apply IaC (Terraform/OpenTofu, Bicep, Pulumi, CloudFormation/CDK) with saved plans, locked remote state, pinned versions, policy checks, safe state surgery and drift control.

**Review**: Get the change or content reviewed before it merges or publishes.
- `/anthropica-devops:review-pull-request`: Open, review and merge pull requests (branch naming, draft until green, Conventional Commit titles, body template, severity-labelled comments, CI triage, squash vs merge vs rebase, stacked PRs).

**Secure**: Check code, dependencies and repository settings for security gaps.
- `/anthropica-devops:secure-repo-governance`: Audit or set up GitHub repo and org governance with `gh api` - rulesets vs branch protection, required checks, reviews, CODEOWNERS, signed commits, tag rules, environments, secret and code scanning, Dependabot, vulnerability reporting, SECURITY.md, Actions permissions - as a gap table.
- `/anthropica-devops:secure-supply-chain`: Harden the software supply chain - lockfiles, dependency audits with reachability triage, pinning actions by SHA and images by digest, update bots, typosquatting and dependency confusion, SBOM, signing and provenance, CI hardening, leaked-secret response.

**Ship**: Run pipelines, decide go/no-go and release.
- `/anthropica-devops:ship-ci-cd`: Design CI/CD pipelines on any platform (stages, build once and promote, environment gates, rolling/blue-green/canary/flags, expand-contract migrations, rollback vs roll-forward, monorepos, speed, OIDC).
- `/anthropica-devops:ship-github-actions`: Author, harden and debug GitHub Actions workflows (least-privilege token, SHA-pinned actions, OIDC to clouds, script-injection safety, caching, matrices, concurrency, environments, reusable workflows).
- `/anthropica-devops:ship-readiness`: Run an evidence-based go/no-go review (CI on the exact SHA, security, migrations, config, flags, observability, rehearsed rollback, on-call, approvals) with blocking vs non-blocking items and a verdict.
- `/anthropica-devops:ship-release`: Version, tag, publish and promote releases (semver, changelogs, annotated tags, release-please/changesets/semantic-release, registries, hotfixes, rollback, yanking a bad version).

**Operate**: Monitor production and respond to incidents; feed what you learn back into Plan.
- `/anthropica-devops:operate-incident`: Run a production incident through triage, mitigation (rollback, flag off, scale, failover), comms, handover and a blameless postmortem; write runbooks.
- `/anthropica-devops:operate-observability`: Design logs, metrics, OpenTelemetry traces, dashboards, SLOs with multiwindow burn-rate alerts, and telemetry cost control.

**Any phase**: Cross-cutting skills for documentation, prompts and token use.
- `/anthropica-devops:core-documentation`: Write or update docs verified against the repo (README, setup and how-to guides, tutorials, reference, explanations, runbooks) using Diataxis and docs-as-code checks.
- `/anthropica-devops:core-prompt-engineering`: Write and debug prompts, system prompts and task specs for Claude, coding agents and subagents (goal, context, constraints, output contract, XML, examples, reasoning) and test them.
- `/anthropica-devops:core-token-efficiency`: Keep Claude Code sessions lean without losing correctness - what loads always vs on use, targeted search and reads, bounded output, subagent fan-out, /compact, /clear, handover notes.

## Choosing the next skill
1. Place the work in a phase from the request and the repository state: no agreed approach yet is Plan; uncommitted or unreviewed changes are Build, Test or Review; a merged but unreleased change is Ship; a production problem is Operate.
2. Recommend one skill and say why in one line. Name the one after it.
3. Use it once the user agrees, or straight away if they already asked you to proceed.
4. If no skill here fits, say so rather than stretching one.
