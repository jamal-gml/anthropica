---
name: build-project-setup
description: Onboard a repo for AI coding assistants (lean CLAUDE.md/AGENTS.md with verified commands, conventions, gotchas) and scaffold new projects (lockfile, formatter, linter, tests, pre-commit, CI, .gitignore, .editorconfig, README) without overwriting files. Use when starting a project, onboarding Claude, or creating or trimming CLAUDE.md/AGENTS.md. Not for pipeline design (use `ship-ci-cd`).
---

# Project Setup

Make a repository fast and safe to work in: verified commands, explicit conventions, and the minimum tooling that catches mistakes.

## Procedure
1. **Inventory first.** `git status`, branch, and list existing: `CLAUDE.md`, `CLAUDE.local.md`, `AGENTS.md`, `.claude/` (settings, rules, skills, agents), `.cursor/rules`, `.github/copilot-instructions.md`, `CONTRIBUTING.md`, `README`, `.editorconfig`, `.gitignore`, `.gitattributes`, lint/format/test configs, pre-commit config, CI workflows, lockfiles, version files (`.nvmrc`, `.python-version`, `.tool-versions`, `rust-toolchain.toml`, `global.json`). Existing files are inputs to merge, not to replace.
2. **Classify the job**: (a) onboarding an existing repo - go to 3; (b) new project - go to 7; (c) both.
3. **Discover real commands** from manifests (`package.json` scripts, `Makefile`, `justfile`, `Taskfile.yml`, `pyproject.toml`, `tox.ini`, `noxfile.py`), CI workflows (the most trustworthy source: they actually run) and README. Identify install, build, run, test (all and single test), lint, format, typecheck, migrate, codegen.
4. **Verify commands** by running the cheap, side-effect-free ones (lint, typecheck, single unit test, `--help`). Record the outcome. Never run deploy, migrate against a shared DB, publish or seed commands to verify them; mark them "not verified".
5. **Collect what is non-obvious**: conventions that differ from the ecosystem default, generated files that must not be edited, required env vars (names only), services needed for tests, flaky areas, ordering rules (codegen before build), monorepo package boundaries, protected branches and commit format.
6. **Write or update the instruction file** following Rules and `references/claude-md-template.md` (read it when creating or restructuring an instruction file). Show the proposed content or diff and get approval before overwriting an existing file.
7. **New project scaffolding** (only what the user asked for plus the essentials below). Read `references/toolchains.md` when choosing per-language tools and init commands. Essentials, in order:
   1. Pin the runtime version (version file and/or `engines`/`requires-python`/`go` directive).
   2. Package manager with a committed lockfile.
   3. `.gitignore` for the language, OS and editor files, `.env*` except `.env.example`.
   4. `.editorconfig` and `.gitattributes` (`* text=auto eol=lf`, binary types marked).
   5. Formatter and linter with config, plus typechecker where the language has one.
   6. Test runner with one real passing test.
   7. Pre-commit hook running format and lint on staged files (fast; heavy checks stay in CI).
   8. CI running the same commands as local (hand to `ship-github-actions` or `ship-ci-cd`).
   9. `README.md` with purpose, prerequisites, setup, run, test; `.env.example` with names and safe placeholder values.
8. **Validate**: run install, format check, lint, typecheck and tests from a clean state; fix only setup problems.
9. **Report** with the Output template.

## Rules
- Instruction files are loaded into every session: every line costs tokens on every task. Target under ~100 lines for a root `CLAUDE.md`; hard ceiling ~200. Move area-specific rules into subdirectory `CLAUDE.md` files or `.claude/rules/`; move procedures into skills; link to docs instead of copying them.
- Include only what changes behaviour and cannot be cheaply discovered from the code:
  | Include | Exclude |
  |---|---|
  | Exact verified commands, especially single-test and pre-commit checks | Generic advice ("write clean code", "add tests") |
  | Conventions that differ from defaults, with an exemplar path | File-by-file descriptions of the tree |
  | Gotchas that caused real mistakes | Content of README/docs (link instead) |
  | Do-not-touch: generated, vendored, migrations already applied | Secrets, tokens, internal hostnames, personal data |
  | Architecture pointers (where X lives, which layer owns Y) | Long examples, full style guides, changelogs |
  | Safety rules (ask before X), commit/branch format | Aspirations that are not enforced |
- Write rules as imperative, testable statements: "Run `pnpm test -- <file>` for a single test", not "Testing is important".
- Mark unverified commands `(not verified)`; never invent commands, coverage numbers, architecture or deployment behaviour.
- One source of truth across tools: if `AGENTS.md` already exists, keep shared content there and have `CLAUDE.md` import it (`@AGENTS.md`) plus Claude-specific notes, rather than duplicating. Do not create `AGENTS.md` unless other agents are used or the user asks.
- Personal or machine-specific notes go in `CLAUDE.local.md` (gitignored), never in the shared file.
- Permissions: put shared, safe allowlists (test, lint, read-only git) in `.claude/settings.json`; never allow destructive or deploy commands by default. Hooks only for deterministic checks that must always run (format on edit, block edits to generated files), opt-in and fast.
- Preserve existing conventions and tool choices. Changing formatter, linter, package manager or test framework is a separate, approved change.
- Formatter and linter must not fight: disable lint style rules that the formatter owns.
- Pin tool versions in config (dev dependencies, pre-commit `rev:`) so local, hook and CI results match.

## Safety gates
- Never overwrite or delete an existing file (instruction files, configs, README, CI). Propose a merge/diff and wait for approval.
- Never run project generators (`create-*`, `init`, `new`) in a non-empty directory without asking; run them in an empty temp directory and copy selectively if needed.
- Ask before adding any dependency, dev tool or pre-commit hook; state what it does and its maintenance cost.
- Ask before running commands that touch databases, cloud resources, package registries or remote git.
- Do not write secrets into any file; `.env.example` holds placeholders only.

## Output
```
## Setup summary
Stack: <languages, frameworks, package manager, runtime version>
Created: <files>
Updated (approved): <files and what changed>
Left unchanged: <existing files deliberately not touched>

## Verified commands
| Purpose | Command | Result |
|---|---|---|

## Not verified / unknown
- <command or fact> - <why>

## Follow-ups
- <e.g. add CI, branch protection, dependency updates>
```

---

**Phase:** Build. **Next:** when this step is done, offer Test: `/anthropica-engineering:test-performance`, `/anthropica-engineering:test-validation`. Phase map: `/anthropica-engineering:workflow`.
