# Anthropica maintainer guide

This repository publishes role-based Claude skill packs as a Claude Code plugin marketplace and as Claude.ai skill ZIPs.

## Source of truth
- Edit only `skills/<name>/SKILL.md`, `agents/<name>.md`, `packs.json`, `scripts/`, `global/` and the hand-written parts of `README.md`.
- Never hand-edit `plugins/`, `.claude-plugin/marketplace.json`, `dist/` or the README table between `<!-- packs:start -->` and `<!-- packs:end -->`. They are generated.
- After any source change run `python scripts/build.py` and commit the generated output in the same commit. CI fails if it is stale.

## Adding or changing a skill
1. Create `skills/<phase>-<topic>/SKILL.md`, lowercase-hyphenated, folder name = `name`. `<phase>` is one of `packs.json` → `phases` (plan, build, create, test, review, secure, ship, operate, core); the build rejects anything else. `core-*` skills are the `core` list added to every pack. `workflow` is reserved: the build generates one per pack (phase map + next-skill advice) and appends a "Phase / Next" line to every skill, so never write either by hand.
2. Frontmatter: `name` and a **single-line** `description` (the build parser reads one line), ideally under 400 characters: what it does, then "Use when ..." with the concrete triggers (task words, file names, errors). Add "Not for X (use `<other-skill>`)" where skills overlap.
3. Body, in this order: `# Title`, one-line purpose, `## Procedure` (numbered steps), `## Rules` (decision rules, thresholds, gotchas), `## Safety gates`, `## Output` (a template). Target 60-150 lines.
4. Write only what Claude would not do unprompted: exact commands, thresholds, decision tables, failure modes, ordering that matters. Cut generic advice ("write clean code", "follow best practices"). Imperative, terse, no filler, tables for mappings.
5. Long or situational material (full checklists, examples, per-tool detail) goes in `skills/<name>/references/<topic>.md`, linked from SKILL.md with when to read it: "Read `references/terraform.md` when the repo uses Terraform." Only the description is always in context; the body loads on use; references load only when read.
6. Don't repeat another skill: refer to it by name (`ship-github-actions`). Industry-standard, vendor-neutral practice, no organisation-specific names or secrets.
7. Skills must tell Claude to ask before destructive or outward-facing actions (deploys, force pushes, production changes, publishing, sending).
8. Add the skill to every relevant pack in `packs.json`. The full `anthropica` pack picks it up automatically. The build prints a note for skills that are only in the full pack.
9. Run the build, then `claude plugin validate .` and `claude plugin validate plugins/<pack>` if the Claude CLI is available.

Renaming a skill breaks users' commands: bump MAJOR and add the old and new names to the migration table in `README.md`.

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
4. The `release` workflow verifies the tag, rebuilds, tests and cross-compiles the CLI, and publishes a GitHub Release with each pack's `<pack>-claude-ai.zip`, every single-skill ZIP, and CLI archives for Windows, macOS and Linux (`anthropica_<os>_<arch>`, no version in the name) with `checksums.txt`, tags `cli/vX.Y.Z` for `go install`, and publishes the npm packages (`anthropica` launcher in `npm/anthropica/`, assembled with the platform packages by `scripts/npm/build.mjs`) through npm trusted publishing.
5. Users of the Claude Code marketplace get the update with `/plugin marketplace update anthropica`.

Ask the repository owner before tagging, publishing a release, or changing repository settings.

## Commits and PRs
Follow `skills/build-commits` and `skills/review-pull-request`: Conventional Commits (`feat(skills): add review-seo`), small focused PRs, squash merge into `main`.
