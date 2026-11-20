# Keeping branches current and stacking PRs

## Update a PR branch
| Situation | Do |
|---|---|
| Your branch, nobody else pushed, team allows rewrites | `git fetch origin && git rebase origin/<base>`, run tests, then `git push --force-with-lease` (ask first) |
| Shared branch, or reviewers mid-review, or repo prefers merges | `git fetch origin && git merge origin/<base>` then push; or `gh pr update-branch <n>` |
| Platform "Update branch" button / merge queue present | let the platform do it; do not rebase manually |
| No conflicts and repo does not require up-to-date branches | leave it; updating just retriggers CI |

After any update: re-run local checks and confirm CI re-ran on the new head SHA.

## Resolve conflicts
1. Understand both sides: `git log --oneline origin/<base> -- <path>` shows what changed upstream.
2. Resolve to preserve both intents; never "take mine" wholesale on files you did not author.
3. Lockfiles: do not hand-merge; take base's version and re-run the install command.
4. Generated files: regenerate.
5. Run tests before pushing. Mention non-trivial resolutions in a PR comment.

## Stacked PRs
Use when a change is too big for one reviewable PR but parts depend on each other.

**Structure**
```
main <- feature/abc-1-schema      (PR 1, base main)
          <- feature/abc-1-api    (PR 2, base feature/abc-1-schema)
               <- feature/abc-1-ui (PR 3, base feature/abc-1-api)
```
- Each PR is independently reviewable, passes CI and leaves `main` working when merged alone (use flags for incomplete features).
- Title each with its position: `feat(api): add export endpoint [2/3]`; link the stack in each body.
- Keep each layer small; reviewers review layer by layer.

**Merging a stack (squash repos)**
1. Merge PR 1 into `main`.
2. Retarget PR 2: `gh pr edit <n2> --base main`.
3. Because PR 1 was squashed, PR 2 still contains PR 1's original commits. Rebase them away: `git rebase --onto origin/main <old-pr1-tip> feature/abc-1-api`, then force-push with lease (ask first).
4. Repeat for each layer. Tools (`git-town`, `ghstack`, `spr`, Graphite, `git rebase --update-refs`) automate this; use whichever the repo already uses.

**Updating a lower layer**
- `git rebase --update-refs` (git 2.38+) on the top branch moves all intermediate branch refs in one pass; push each with `--force-with-lease`.

**Avoid stacks when** reviewers are scarce (stacks stall), layers cannot ship independently, or the repo uses merge commits and reviewers find the retargeted diffs confusing. Consider one PR with well-separated commits and "review by commit" instead.

## Long-lived branches
- Integrate from base at least daily; drift beyond a few days multiplies conflicts.
- Prefer feature flags and small merges over long-lived feature branches (trunk-based development).
- Release and environment branches receive changes only via promotion PRs or cherry-picks of merged fixes (`ship-release`).
