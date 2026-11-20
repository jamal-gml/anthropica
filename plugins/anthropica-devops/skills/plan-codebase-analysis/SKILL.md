---
name: plan-codebase-analysis
description: Orient in an unfamiliar repository and trace a behaviour end to end (entry points, call paths, data flow, git history, tests as spec, conventions) with confidence-labelled findings. Use before broad or cross-module changes, when asked how X works or where Y is handled, or when a bug spans several files. Not for designing new systems (use `plan-architecture`).
---

# Codebase Analysis

Build a verified, minimal map of the code that matters for the task, read-only, before anyone edits.

## Procedure
1. **Frame the question.** Restate the target as something observable: a route, command, job, UI action, error message or data field. Note what a full answer must explain (entry, decision points, side effects, outputs).
2. **Orient in under 10 tool calls.**
   - `git status`, `git branch --show-current`, `git log --oneline -15` (state and recent intent).
   - Instruction files: `CLAUDE.md`, `AGENTS.md`, `.claude/rules/`, `CONTRIBUTING.md`, `docs/adr/`.
   - Manifests give stack and scripts: `package.json`, `pyproject.toml`, `go.mod`, `pom.xml`/`build.gradle`, `*.csproj`, `Cargo.toml`, `Gemfile`, `composer.json`. Workspace files (`pnpm-workspace.yaml`, `go.work`, `nx.json`, `turbo.json`, Cargo `[workspace]`) mean a monorepo: find the owning package first.
   - Entry points and wiring: `main`/`index`/`app`/`server` files, route registration, DI container, `Dockerfile` `CMD`/`ENTRYPOINT`, CI build/test steps, `Procfile`, deployment manifests.
   - Tree to depth 2, excluding `node_modules`, `vendor`, `dist`, `build`, `.venv`, `target` and generated code.
3. **Locate the anchor.** Search in this order and stop at the first solid hit:
   1. Exact user-visible strings (error text, log message, UI label, i18n key) - most unique.
   2. Route path or CLI command name.
   3. Symbol name, then its definition and references (LSP or `ctags` if available, else word-boundary search `\bname\b`).
   4. Config key, env var, DB table or column name.
   Read `references/search-recipes.md` when the stack is unfamiliar or the first two searches miss.
4. **Trace the call path.** From the anchor follow callers up to the entry point and callees down to side effects (DB, queue, HTTP, file, cache, email). Record each hop as `file:line - function - what it decides or changes`. Stop a branch when it leaves the question's scope or enters a third-party library.
5. **Trace the data.** For the key entity find where it is defined (schema/model), validated, transformed, persisted, read back and serialised. Note every place the shape or units change.
6. **Read the tests as spec.** Find tests that exercise the path (search for the function, route or fixture). They show intended behaviour, edge cases and how to run that slice. Missing tests are a finding.
7. **Use history for "why".** `git log -S'<string>' --oneline` (when a string appeared or vanished), `git log -L :<func>:<file>` (one function's history), `git blame -w -C -L a,b <file>` then `git show <sha>` for the message and linked issue. Do this before calling odd code a bug.
8. **Extract conventions** the change must follow: layering, naming, error handling, DI, logging, validation library, test style and location, file placement. Cite one exemplar file for each.
9. **Check what changes the path at runtime**: env vars, feature flags, per-environment config, build tags, tenant settings.
10. **Report** with the Output template. Stop there unless implementation was requested.

## Rules
- Read the smallest slice that answers the question: search, then targeted line ranges. Never read lockfiles, minified bundles, snapshots or generated code unless they are the subject.
- Folder names are hypotheses. Confirm architecture from imports, registration code and runtime config.
- Label every claim:
  | Label | Meaning |
  |---|---|
  | Verified | Read the code or ran it; cite `file:line` |
  | Inferred | Follows from verified facts plus a convention; say which |
  | Unverified | Plausible, not checked; say what would confirm it |
- Dynamic dispatch hides callers (DI, reflection, decorators, event buses, string-keyed handlers, ORM hooks, middleware, signals). Search for the registration key or event name, not only the symbol.
- When several implementations match (interfaces, plugins, per-tenant variants, old and new copies in a monorepo), find which one the entry point actually wires for the case in question.
- Prefer a cheap, side-effect-free check (run one narrow test, a `--help`, a dry-run) over reasoning about behaviour.
- Delegate wide fan-out searches (many naming variants or packages) to a search subagent when available; keep the conclusions, not file dumps.
- Note dead code, TODO/FIXME/HACK near the path and skipped tests: they are risks for the change.
- Time-box: if the anchor is not found after ~6 searches, report what was tried and ask for a hint (URL, screenshot, log line) instead of reading more files.

## Safety gates
- Read-only. Do not edit, format, install dependencies, run migrations, or start anything that writes to shared resources without asking.
- Do not print secrets found in `.env`, config or history; reference the file and key name only.
- Do not check out branches or stash the user's changes to look around; use `git show <ref>:<path>` instead.

## Output
```
## Question
<one line>

## Answer (short)
<2-4 lines>

## Map
Entry: <file:line> <route/command/job>
Path:
1. <file:line> <function> - <what it does or decides>
2. ...
Side effects: <tables / queues / external calls / files>
Data shape changes: <where and how>
Runtime switches: <flags, env vars, config>

## Conventions to follow
- <convention> - exemplar <file>

## Tests
- Covering: <test files, command to run them>
- Gaps: <untested branches>

## Risks and unknowns
- [Inferred|Unverified] <item> - <how to confirm>

## Recommended next step
<smallest change, or the question that blocks it>
```

---

**Phase:** Plan. **Next:** when this step is done, offer Build: `/anthropica-devops:build-commits`, `/anthropica-devops:build-containers`, `/anthropica-devops:build-iac`. Phase map: `/anthropica-devops:workflow`.
