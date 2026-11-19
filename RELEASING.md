# Releasing and deploying Anthropica

One tag publishes everything: the Claude Code marketplace update, Claude.ai skill packs and the `anthropica` CLI binaries. The work is done by [`.github/workflows/release.yml`](.github/workflows/release.yml).

## What a release produces

| Output | Where | Used by |
|---|---|---|
| Updated plugins and marketplace | `main` branch (`plugins/`, `.claude-plugin/marketplace.json`) | Claude Code: `/plugin marketplace update anthropica` |
| `<pack>-claude-ai.zip` (one per role pack) | GitHub Release assets | Claude.ai / Claude Desktop skill uploads |
| `<skill>.zip` (one per skill) | GitHub Release assets | Uploading a single skill |
| `anthropica_<os>_<arch>.zip` / `.tar.gz` | GitHub Release assets | CLI downloads (Windows, macOS, Linux; amd64 and arm64) |
| `checksums.txt` | GitHub Release assets | SHA-256 verification of the CLI archives |
| Tag `cli/vX.Y.Z` | Git tag, created by the workflow | `go install github.com/jamal-gml/anthropica/cli@latest` |

Asset names carry no version, so `https://github.com/jamal-gml/anthropica/releases/latest/download/<file>` always points at the newest release.

## Before you release

- [ ] Everything for the release is merged into `main` through pull requests, and CI on `main` is green.
- [ ] `packs.json` → `version` is bumped (semver: MINOR for new skills, packs or CLI commands; PATCH for fixes; MAJOR for removing or renaming a skill, pack or command).
- [ ] `python scripts/build.py` was run after the bump and its output committed (CI enforces this).
- [ ] From `cli/`: `go vet ./...` and `go test ./...` pass.

## Cut the release

The tag must be exactly `v` + the `packs.json` version, on the commit on `main` that passed CI.

**CMD or PowerShell:**
```bat
git switch main
git pull --ff-only
python -c "import json; print(json.load(open('packs.json'))['version'])"
git tag -a v1.2.0 -m "v1.2.0"
git push origin v1.2.0
```

Then watch it:
```bat
gh run watch --repo jamal-gml/anthropica
gh release view v1.2.0 --repo jamal-gml/anthropica
```

The workflow:
1. Fails unless the tag matches `packs.json`.
2. Rebuilds the packs and fails if the committed generated files are stale.
3. Tests the CLI and cross-compiles it for six platforms with the version embedded.
4. Publishes the GitHub Release with generated notes and every asset above.
5. Creates the `cli/vX.Y.Z` tag for `go install`.

## After the release

- Smoke-test a download (see [Install the CLI](README.md#install-the-cli)) and `anthropica version`.
- Claude Code users update with `/plugin marketplace update anthropica`, then `/reload-plugins`.
- Announce changes to teams, calling out new or removed skills per pack.

## If something goes wrong

- **Workflow failed before publishing:** fix on a branch, merge, delete the tag (`git push origin :refs/tags/v1.2.0` and `git tag -d v1.2.0`), and tag the new commit.
- **Release published with a bad asset:** do not replace assets or move the tag. Fix forward with a new PATCH version.
- **Bad plugin content on `main`:** revert the change in a PR and release a new PATCH version. Users get it on their next marketplace update.

Never move or reuse a published tag, and never publish from a local machine; releases are built only by the workflow.
