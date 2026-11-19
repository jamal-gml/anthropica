# Anthropica maintainer guide

This repository publishes role-based Claude skill packs as a Claude Code plugin marketplace and as Claude.ai skill ZIPs.

## Source of truth
- Edit only `skills/<name>/SKILL.md`, `agents/<name>.md`, `packs.json`, `scripts/`, `global/` and the hand-written parts of `README.md`.
- Never hand-edit `plugins/`, `.claude-plugin/marketplace.json`, `dist/` or the README table between `<!-- packs:start -->` and `<!-- packs:end -->`. They are generated.
- After any source change run `python scripts/build.py` and commit the generated output in the same commit. CI fails if it is stale.

## Adding or changing a skill
1. Create `skills/<name>/SKILL.md`. The name is lowercase-hyphenated and equals the folder name.
2. Frontmatter: `name` and `description`. The description says what the skill does **and** when to use it ("Use when ..."), in under 1024 characters. The build rejects descriptions without a trigger.
3. Body: short, imperative, checklist-style; under about 150 lines. Put the method, then the safety gates, then the expected output. Industry-standard practice, no organisation-specific names or secrets.
4. Skills must tell Claude to ask before destructive or outward-facing actions (deploys, force pushes, production changes, publishing).
5. Add the skill to every relevant pack in `packs.json`. The full `anthropica` pack picks it up automatically. The build prints a note for skills that are only in the full pack.
6. Run the build, then `claude plugin validate .` and `claude plugin validate plugins/<pack>` if the Claude CLI is available.

## Adding a role pack
Add an entry to `packs.json` → `packs` with `name` (`anthropica-<role>`), `title`, `role`, `description`, `category`, `skills` and `agents`. Core skills (`core`) are added to every pack automatically. Rebuild.

## CLI (`cli/`)
- Go module `github.com/jamal-gml/anthropica/cli`. `main.go` only calls `cli.Execute()`; each command is one file in `internal/cli/`, registered in `root.go`.
- Commands write to `cmd.OutOrStdout()`/`cmd.ErrOrStderr()` and are tested in-process in `cli_test.go`. Tests never call the real GitHub API (`githubAPIBase` points at a test server).
- Before committing CLI changes, from `cli/`: `go fmt ./...`, `go mod tidy`, `go vet ./...`, `go test ./...`, `go build -o anthropica.exe .`. CI runs the same on Linux and Windows.
- Keep it provider-neutral and safe: `init` must never overwrite files, the GitHub commands stay read-only, and tokens are never printed.
- Its version comes from the release tag via `-ldflags`; there is no version constant to bump.

## Versioning
- One version for everything, in `packs.json` → `version` (semver). The build copies it into every plugin and the marketplace.
- Bump MINOR for new skills or packs, PATCH for wording fixes, MAJOR for removing or renaming a skill or pack (that breaks users' commands).
- Claude Code only offers users an update when the version changes, so bump it in every release.

## Releasing
Full checklist and recovery steps: `RELEASING.md`.

1. Work on a branch (`feature/`, `fix/`, `docs/`, `chore/`), open a PR into `main`, and merge once CI is green. Do not push directly to `main`.
2. In a release PR, bump `packs.json` version, rebuild, and merge.
3. Tag the merge commit on `main`: `git tag -a vX.Y.Z -m "vX.Y.Z" && git push origin vX.Y.Z`. The tag must equal `v` + the `packs.json` version.
4. The `release` workflow verifies the tag, rebuilds, tests and cross-compiles the CLI, and publishes a GitHub Release with each pack's `<pack>-claude-ai.zip`, every single-skill ZIP, and CLI archives for Windows, macOS and Linux (`anthropica_<os>_<arch>`, no version in the name) with `checksums.txt`, then tags `cli/vX.Y.Z` for `go install`.
5. Users of the Claude Code marketplace get the update with `/plugin marketplace update anthropica`.

Ask the repository owner before tagging, publishing a release, or changing repository settings.

## Commits and PRs
Follow `skills/github-commits` and `skills/github-pull-requests`: Conventional Commits (`feat(skills): add seo-audit`), small focused PRs, squash merge into `main`.
