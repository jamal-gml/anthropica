---
name: project-bootstrap
description: Establish a new or existing software project for Claude with concise instructions, conventions, quality gates and optional skills. Use when starting a project, onboarding Claude to a repository, or creating CLAUDE.md/AGENTS.md.
disable-model-invocation: true
---

# Project Bootstrap

Prepare a project for Claude without generating a bloated instruction file.

## Discovery
1. Inspect repository root, current branch, working-tree status and existing instruction files.
2. Identify languages, frameworks, package managers, build/test/lint commands, deployment targets and CI.
3. Inspect only the key manifests, entry points, CI workflows and representative tests.
4. Ask before replacing or deleting existing instructions or configuration.

## Deliverables
- A concise root `CLAUDE.md` containing project purpose, architecture pointers, commands, conventions and safety constraints.
- Optional `AGENTS.md` only when cross-agent instructions are useful.
- Scoped `.claude/rules/` for path-specific or domain-specific constraints.
- Skills for repeatable procedures; do not copy entire skill bodies into `CLAUDE.md`.
- Hooks only when a deterministic check cannot be handled more simply; keep them opt-in.

## Rules
- Record verified commands only. Mark unknown commands as unknown.
- Link to deeper documentation rather than embedding long reference material.
- Do not invent architecture, test coverage, deployment behaviour or secrets.
- Preserve existing conventions unless the user approves a change.
- Keep generated instructions short, testable and repository-specific.
- Finish with a summary of files created/changed and how to validate them.
