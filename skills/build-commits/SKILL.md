---
name: build-commits
description: Stage and write Git commits to standard (Conventional Commits incl. breaking changes, atomic hunks, verified identity and signing, secret scanning, autosquash on unpushed work). Use when staging, committing, amending, splitting or rewording commits, cleaning history before a push, or after a secret was committed. Not for PRs (use `review-pull-request`).
---

# Git Commits

Produce small, revertable commits with the right identity, honest messages and no secrets.

## Procedure
1. **Read local rules first:** `CONTRIBUTING.md`, `CLAUDE.md`, `commitlint.config.*`, `.husky/`, `.pre-commit-config.yaml`, `lefthook.yml`. Repository rules override this skill.
2. **Check branch:** `git branch --show-current`. On `main`, `master`, `develop` or a release branch, stop and branch first (naming in `review-pull-request`). Never commit to bot branches (`dependabot/*`, `renovate/*`, `release-please--*`, `changeset-release/*`).
3. **Check identity:** `git config user.name; git config user.email` must be the user's. If unset or wrong, set per repo (`git config user.email ...`) or per command (`git -c user.name=... -c user.email=... commit`). Never invent one; ask.
4. **Inspect:** `git status --short`, `git diff`, `git diff --staged`. Know every file before staging it.
5. **Stage by intent, not by directory:**
   - Files: `git add <path>`. Never `git add -A` / `git add .` without reading `git status` first.
   - Hunks: `git add -p <path>`; with no interactive terminal, write the hunks to a patch and `git apply --cached hunks.diff`.
   - Unstage: `git restore --staged <path>`. Re-check with `git diff --staged --stat`.
6. **Scan staged content for secrets** (patterns below). Run the repo hook, or `gitleaks protect --staged` / `detect-secrets-hook` if installed.
7. **Run the fast checks the hooks run** (format, lint, unit tests for touched modules). Fix failures; never bypass.
8. **Write the message** (format below). Use `git commit -F msg.txt` for multi-line bodies.
9. **Verify:** `git log -1 --format='%h %an <%ae> | %cn <%ce>%n%B'`; if signing is required, `git log -1 --show-signature`.
10. **Before push:** `git log origin/<base>..HEAD --format='%h %an <%ae> %s'` and `git diff origin/<base>..HEAD` for a final identity and secret check. Ask before pushing.

## Message format (Conventional Commits 1.0)
```
<type>(<scope>)!: <imperative summary, <=72 chars, no trailing period>

<body: why, what it affects, alternatives rejected; wrap at 72>

BREAKING CHANGE: <what breaks and how to migrate>
Refs: #123
Fixes ABC-42
```
| Type | Use for | Release effect |
|---|---|---|
| `feat` | new user-visible capability | minor |
| `fix` | bug fix | patch |
| `perf` | faster, same behaviour | patch |
| `refactor` | restructure, same behaviour | none |
| `docs`, `test`, `style` | docs, tests, formatting only | none |
| `build`, `ci`, `chore` | build system and deps, pipelines, tooling | none (some tools patch on `build(deps)`) |
| `revert` | `revert: <original summary>` + body `This reverts commit <sha>.` | as original |

- **Breaking:** `!` after type/scope and/or a `BREAKING CHANGE:` footer (exact uppercase spelling, colon, space). Either forces a major bump (minor while `0.y.z`); the footer carries migration steps.
- **Scope:** lowercase module/package (`api`, `auth`, `deps`). In monorepos use the package name release tooling maps. Omit rather than invent.
- **Footers:** `Token: value` or `Token #value`, one blank line before the block. Issue keys go in a footer or as a summary suffix `(ABC-42)`, per repo convention.
- **Body:** explain why; the diff shows what. Note behaviour changes, risk and follow-ups.
- **Forbidden:** "WIP", "fix stuff", "update file", "address review", ticket-only or emoji-only summaries, past tense ("added").

## Rules
- **Atomic:** each commit builds, passes tests and reverts cleanly alone. Separate formatting/renames, dependency bumps, refactors and features. Rename with `git mv` in its own commit so history follows the file.
- **Split heuristic:** if the summary needs "and", split. Over ~400 hand-written changed lines in one commit needs a stated reason.
- **Generated files** (lockfiles, codegen) go in the commit that caused them; regenerate with the tool, never hand-edit.
- **Never commit:** `.env*` (except `.env.example` with placeholders), private keys (`*.pem`, `id_rsa*`, `*.p12`, `*.pfx`, `*.key`), cloud credential files, `*.tfstate`, secret `*.tfvars`, kubeconfigs, DB dumps, binaries over ~5 MB (LFS if the repo uses it), `node_modules/`, build output, IDE folders. Add the pattern to `.gitignore` in the same change.
- **Secret patterns to grep the staged diff for:** `AKIA[0-9A-Z]{16}`, `-----BEGIN [A-Z ]*PRIVATE KEY-----`, `ghp_|gho_|ghs_|github_pat_`, `glpat-`, `xox[baprs]-`, `sk_live_|rk_live_`, `AIza[0-9A-Za-z_-]{35}`, JWTs `eyJ[A-Za-z0-9_-]{10,}\.`, `(password|secret|token|api_key)\s*[:=]\s*\S{8,}`, URLs with `://user:pass@`.
- **Rewrite only unpushed history.** `git branch -r --contains <sha>` empty = unpushed. If a remote branch others use contains it, add a new commit or `git revert <sha>`; never rewrite.
- **Fixups:** `git commit --fixup=<sha>`, then `GIT_SEQUENCE_EDITOR=: git rebase -i --autosquash <base>` (runs without an editor). `--fixup=reword:<sha>` (git 2.32+) changes only the message.
- **Amend** only the newest unpushed commit: `git commit --amend --no-edit` (content) or `--amend -F msg.txt` (message). Add `--reset-author` if the identity was wrong.
- **Own branch already pushed and not shared:** after rewriting, `git push --force-with-lease=<branch>:<expected-sha>` only; never plain `--force`, never on shared branches.
- **Hooks and signing:** never `--no-verify` or `-c commit.gpgsign=false` to get past a failure; fix the cause or ask.
- **Rebase conflicts:** resolve, `git add`, `git rebase --continue`; if unsure, `git rebase --abort` and ask. Never leave "rebase in progress".
- **Safety net:** `git branch backup/<name>` before any rewrite; `git reflog` recovers lost commits.

Read `references/history-editing.md` when splitting, reordering, squashing or rewording several commits, or recovering from a bad rebase.
Read `references/signing.md` when commits must be signed or a signature shows as unverified.
Read `references/secret-committed.md` immediately when a secret has been committed or pushed.

## Safety gates
- Ask before: pushing, any force-push (even with lease), rewriting a commit that exists on a remote, deleting branches, `git reset --hard`, `git clean -fd`, `git stash drop`.
- Secret in a pushed commit: stop, tell the user, rotate first. Rewriting history alone does not fix it.
- If the repo or user wants review before commit, show files, branch and full message, then wait.

## Output
```
Branch: <branch>   Base: <base>
Commits:
  <sha7> <type(scope)!: summary>   (<n> files, +<a>/-<d>)
Identity: <name> <email>   Signed: yes | no | not required
Checks: <hooks/lint/tests> -> pass | fail (<detail>)
Secret scan: clean | <finding and action taken>
Pushed: no (awaiting approval) | yes -> <remote>/<branch>
```
