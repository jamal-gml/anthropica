---
name: review-pull-request
description: Open, review and merge pull requests (branch naming, draft until green, Conventional Commit titles, body template, severity-labelled comments, CI triage, squash vs merge vs rebase, stacked PRs). Use when starting a branch, opening or updating a PR, reviewing a diff, answering review, fixing failing PR checks, or merging. Not for commits (use `build-commits`).
---

# Pull Requests

Get a focused change reviewed and merged through the repository's gates, never around them.

## Procedure: author
1. **Branch from fresh base:** `git fetch origin && git switch -c <prefix>/<key>-<slug> origin/<default>`. Find the default with `gh repo view --json defaultBranchRef -q .defaultBranchRef.name`.
2. **Read repo rules:** `.github/pull_request_template.md`, `CONTRIBUTING.md`, `CODEOWNERS`, required checks (`gh api repos/{owner}/{repo}/rules/branches/<base>`).
3. **Check locally** what CI runs (lint, types, tests, build). Commits per `build-commits`.
4. **Ask, then push:** `git push -u origin HEAD`.
5. **Open as draft:** `gh pr create --draft --base <base> --title "<type(scope): summary>" --body-file pr.md`.
6. **Watch CI:** `gh pr checks --watch --fail-fast`. Triage failures (Rules). Only when green and the body is complete: `gh pr ready`.
7. **Request reviewers** from `CODEOWNERS` or the owning team: `gh pr edit --add-reviewer <user|org/team>`.
8. **Respond to review:** address every thread (fix, or reply with the reason); push follow-up commits; re-request review (`gh pr edit --add-reviewer`) after substantive changes.
9. **Merge** when gates pass and the user approves (Merge strategy below). Delete the branch; confirm the linked issue closed.

## Procedure: reviewer
1. **Context:** `gh pr view <n> --json title,body,baseRefName,headRefName,files,reviewDecision,statusCheckRollup` and the linked issue. Know the intent before reading code.
2. **Diff:** `gh pr diff <n>`; check out for anything non-trivial: `gh pr checkout <n>` then run tests.
3. **Review in passes:** (a) does it do what the issue asks, nothing more; (b) correctness and edge cases; (c) tests prove (b); (d) security and data; (e) compatibility, migrations, rollout; (f) readability. Read `references/reviewing.md` for the full checklist and comment templates.
4. **Comment** with a severity label and `file:line`; propose code via suggestion blocks.
5. **Submit one review:** `gh pr review <n> --approve | --request-changes | --comment --body-file review.md`. Request changes only for blocking items.

## Rules
- **Branch prefixes:** `feature/`, `fix/`, `hotfix/`, `chore/`, `ci/`, `docs/`, `refactor/`; lowercase, hyphenated, issue key first (`feature/abc-123-export-csv`). Repo convention wins.
- **Size:** aim under ~400 changed lines of hand-written code, one concern per PR. Larger: split into a stack (read `references/branch-maintenance.md`). Generated files and mass renames in their own PR or commit.
- **Title** = Conventional Commit (`feat(api)!: drop v1 export`); with squash merge it becomes the commit on the base branch and drives release tooling. Under 72 chars.
- **Body** (when no template): `## What`, `## Why`, `## How tested` (commands and results), `## Risk and rollback`, `## Notes for reviewers`, then `Fixes #123` / `Closes ABC-42`. Add screenshots for UI, migration and rollout notes for schema or config changes, and `BREAKING CHANGE:` text for breaking ones.
- **Draft until green.** Mark ready only with CI passing, conflicts resolved and the body complete.
- **Comment severity labels:** `blocking:` (must change before merge: bug, security, data loss, broken contract, missing test for new behaviour), `suggestion:` (better approach, author decides), `question:` (need understanding), `nit:` (style, never blocks), `praise:`. Unlabelled comments are read as blocking; label them.
- **CI triage:** `gh pr checks <n>`, then `gh run view <run-id> --log-failed`. Find the first real error. Compare with the base branch: `gh run list --branch <base> --workflow <file> -L 5`. If it fails on base too, it is not this PR's bug: report it, do not "fix" it here. Flaky: re-run once (`gh run rerun <id> --failed`); a second failure is real. Read `references/ci-triage.md` for failure patterns.
- **Keep current:** rebase or merge base into the branch only when there are conflicts, the repo requires up-to-date branches, or base changed code you depend on. Rebase only your own unshared branch; otherwise `git merge origin/<base>` or `gh pr update-branch`.
- **Never:** self-approve, approve your own changes via a second account, merge with `--admin`, disable or rename required checks, dismiss others' reviews to get through, push to someone else's branch without asking, or resolve a reviewer's thread for them when the point is unaddressed.

## Merge strategy
| Strategy | Use when | Effect |
|---|---|---|
| Squash | default for feature/fix PRs into the default branch | one commit; PR title becomes the message; branch history discarded |
| Merge commit | environment promotions (`main` -> `staging` -> `production`), long-lived release branches, stacked series that must keep commits | preserves all commits and the merge point; tags and SHAs stay valid |
| Rebase | linear-history repos where every commit is already atomic and well-written | replays commits, new SHAs; no merge commit |

- Repo settings and rulesets override this table; check with `gh repo view --json squashMergeAllowed,mergeCommitAllowed,rebaseMergeAllowed`.
- Use the merge queue or auto-merge if configured: `gh pr merge <n> --squash --auto --delete-branch`.
- Before merging verify: required checks green, required approvals present, no unresolved blocking threads, title is a valid Conventional Commit, branch up to date if required.

## Safety gates
- Ask before: pushing, marking ready, requesting reviewers on someone else's behalf, merging, closing PRs, deleting branches, posting review comments as the user.
- Never bypass branch protection or rulesets, even when you have permission to.
- For hotfixes, the gate is shorter, not absent: still needs CI and an approver.

## Output
```
PR: <url>  (#<n>, draft | ready | merged)
Base <- head: <base> <- <branch>
Title: <type(scope): summary>
Checks: <passed>/<total>  failing: <names + first error line>
Reviews: <approved | changes requested | pending> (<reviewers>)
Blocking: <list or "none">
Next action: <what and who>
```
For a review: `Verdict: approve | request changes | comment`, then findings grouped by `blocking`, `suggestion`, `question`, `nit`, each with `file:line` and a concrete fix.

---

**Phase:** Review. **Next:** when this step is done, offer Secure: `/anthropica-devops:secure-repo-governance`, `/anthropica-devops:secure-supply-chain`. Phase map: `/anthropica-devops:workflow`.
