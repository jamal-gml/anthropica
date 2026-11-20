---
name: secure-repo-governance
description: Audit or set up GitHub repo and org governance with `gh api` - rulesets vs branch protection, required checks, reviews, CODEOWNERS, signed commits, tag rules, environments, secret and code scanning, Dependabot, vulnerability reporting, SECURITY.md, Actions permissions - as a gap table. Use when creating or hardening a repo, required checks are stuck, or settings must match a standard. Not for workflow code (use `ship-github-actions`).
---

# Repository Governance

Read the effective configuration, compare it with the baseline, and propose exact, reversible changes. Never change settings without approval.

## Procedure
1. **Identify scope and access.** `gh auth status` (scopes: `repo`, plus `admin:org`/`read:org` for org settings; some endpoints need admin and return 404/403 otherwise - record "not visible" rather than "disabled"). Get `default_branch`, visibility, and the plan (some features need GitHub Advanced Security on private repos).
2. **Collect current state** read-only. Read `references/audit-commands.md` for the full `gh api` command set (repo, rulesets, classic protection, Actions, security features, environments, access, org).
3. **Check effective branch rules,** not just one source: `gh api repos/{owner}/{repo}/rules/branches/{branch}` returns rules from all repo and org rulesets that apply; also read classic protection. Both layer: the most restrictive combination wins.
4. **Check required status checks actually match** check-run names produced by CI (rules below). List recent names with `gh api repos/{owner}/{repo}/commits/{sha}/check-runs --jq '.check_runs[].name'`.
5. **Check files:** `CODEOWNERS` (and `gh api repos/{owner}/{repo}/codeowners/errors`), `SECURITY.md`, `.github/dependabot.yml` covering every ecosystem and directory, PR/issue templates, `.gitignore` for secrets and env files, `LICENSE` for public repos.
6. **Compare with the baseline** below and the team's written standard if one exists (the stricter wins; flag conflicts). Rate each gap.
7. **Propose changes** as `gh api` calls, ruleset JSON or file diffs. Read `references/templates.md` for ruleset JSON, CODEOWNERS, SECURITY.md and dependabot.yml templates. Prefer rulesets for new configuration; migrate classic protection only with a plan to avoid a window without protection.
8. **After approval, apply and verify:** re-read the setting, then prove it with a harmless test (a draft PR shows the required checks and reviewers; `gh api .../rules/branches/{branch}` shows the rule).
9. Hand dependency, pinning and CI-token findings to `secure-supply-chain`; workflow edits to `ship-github-actions`.

### Baseline (default and release branches, tags, settings)
| Area | Recommended |
|---|---|
| Branch rules | Ruleset on `~DEFAULT_BRANCH` and release branches: PR required, >= 1 approval (2 for high-risk), dismiss stale approvals, require approval of most recent push, code-owner review, conversation resolution, required checks with strict up-to-date (or merge queue), block force push and deletion, linear history if squash-only |
| Bypass | Empty, or a break-glass team with "pull requests only" bypass mode; no individual users; admins not exempt by default |
| Tags | Tag ruleset on `v*`: restrict creation/update/deletion to release role or workflow; legacy tag protection is retired - use rulesets |
| Signed commits | Required where policy demands provenance of authorship; check merge-method compatibility and bot/automation signing first |
| Merge settings | Squash (PR title as message), auto-delete head branches, merge commits only if promotions need them, auto-merge allowed only with required checks |
| Environments | `production`: required reviewers, prevent self-review, deployment branches/tags restricted to protected refs, environment-scoped secrets, admin bypass off |
| Security features | Secret scanning + push protection (+ validity checks, non-provider patterns where available), Dependabot alerts + security updates, code scanning (CodeQL default setup or equivalent) on PRs and default branch, private vulnerability reporting on |
| Actions | Default `GITHUB_TOKEN` read-only, Actions cannot approve PRs, allowed actions = selected (owner-owned + verified + explicit list), require full-SHA pinning if available, fork PR workflows need approval, self-hosted runners not on public repos |
| Access | Teams not individuals, least privilege (write for contributors, maintain/admin for few), outside collaborators reviewed, deploy keys read-only unless required, webhooks with secrets and HTTPS |
| Org | 2FA/SSO required, base permission read or none, restrict repo creation/visibility change/deletion/forking of private repos, fine-grained PAT approval and classic PAT restriction, app installation approval, audit-log streaming, org-wide security configuration applied to all repos |

## Rules
- **Required check names:** the context is the job's `name:` (or job ID), not the workflow name; matrix jobs add `(param, ...)`; reusable-workflow jobs appear as `caller / callee`. Renaming a job silently blocks merges ("Expected - waiting for status"). Pin the check's source app (GitHub Actions) in rulesets so another app cannot satisfy it.
- **Skipped workflows:** a required check from a workflow filtered by `paths:` or `branches:` never reports and blocks the PR. Use an always-running job that reports success when nothing relevant changed, or job-level `if:` (skipped jobs report success).
- **Merge queue:** required workflows must also trigger on `merge_group`, or the queue stalls.
- **CODEOWNERS:** last matching pattern wins; owners need write access and teams must be visible; invalid lines are ignored silently (check `codeowners/errors`). Own `CODEOWNERS`, `.github/workflows/` and security-sensitive paths by a security-aware team.
- **Rulesets vs classic:** rulesets layer (repo + org), support tags and push rules, bypass lists and evaluate mode (Enterprise) - use evaluate mode before enforcing org-wide. Classic protection has one rule per pattern and an "include administrators" toggle.
- A 404 from a protection or settings endpoint can mean "not configured" or "no permission": confirm with your permission level before reporting a gap.
- Rate gaps: Critical = no protection on default/release branch, admin-only bypass anyone uses, public repo with secrets exposure path, `GITHUB_TOKEN` write-all with `pull_request_target`; High = no required reviews or checks, push protection off, prod environment without reviewers; Medium = no code scanning, Dependabot incomplete, stale approvals kept; Low = templates, labels, hygiene.

### Common failure modes
- Protection on `main` only, while release branches, `prod` branches and tags are open.
- Approval required but stale approvals kept, so the author pushes new code after review.
- CODEOWNERS present but code-owner review not required, or the owning team lacks write access.
- Admins or a broad team on the bypass list "temporarily", used routinely.
- Push protection on, but bypass reasons never reviewed.
- Environment reviewers set, but deployment branch policy `null`, so any branch can deploy with production secrets.
- Dependabot alerts on, but no `dependabot.yml`, so version updates (including Actions) never arrive.

## Safety gates
- Read-only until the user approves each change; show the exact `gh api -X PUT/POST/PATCH` call or ruleset JSON first. Ask before changing org-level settings, bypass lists, visibility, or deleting protection.
- Enable new rulesets in `evaluate` mode or on a test branch where possible; never leave a branch unprotected mid-migration.
- Never weaken protection to get a PR merged, use admin bypass, or self-approve. Never print secret values (list secret names only).

## Output
```markdown
## Governance audit: <owner>/<repo> (date, viewer permission)
| Area | Setting | Current | Recommended | Gap risk | Evidence (command) |
|---|---|---|---|---|---|
| Branch rules | Required approvals | 0 | >= 1, dismiss stale | High | `rules/branches/main` |

**Not visible with current permissions:** ...
**Conflicts with team standard:** ...
**Proposed changes (awaiting approval):** numbered `gh api` calls / ruleset JSON / file diffs, each with rollback
```

---

**Phase:** Secure. **Next:** when this step is done, offer Ship: `/anthropica-infrastructure:ship-github-actions`. Phase map: `/anthropica-infrastructure:workflow`.
