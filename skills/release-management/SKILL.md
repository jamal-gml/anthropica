---
name: release-management
description: Version, tag, publish and promote releases and packages (semantic versioning, changelogs, GitHub Releases, package registries, environment promotion, rollback). Use when cutting a release, publishing a package or container image, promoting a build between environments, or rolling back.
---

# Release Management

Run `release-readiness` first for anything going to production. A release that reaches users requires the user's explicit approval in the current conversation, after stating the version, commit SHA and changes.

## Versioning
- Semantic Versioning: MAJOR for breaking changes, MINOR for features, PATCH for fixes. Pre-releases: `1.4.0-rc.1`.
- Derive the bump from Conventional Commits since the last tag (`git log <last-tag>..HEAD --oneline`).
- One source of truth for the version (package manifest, `version.txt`, or the tag). Keep all copies in sync.

## Changelog
- Group by Added, Changed, Fixed, Removed, Security (Keep a Changelog) or by commit type.
- Write for consumers: what changed and what they must do. Call out breaking changes and migrations first.
- Prefer automation already in the repo (release-please, semantic-release, changesets) over manual edits.

## Tag and GitHub Release
```bash
git switch main && git pull --ff-only
git tag -a v1.4.0 -m "v1.4.0"        # signed (-s) if the repo requires it
git push origin v1.4.0
gh release create v1.4.0 --title "v1.4.0" --notes-file CHANGELOG-excerpt.md dist/*
```
- Tag the exact commit that passed CI. Never move or reuse a published tag.
- Build release artifacts in CI from the tag, not on a laptop. Attach checksums, and an SBOM or provenance where supported.

## Packages and images
- Publish from CI with scoped, short-lived credentials (OIDC trusted publishing for npm and PyPI where available).
- Container images: tag with the version and the commit SHA, never only `latest`. Push to the agreed registry.
- Never republish an existing version. Fix forward with a new patch version.

## Promotion through environments
1. Deploy the same immutable artifact to each environment (build once, promote many).
2. Order: development → staging → production, each gated by checks and, for production, an environment approval.
3. Verify after each step: health checks, error rates, key user flows.
4. Record what was deployed where (deployment record, release notes, ticket update).

## Rollback
- Know the rollback path before deploying: previous artifact version, feature flag, or revert PR.
- Database migrations must be backwards compatible with the previous release, or have a tested down path.
- Roll back first, investigate second, when users are affected.

## Output
Version, commit SHA, changelog summary, artifacts and their locations, environments deployed, verification results, and the rollback command.
