# Triage failing PR checks

## Commands
```bash
gh pr checks <n>                                   # state of every check
gh pr checks <n> --required                        # only the ones that gate merge
gh run list --branch <head-branch> -L 5            # recent runs for the branch
gh run view <run-id> --log-failed                  # logs of failed steps only
gh run view <run-id> --job <job-id> --log          # full log of one job
gh run list --branch <base> --workflow <file> -L 5 # is the base branch green?
gh run rerun <run-id> --failed                     # re-run only failed jobs (once, for flakes)
gh run download <run-id> -n <artifact>             # test reports, screenshots
```

## Decision order
1. **Required or optional?** Optional failures still need a note in the PR, but do not block.
2. **Fails on base too?** Then it is pre-existing: report it, link the base run, do not fix it in this PR unless asked.
3. **First real error.** Scroll up from the end; the last line is usually "Process completed with exit code 1".
4. **Reproduce locally** with the same command and versions the job uses (read the workflow file).
5. **Fix in a new commit** with the right type (`fix`, `test`, `ci`). Push only after the user agrees.

## Failure patterns
| Signal | Likely cause | Action |
|---|---|---|
| Passes locally, fails in CI | version drift (runtime, tool, lockfile), env var missing, case-sensitive filesystem, timezone/locale | pin versions to match CI; use the CI image locally |
| Lockfile / "frozen" install fails | manifest changed without regenerating lockfile | run the package manager's install, commit the lockfile |
| Format or lint only | hook not run locally | run the formatter; commit as `style:` or fold into the commit |
| Fails intermittently, different tests each time | flaky test, shared state, ordering, real time/network | re-run once; if it recurs, report as flaky with evidence; never mark the check optional |
| Timeout | hung test, missing service, deadlock, slow runner | find the last log line before hang; add timeouts |
| `Resource not accessible by integration` / 403 | token lacks permission; fork PR has read-only token | expected for forks; do not add secrets to fork runs (`ship-github-actions`) |
| Secret empty in fork PR | secrets are not passed to fork workflows | by design; maintainers run privileged checks after review |
| Check "Expected - waiting for status" forever | required check name no longer produced (job renamed, path filter skipped it) | report to maintainers; never rename a required check without updating rules |
| Coverage threshold failed | new code untested | add tests rather than lowering the threshold |
| Dependency audit failed | new advisory, possibly unrelated to this PR | check if base fails too; bumps go through the repo's dependency bot where possible |
| Merge conflict / "branch out of date" | base moved | update the branch (see `branch-maintenance.md`) |

## Never
- Add `continue-on-error`, `|| true`, `skip` markers or lower thresholds to turn a check green.
- Retry more than once without investigating.
- Ask an admin to bypass a required check for a non-emergency.
