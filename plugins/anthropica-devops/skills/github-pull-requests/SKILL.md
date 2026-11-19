---
name: github-pull-requests
description: Create, review and merge GitHub pull requests to industry standard (branch naming, small focused PRs, clear descriptions, CI gates, review etiquette, merge strategy). Use when starting a branch, opening or updating a PR, responding to review, or merging.
---

# GitHub Pull Requests

## Branch
1. `git fetch origin` and branch from the latest default branch: `git switch -c <prefix>/<ticket>-<slug> origin/main`.
2. Prefixes: `feature/`, `fix/`, `hotfix/`, `chore/`, `ci/`, `docs/`, `refactor/`. Lowercase, hyphenated, ticket key first when one exists (`feature/abc-123-export-csv`). Follow the repository's convention if it differs.
3. One branch per change. Never push to bot-owned branches (`dependabot/*`, `renovate/*`, release-bot branches).

## Size and scope
- Aim for under ~400 changed lines of reviewable code. Split large work into stacked or sequential PRs.
- Keep refactors, formatting and generated files out of feature PRs, or isolate them in separate commits.

## Opening the PR
- Run the repository's lint, type check and tests locally first.
- Open as **draft** until CI is green and the description is complete, then mark ready for review.
- Title: a Conventional Commit summary (`feat(api): add CSV export`). With squash merges it becomes the commit on `main`.
- Use the repository's PR template if one exists. Otherwise:
  ```
  ## What
  ## Why
  ## How it was tested
  ## Risk and rollback
  Fixes #123
  ```
- Link the issue with a closing keyword. Add screenshots for UI changes and migration notes for schema changes.
- Request reviewers from `CODEOWNERS` or the owning team. Do not self-approve.

## Commands
```bash
gh pr create --draft --base main --title "feat(api): add CSV export" --body-file pr.md
gh pr checks --watch
gh pr ready
gh pr view --comments
```

## Responding to review
- Address every comment: fix it, or reply with the reason. Resolve threads only when addressed.
- Push follow-up commits rather than force-pushing during review, unless the team prefers rebasing.
- Re-request review after significant changes.

## Reviewing someone else's PR
- Review the diff against the stated intent. Check correctness, tests, security, naming and backwards compatibility.
- Label comments by weight: blocking, suggestion, nit. Be specific and cite `file:line`.

## Merging
- Merge only with required checks green and required approvals. Never use admin bypass, never disable checks to merge.
- Default: **squash merge** feature branches into `main`. Use merge commits for environment promotions (`main` → `staging` → `production`) to preserve history. Follow the repository's settings if they differ.
- Delete the branch after merge. Confirm the linked issue closed.

## Output
Report the PR URL, base and head branches, check status and what is blocking merge, if anything.
