# Editing local history

Only for commits no one else has pulled. Confirm first: `git branch -r --contains <sha>` prints nothing, or the branch is yours alone and the user agreed to a force-push.

Always start with a safety branch: `git branch backup/$(git branch --show-current)-$(date +%s)`.

## Non-interactive rebase (agents have no editor)
`git rebase -i` opens an editor. Drive it with environment variables instead:
- `GIT_SEQUENCE_EDITOR` edits the todo list; `GIT_EDITOR` edits messages.
- Autosquash prepared fixups: `GIT_SEQUENCE_EDITOR=: git rebase -i --autosquash <base>`.
- Custom todo: write it to a file and `GIT_SEQUENCE_EDITOR="cp todo.txt" git rebase -i <base>`.
- Keep message unchanged during `--continue`: `GIT_EDITOR=true git rebase --continue`.

Todo verbs: `pick`, `reword` (r), `edit` (e), `squash` (s, merges messages), `fixup` (f, drops message), `fixup -C` (keep this commit's message), `drop` (d), `exec` (x, run a command).

## Recipes
| Goal | Command |
|---|---|
| Fold a fix into an earlier commit | `git commit --fixup=<sha>` then autosquash |
| Reword an older commit only | `git commit --fixup=reword:<sha> -m "<new message>"` (or `--allow-empty`) then autosquash |
| Squash whole branch into one | `git reset --soft $(git merge-base HEAD <base>)` then `git commit -F msg.txt` |
| Split the last commit | `git reset HEAD~1` (keeps changes unstaged), then stage and commit in parts |
| Split an older commit | todo with `edit <sha>`; at the stop: `git reset HEAD~1`, commit parts, `git rebase --continue` |
| Reorder commits | custom todo with lines reordered |
| Drop a commit | custom todo with `drop <sha>`, or `git rebase --onto <sha>~1 <sha>` |
| Move commits to a new branch | `git branch new-branch`; `git reset --hard <base>` on the old one (ask first) |
| Update onto latest base | `git fetch origin && git rebase origin/<base>` |
| Run tests on every commit | `git rebase -x "<test cmd>" <base>` |

## Conflicts
1. `git status` lists conflicted paths. Resolve each, `git add <path>`.
2. `GIT_EDITOR=true git rebase --continue`.
3. Same conflict repeating across commits: enable `git config rerere.enabled true` for the repo.
4. Lost or unsure: `git rebase --abort` returns to the pre-rebase state.

## Verify after rewriting
- `git range-diff <old-base>..backup/<name> <new-base>..HEAD` shows what changed between the old and new series.
- `git diff backup/<name> HEAD` must be empty when you only reorganised commits (no content change intended).
- Re-run the build and tests on the tip; for important series, `git rebase -x` to test each commit.

## Recovery
- `git reflog` lists every previous HEAD: `git reset --hard HEAD@{n}` restores one (ask first; discards working changes).
- Lost branch tip: `git reflog show <branch>` or `git fsck --lost-found`.
- Delete the backup branch only after the user confirms the result.

## Pushing a rewritten branch
- Own, unshared branch only: `git push --force-with-lease=<branch>:<sha-you-last-fetched> origin <branch>`. The lease fails if someone else pushed; then fetch and inspect instead of forcing.
- Add `--force-if-includes` (git 2.30+) to catch the case where a background fetch updated the remote-tracking ref.
- Never on `main`, release or other shared branches; protected branches will (and should) reject it.
- During review, prefer follow-up commits over rewriting unless the team rebases; rewriting resets reviewers' "viewed" state.
