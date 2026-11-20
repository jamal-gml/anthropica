# Release automation tools

Use what the repo already has. Introduce a tool only with the user's agreement.

## Comparison
| Tool | Model | Version source | Best for | Watch out |
|---|---|---|---|---|
| release-please | bot keeps a release PR open; merging it tags and releases | Conventional Commits since last release | repos with squash-merged Conventional Commit PR titles; multi-language monorepos (manifest mode) | needs accurate commit types; its PR must trigger CI (use an app token); do not edit its branch by hand |
| changesets | contributors add `.changeset/*.md` files declaring bump per package; bot opens "Version Packages" PR | changeset files | JS/TS monorepos, explicit human-chosen bumps and notes | contributors forget changesets (add a CI check: `changeset status --since=origin/main`) |
| semantic-release | fully automatic on every push to release branch: analyse commits, tag, publish | Conventional Commits | single packages wanting continuous release | releases on every qualifying merge; harder to batch; needs `fetch-depth: 0` |
| goreleaser | builds and publishes binaries, archives, checksums, images, SBOMs from a tag | git tag | Go (and other) CLI/binary distribution | pairs with a versioning tool or manual tag |
| Manual (script + checklist) | human bumps and tags | manifest | small repos, low cadence | drift between manifest and tag; add a CI check |

## release-please essentials
- Files: `release-please-config.json` (packages, release-type, changelog sections), `.release-please-manifest.json` (current versions).
- Bootstrapping existing projects: set the manifest to the current version; optionally `bootstrap-sha` so history before it is ignored.
- Force a specific version: commit with footer `Release-As: 2.0.0`.
- Pre-1.0: `bump-minor-pre-major: true` makes breaking changes bump minor.
- Never push to `release-please--*` branches; change the config or commits instead.

## changesets essentials
- `npx changeset` (interactive) or write `.changeset/<name>.md` with frontmatter `"pkg": minor` and a summary.
- `changeset version` applies bumps and changelogs; `changeset publish` publishes and tags.
- Pre-release mode: `changeset pre enter rc` / `pre exit`; remember to exit.
- Linked/fixed groups in `.changeset/config.json` for packages that must share versions.

## semantic-release essentials
- Config in `.releaserc` / `release.config.js`: `branches` (`main`, prerelease channels like `{ name: 'beta', prerelease: true }`), plugins (commit-analyzer, release-notes-generator, changelog, npm, github, git).
- Dry run: `npx semantic-release --dry-run` (needs tokens to verify; do not publish).
- Commits like `chore(release)` committed back to the branch need permission to push to a protected branch: prefer not committing back (changelog in release notes only).

## Common to all
- Commit/PR title lint (commitlint) in CI keeps automation honest.
- Squash-merge with Conventional Commit PR titles produces one clean entry per PR.
- Tag must match package versions; add a CI check.
- The release job runs on tag or bot PR merge, with least privilege: `contents: write` only in the job that tags, `id-token: write` only in publish jobs.
- Test the flow with a pre-release before relying on it for a stable release.

## Migrating between tools
1. Freeze releases; note the last version per package.
2. Configure the new tool with that version as its baseline (manifest/bootstrap).
3. Remove the old tool's config and workflow in the same PR.
4. Cut a pre-release to prove the pipeline; then resume.
