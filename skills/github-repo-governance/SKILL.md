---
name: github-repo-governance
description: Set up or audit GitHub repository governance (rulesets and branch protection, CODEOWNERS, required checks, Dependabot, secret scanning, code scanning, environments, templates). Use when creating a repository, hardening an existing one, or reviewing repository settings against a team standard.
---

# GitHub Repository Governance

Inspect current settings before proposing changes, and change repository or organisation settings only with explicit approval.

## Audit commands
```bash
gh api repos/{owner}/{repo} --jq '{default_branch, visibility, allow_squash_merge, allow_merge_commit, delete_branch_on_merge}'
gh api repos/{owner}/{repo}/rulesets
gh api repos/{owner}/{repo}/branches/{branch}/protection   # legacy protection
gh api repos/{owner}/{repo}/environments
gh api repos/{owner}/{repo} --jq .security_and_analysis
```

## Baseline standard
**Default and release branches** (prefer rulesets over legacy protection):
- Require a pull request with at least one approval. Dismiss stale approvals on new pushes. Require code-owner review where `CODEOWNERS` exists.
- Require status checks (CI, lint, PR title) and require branches to be up to date.
- Require conversation resolution. Block force pushes and deletion.
- Require linear history if the team squash-merges. Require signed commits if policy demands it.
- Keep the bypass list empty or limited to a break-glass role.

**Repository files**
- `CODEOWNERS`: map paths to teams, not individuals where possible.
- `.github/dependabot.yml`: cover every ecosystem in the repo, including `github-actions`. Group minor and patch updates.
- `.github/pull_request_template.md` and issue templates.
- `SECURITY.md` with a disclosure route. `LICENSE` for public repositories.
- `.gitignore` covering secrets, environment files and build output.

**Security features**
- Secret scanning and push protection on. Dependabot alerts and security updates on.
- Code scanning (CodeQL or equivalent) for supported languages.

**Merge settings**
- Allow squash merge (PR title as the commit message). Enable auto-delete of head branches.
- Allow merge commits only if environment promotions need them.

**Environments**
- `production`: required reviewers, deployment branch restricted to `main` or release tags, environment-scoped secrets.

## Output
A table of each setting: current value, recommended value, risk if left unchanged. Then the exact changes proposed, as `gh api` calls or file diffs, for approval. Never print secret values.
