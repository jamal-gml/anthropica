# Hotfixes, rollback and withdrawing a release

## Decide first
| Situation | Action |
|---|---|
| Users impacted now, previous version compatible with current data | **Roll back** the deployment to the previous artifact/digest or flip the flag, then fix |
| Rollback unsafe (irreversible migration, external contract already changed) | **Roll forward** with a minimal hotfix |
| Defect in a published library/package | publish a fixed PATCH; then deprecate/yank the bad version |
| Security vulnerability in a published version | coordinated fix release + advisory (GHSA/CVE); yank only if exploitation risk outweighs breakage |
| Bad release not yet promoted to production | stop promotion; mark release as not for use; fix forward |

Rollback first, investigate second, when users are affected (`operate-incident`).

## Rollback mechanics
- Redeploy previous version by digest/artifact ID through the normal pipeline (fastest safe path), e.g. GitOps revert of the version bump, `kubectl rollout undo`, re-run the deploy job for the previous release.
- Do not rebuild old source for a rollback; redeploy the old artifact.
- Database: app rollback must not need schema rollback (expand/contract, `ship-ci-cd`). If it does, stop and escalate.
- Verify: health checks, error rate back to baseline, key flows.
- Record the rollback in the deploy log and the incident.

## Hotfix flow
**Trunk-based (default branch is releasable)**
1. Branch `hotfix/<key>-<slug>` from the default branch, fix with a regression test.
2. Fast-tracked PR: CI must pass and one approval is still required.
3. Release next PATCH from the default branch through the normal (expedited) promotion.

**Default branch has unreleased work that must not ship**
1. Branch from the production tag: `git switch -c hotfix/<key>-<slug> vX.Y.Z`, or from `release/X.Y`.
2. Fix with a regression test; tag `vX.Y.(Z+1)` from that branch after CI.
3. Forward-port: cherry-pick (`git cherry-pick -x <sha>`) or merge the fix into the default branch in the same day, so the next release does not regress.
4. Automation tools: use their release-branch support (release-please `target-branch`, changesets on a maintenance branch, semantic-release maintenance branches `X.Y.x`).

Gates are shortened, never removed: CI green, one reviewer, production approval with version and SHA.

## Withdrawing a published version
Never delete-and-reuse a version number. Mark it bad and publish a fix.

| Registry | Command | Effect |
|---|---|---|
| npm | `npm deprecate <pkg>@<version> "Broken: use X.Y.Z"` | install warns; still installable |
| npm | `npm unpublish <pkg>@<version>` | only within 72 h and limited conditions; avoid - breaks lockfiles; version can never be reused |
| npm dist-tag | `npm dist-tag add <pkg>@<good> latest` | moves `latest` back to a good version |
| PyPI | yank in the web UI (or API) with reason | excluded from resolution unless pinned exactly |
| crates.io | `cargo yank --version X.Y.Z <crate>` | new lockfiles avoid it; existing keep working |
| Go | add `retract vX.Y.Z // reason` to `go.mod`, release a new version | `go get` avoids retracted versions |
| NuGet | `dotnet nuget delete` (unlists) | hidden from search; restorable by exact version |
| RubyGems | `gem yank <gem> -v X.Y.Z` | removed from index |
| Maven Central | cannot remove | publish fix; announce |
| Container image | do not delete or retag the digest; move floating tags (`X.Y`, `latest`) back to the good digest; publish fix | pinned digests keep working |
| GitHub Release | edit notes: "Withdrawn: <reason>, use vX.Y.Z"; mark as pre-release so it is not "Latest" | tag stays immutable |

Ask before any of these; they affect every consumer.

## Communicate
- Release notes of the bad version: what is wrong, who is affected, which version to use.
- Advisory for security issues, with affected and fixed ranges.
- Notify downstream teams and support; update the incident or tracker.
