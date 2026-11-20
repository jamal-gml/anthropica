# Debugging workflow runs

## Read the failure
```bash
gh run list --workflow <file> -L 10              # pattern: always failing, or since when?
gh run view <run-id>                             # jobs and their conclusions
gh run view <run-id> --log-failed                # failed steps only
gh run view <run-id> --job <job-id> --log        # full job log
gh run download <run-id>                         # artifacts (reports, screenshots)
gh api repos/{owner}/{repo}/actions/runs/<run-id>/jobs -q '.jobs[] | {name, conclusion, runner_name}'
```
1. Find the first error, not the final `exit code 1`.
2. Same failure on the default branch? Then it is pre-existing or environmental.
3. Diff what changed since the last green run: `gh run view <green-id> --json headSha` then `git diff <sha>..<sha>`, including action versions and runner image (`Runner Image` in "Set up job").

## Symptom table
| Symptom | Likely cause | Fix |
|---|---|---|
| Workflow did not trigger | wrong branch/paths filter, file not on default branch (for `schedule`, `workflow_dispatch`, `issue_comment`), YAML error, triggered by `GITHUB_TOKEN` | `actionlint`; check filters; check the Actions tab for "workflow file issue" |
| Job skipped | `if:` false, `needs` job failed/skipped (dependent jobs skip unless `if: always()`/`!cancelled()`) | inspect `needs.<job>.result` |
| Required check stuck "Expected" | path-filtered or renamed workflow never reported | aggregator job, see `patterns.md` |
| `Resource not accessible by integration` | missing `permissions:` scope, or fork PR read-only token | add the specific scope to the job; forks by design |
| `Permission denied (publickey)` / 403 on push | `persist-credentials: false` with a later push, or token lacks `contents: write`, or branch protection | grant scope or use an app token; never bypass protection |
| OIDC `AssumeRoleWithWebIdentity` / `AADSTS70021` | `sub` mismatch, missing `id-token: write` | see `oidc.md` claim debugging |
| Secret is empty | secret in environment but job has no `environment:`, fork PR, name typo, Dependabot runs use Dependabot secrets | check scope |
| Works on ubuntu, fails on windows | shell default is `pwsh`; path separators; line endings | set `shell: bash`; `.gitattributes` |
| Cache never hits | key includes volatile data; cache saved only on success; different branch scope | stable key on lockfile hash; `restore-keys` |
| Slow job | no cache, serial jobs, full clone, oversized matrix | see `patterns.md` Speed |
| Hangs until timeout | interactive prompt (apt, npm login, git credential), service not ready, test deadlock | non-interactive flags, health-check services, `timeout-minutes` |
| `No space left on device` | large images/build outputs | prune, larger runner, clean tool caches |
| Rate limit / 429 | many API calls with `GITHUB_TOKEN` or registry pulls | cache, authenticate pulls, back off |
| Action broke without changes | unpinned action tag moved, or runner image updated | pin to SHA; pin runner (`ubuntu-24.04` not `-latest`) when stability matters |

## More visibility
- Re-run with debug: `gh run rerun <id> --debug` (enables step and runner debug logs) - only after reading normal logs.
- Add a temporary diagnostic step: `env | sort | grep -v -i -E 'token|secret|key|pass'`, `df -h`, tool `--version`. Remove before merging.
- Context dumps: `echo '${{ toJSON(github.event) }}'` is unsafe for untrusted events (injection); write it via `env:` to a file instead.

## Reproduce locally
- Run the same commands in the same container image (`docker run -it <image>`).
- `act` can emulate many workflows but differs from hosted runners (no OIDC, different images); do not treat an `act` pass as proof.

## Fix rules
- Fix the cause in a normal PR. Never add `continue-on-error`, delete required checks, weaken `permissions:` or remove environment gates to get green.
- Re-run once for a suspected flake (`gh run rerun <id> --failed`); repeated failure is real.
- Ask before re-running deploy jobs.
