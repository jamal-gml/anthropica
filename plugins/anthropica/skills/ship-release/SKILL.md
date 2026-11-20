---
name: ship-release
description: Version, tag, publish and promote releases (semver, changelogs, annotated tags, release-please/changesets/semantic-release, registries, hotfixes, rollback, yanking a bad version). Use when cutting or tagging a release, publishing a package or image, writing release notes, promoting a version, shipping a hotfix or withdrawing a release. Not for go/no-go (use `ship-readiness`).
---

# Release Management

Turn a verified commit into an immutable, documented, traceable version, and be able to withdraw it.

## Procedure
1. **Readiness:** for anything reaching users, run `ship-readiness` first and have a GO.
2. **Find the release mechanism:** look for `release-please-config.json`, `.release-please-manifest.json`, `.changeset/`, `.releaserc*`/`release.config.*`, `goreleaser.yml`, a release workflow, `RELEASING.md`. If automation exists, use it; never hand-tag around it.
3. **Determine the version:**
   - Last tag: `git describe --tags --abbrev=0` (or `git tag --sort=-v:refname | head`).
   - Changes: `git log <last-tag>..origin/<default> --format='%h %s'`.
   - Apply the bump table below; the highest change wins.
4. **Update version sources** (manifest, lockfile if it embeds version, `version.txt`, chart `appVersion`) in one release PR, plus the changelog. Merge via `review-pull-request`.
5. **Confirm the commit:** `git fetch origin && git rev-parse origin/<default>`; CI must be green on that exact SHA (`gh run list --commit <sha>`).
6. **Ask, then tag** the merge commit (annotated, signed if required):
   `git tag -a vX.Y.Z <sha> -m "vX.Y.Z"` then `git push origin vX.Y.Z`. Never `git push --tags` (pushes stray local tags).
7. **Build and publish from CI on the tag,** not from a laptop: artifacts, checksums, SBOM and provenance (`secure-supply-chain`).
8. **Release notes:** `gh release create vX.Y.Z --verify-tag --title "vX.Y.Z" --notes-file notes.md` (or `--generate-notes` then edit). Pre-releases: `--prerelease`.
9. **Promote** the same artifact through environments (gates in `ship-ci-cd`); verify after each step.
10. **Record** version, SHA, artifacts and environments in the release, tracker and deploy log.

## Version bump rules (SemVer 2.0)
| Change since last release | 1.x and later | 0.x |
|---|---|---|
| Breaking API/CLI/config/schema/behaviour (`!`, `BREAKING CHANGE:`) | MAJOR | MINOR |
| New feature, new optional parameter, deprecation notice (`feat`) | MINOR | MINOR or PATCH |
| Bug fix, perf, security fix without API change (`fix`, `perf`) | PATCH | PATCH |
| Docs, tests, CI, refactor only | no release (or PATCH if shipped artifact changed) | same |

- Removing or renaming anything public, raising minimum runtime/platform versions, and changing defaults users rely on are breaking.
- Dependency bumps: PATCH unless they change your public surface.
- Pre-releases: `X.Y.Z-rc.N`, `-beta.N`; they sort before `X.Y.Z`. Build metadata `+sha` is ignored for precedence.
- Tag format `vX.Y.Z` unless the ecosystem/repo uses another; monorepos `<package>@X.Y.Z` or `<package>-vX.Y.Z`.

## Changelog and release notes
- Keep a Changelog sections: Added, Changed, Deprecated, Removed, Fixed, Security. Newest first, dated `## [X.Y.Z] - YYYY-MM-DD`, `## [Unreleased]` on top.
- Order notes by user impact: **Breaking changes with migration steps first**, then security fixes (with CVE/GHSA IDs), features, fixes.
- Write for consumers ("Export now supports CSV"), link PRs/issues, credit contributors. Exclude `chore`, `ci`, `test`, internal refactors.
- Generated notes are a draft; edit for clarity.

## Rules
- **Immutable:** never move, delete or reuse a published tag or version. A bad release gets a new version.
- **Tag the exact SHA that passed CI** and is on the default (or release) branch: `git branch -r --contains <sha>`.
- **One source of truth** for the version; CI checks the tag equals the manifest version.
- **Images:** tag `X.Y.Z`, `X.Y`, and the commit SHA; deploy by digest. `latest` only as a convenience, never as a deploy reference.
- **Credentials:** publish with OIDC trusted publishing or short-lived scoped tokens in CI; never from personal tokens on a laptop.
- **Automation bots** need tokens that trigger CI on their release PRs (app token, not the default CI token).
- **Release branches** (`release/X.Y`) only when supporting several versions; fixes land on default first, then cherry-pick (`git cherry-pick -x <sha>`).

Read `references/automation-tools.md` when choosing or configuring release-please, changesets, semantic-release or goreleaser.
Read `references/registries.md` when publishing to npm, PyPI, Maven, NuGet, crates.io, Go modules, Helm or container registries.
Read `references/hotfix-and-withdrawal.md` when shipping a hotfix, rolling back, or yanking/deprecating a bad release.

## Safety gates
- Ask before: pushing tags, creating or editing GitHub/GitLab releases, publishing to any registry, promoting to production, yanking/deprecating versions, deleting anything.
- Production promotion requires the user's explicit "yes" in the current conversation after you state version, SHA, changes and rollback plan.
- If publish partly fails (some packages/platforms published), stop and report exactly what is live; do not retry blindly or republish the same version.

## Output
```
Release: vX.Y.Z  (previous vA.B.C, bump: major|minor|patch because <reason>)
Commit: <sha> on <branch>, CI <green run url>
Tag: vX.Y.Z annotated, signed: yes|no
Artifacts: <name -> location, checksum/digest>   Provenance/SBOM: <yes/no>
Published to: <registry: package@version>
Notes: <url>  Breaking: <none | summary + migration>
Environments: staging <deployed/verified>, production <pending approval | deployed at time>
Rollback: <command or previous version/digest>
```

---

**Phase:** Ship. **Next:** when this step is done, offer Operate: `/anthropica:operate-incident`, `/anthropica:operate-observability`. Phase map: `/anthropica:workflow`.
