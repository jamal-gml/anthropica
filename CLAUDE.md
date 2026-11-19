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

## Versioning
- One version for everything, in `packs.json` → `version` (semver). The build copies it into every plugin and the marketplace.
- Bump MINOR for new skills or packs, PATCH for wording fixes, MAJOR for removing or renaming a skill or pack (that breaks users' commands).
- Claude Code only offers users an update when the version changes, so bump it in every release.

## Releasing
1. Work on a branch (`feature/`, `fix/`, `docs/`, `chore/`), open a PR into `main`, and merge once CI is green. Do not push directly to `main`.
2. In a release PR, bump `packs.json` version, rebuild, and merge.
3. Tag the merge commit on `main`: `git tag -a vX.Y.Z -m "vX.Y.Z" && git push origin vX.Y.Z`. The tag must equal `v` + the `packs.json` version.
4. The `release` workflow verifies the tag, rebuilds, and publishes a GitHub Release with each pack's `<pack>-claude-ai.zip` and every single-skill ZIP attached.
5. Users of the Claude Code marketplace get the update with `/plugin marketplace update anthropica`.

Ask the repository owner before tagging, publishing a release, or changing repository settings.

## Commits and PRs
Follow `skills/github-commits` and `skills/github-pull-requests`: Conventional Commits (`feat(skills): add seo-audit`), small focused PRs, squash merge into `main`.
