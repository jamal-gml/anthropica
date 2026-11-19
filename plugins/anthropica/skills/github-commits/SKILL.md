---
name: github-commits
description: Write and organise Git commits to industry standard (Conventional Commits, atomic changes, clear messages, correct identity, no secrets). Use when staging, committing, amending, splitting or rewording commits, or when cleaning up history before a push.
---

# GitHub Commits

## Before committing
1. Run `git status` and `git diff --staged`; commit only what belongs to this change.
2. Read the repository's own rules first (`CONTRIBUTING.md`, `CLAUDE.md`, commitlint or husky config). Repository rules override this skill.
3. Confirm the author identity: `git config user.name` and `user.email` must be the user's. Set them per repository if needed; never invent an identity.
4. Never commit secrets, `.env` files, credentials, private keys, large binaries or build output. If one is staged, unstage it and add it to `.gitignore`.
5. Do not commit directly to a protected or default branch (`main`, `master`, `develop`, release branches) unless the user explicitly asks. Create a branch instead (see `github-pull-requests`).

## One commit, one logical change
- Each commit should build and pass tests on its own and be revertable without collateral damage.
- Split unrelated changes with `git add -p`. Keep formatting-only or rename-only changes in their own commit.
- Do not mix dependency bumps, refactors and features in one commit.

## Message format (Conventional Commits 1.0)
```
<type>(<scope>)<!>: <imperative summary, 72 chars max, no full stop>

<body: why the change is needed and what it affects, wrapped at 72>

<footers: BREAKING CHANGE: ..., Refs: #123, Fixes ABC-42>
```
- Types: `feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `build`, `ci`, `chore`, `style`, `revert`.
- Release effect: `feat` is minor, `fix`/`perf` is patch, `!` or a `BREAKING CHANGE:` footer is major. Choose the type by its release effect.
- Imperative summary ("add", not "added"). Scope is the module or area, lowercase, optional.
- The body explains *why*; the diff already shows *what*. Omit it only for trivial changes.
- Reference issues in footers with the tracker's keyword (`Fixes #12`, `Closes ABC-42`).
- Never use "WIP", "fix stuff", "update file", emoji-only or ticket-only messages.

## History hygiene
- Amend, squash or reword only commits that have not been pushed to a shared branch.
- Never force-push a shared branch. On your own branch use `git push --force-with-lease`, never plain `--force`.
- Never skip hooks (`--no-verify`) or disable signing to get a commit through; fix the cause.
- If the repository requires signed commits, sign them with the user's configured GPG or SSH key.

## Before pushing
- `git log origin/<base>..HEAD --format='%h %an <%ae> %s'` to check authors and messages.
- Re-scan the outgoing diff for secrets: `git diff origin/<base>..HEAD`.

## Output
When the user wants to review first, show the files, branch and full message and wait for approval. After committing, report the hash and summary.
