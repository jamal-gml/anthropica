# A secret was committed

Assume any secret that reached a remote is compromised. Forks, clones, CI logs, caches and scanners may already hold it. Rotation is the fix; history rewriting is cleanup.

## Decide by where it is
| Where the secret is | Action |
|---|---|
| Staged only | `git restore --staged <file>`, move the value to env/secret store, add path to `.gitignore`. No rotation needed if it never left the machine. |
| Committed, not pushed | Remove it, then rewrite the unpushed commit (`--amend`, or fixup + autosquash); confirm with `git log -p -S'<distinctive fragment>' --all`. Rotate anyway if the working copy is synced, backed up or mirrored automatically. |
| Pushed to any remote (even a feature branch, even for seconds) | **Stop. Tell the user. Rotate first.** Then clean up. |
| In a public repo or a merged default branch | Rotate immediately; treat as an incident (`operate-incident`). |

## Pushed: order matters
1. **Tell the user** what leaked (type, file, commit, branch), without printing the value.
2. **Rotate/revoke** at the issuer: create a new credential, deploy it, then revoke the old one. Revoke first only if active abuse is likely and an outage is acceptable. The user or credential owner does this; do not do it without approval.
3. **Check for use:** issuer audit logs (cloud audit trail, API provider dashboard) for activity since the commit time.
4. **Remove from the code** in a normal commit: read the value from env or a secret manager; add ignore rules and a pre-commit secret scanner.
5. **Rewrite history only if the user and repo admins agree** (it disrupts every clone and open PR):
   - `git filter-repo --replace-text expressions.txt` or `git filter-repo --path <file> --invert-paths` (preferred over `filter-branch`; BFG is an alternative).
   - Requires a force-push to protected branches: admin action, coordinated; everyone re-clones.
   - Ask the hosting provider to purge cached views and dangling commits/PR refs if needed.
6. **Record** what was rotated and when, for the incident or security log.

## Do not
- Do not paste the secret into chats, issues, PR comments or commit messages ("remove key AKIA...").
- Do not consider `git revert` or a deletion commit a fix: the value stays in history.
- Do not skip rotation because "the branch was private" or "it was deleted quickly".
- Do not force-push shared branches on your own initiative.

## Prevent recurrence
- Pre-commit: `gitleaks`, `detect-secrets` or `trufflehog` via the repo's hook framework.
- Server side: enable push protection / secret scanning (`secure-repo-governance`).
- Keep `.env.example` with placeholders; load real values from a secret manager or CI secrets (OIDC where possible, see `ship-ci-cd`).
