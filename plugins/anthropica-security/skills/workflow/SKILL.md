---
name: workflow
description: Phase map of the Anthropica Security skills (Plan > Build > Secure > Ship > Operate) and advice on which skill to use next. Use when starting a task, when unsure which skill applies, or after finishing a step to pick the next one.
---

# Anthropica Security workflow

Phases run in this order. Skip phases that do not apply; loop back when a later phase finds a problem.

**Plan**: Understand the code, requirements, risks and approach before changing anything.
- `/anthropica-security:plan-codebase-analysis`: Orient in an unfamiliar repository and trace a behaviour end to end (entry points, call paths, data flow, git history, tests as spec, conventions) with confidence-labelled findings.
- `/anthropica-security:plan-threat-model`: Threat-model a system or feature - scope, data-flow diagram with trust boundaries, STRIDE per element, abuse cases, LINDDUN privacy, likelihood x impact rating, mitigations, residual risk - as a threat table.

**Build**: Write code, infrastructure and commits.
- `/anthropica-security:build-iac`: Write, review and apply IaC (Terraform/OpenTofu, Bicep, Pulumi, CloudFormation/CDK) with saved plans, locked remote state, pinned versions, policy checks, safe state surgery and drift control.

**Secure**: Check code, dependencies and repository settings for security gaps.
- `/anthropica-security:secure-code-review`: Security-review code or a diff by tracing untrusted input source to sink (injection, IDOR/authZ, authN/session, SSRF, deserialization, path traversal, XSS, CSRF, crypto, secrets, logging, races, mass assignment), proving reachability, rating severity and giving fixes.
- `/anthropica-security:secure-repo-governance`: Audit or set up GitHub repo and org governance with `gh api` - rulesets vs branch protection, required checks, reviews, CODEOWNERS, signed commits, tag rules, environments, secret and code scanning, Dependabot, vulnerability reporting, SECURITY.md, Actions permissions - as a gap table.
- `/anthropica-security:secure-supply-chain`: Harden the software supply chain - lockfiles, dependency audits with reachability triage, pinning actions by SHA and images by digest, update bots, typosquatting and dependency confusion, SBOM, signing and provenance, CI hardening, leaked-secret response.

**Ship**: Run pipelines, decide go/no-go and release.
- `/anthropica-security:ship-github-actions`: Author, harden and debug GitHub Actions workflows (least-privilege token, SHA-pinned actions, OIDC to clouds, script-injection safety, caching, matrices, concurrency, environments, reusable workflows).

**Operate**: Monitor production and respond to incidents; feed what you learn back into Plan.
- `/anthropica-security:operate-incident`: Run a production incident through triage, mitigation (rollback, flag off, scale, failover), comms, handover and a blameless postmortem; write runbooks.

**Any phase**: Cross-cutting skills for documentation, prompts and token use.
- `/anthropica-security:core-documentation`: Write or update docs verified against the repo (README, setup and how-to guides, tutorials, reference, explanations, runbooks) using Diataxis and docs-as-code checks.
- `/anthropica-security:core-prompt-engineering`: Write and debug prompts, system prompts and task specs for Claude, coding agents and subagents (goal, context, constraints, output contract, XML, examples, reasoning) and test them.
- `/anthropica-security:core-token-efficiency`: Keep Claude Code sessions lean without losing correctness - what loads always vs on use, targeted search and reads, bounded output, subagent fan-out, /compact, /clear, handover notes.

## Choosing the next skill
1. Place the work in a phase from the request and the repository state: no agreed approach yet is Plan; uncommitted or unreviewed changes are Build, Test or Review; a merged but unreleased change is Ship; a production problem is Operate.
2. Recommend one skill and say why in one line. Name the one after it.
3. Use it once the user agrees, or straight away if they already asked you to proceed.
4. If no skill here fits, say so rather than stretching one.
