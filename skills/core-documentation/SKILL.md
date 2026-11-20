---
name: core-documentation
description: Write or update docs verified against the repo (README, setup and how-to guides, tutorials, reference, explanations, runbooks) using Diataxis and docs-as-code checks. Use when writing or fixing a README, CONTRIBUTING, setup guide, config reference or runbook, or when a change makes docs stale. Not for ADRs (use `plan-architecture`), changelogs (use `ship-release`) or postmortems (use `operate-incident`).
---

# Technical Documentation

Write the smallest correct document that lets the intended reader succeed, and prove every command in it works.

## Procedure
1. **Reader and task.** Name the reader (new contributor, operator on call, API consumer, evaluator) and the one thing they need to do or understand.
2. **Pick the type** (Diataxis, table below). One document = one type; link between types instead of mixing them.
3. **Gather facts from the repo,** not memory: build files, scripts, CI workflows, env samples, config schemas, entry points, existing docs. Note anything you cannot verify.
4. **Outline**, then write. Start from the templates: README skeleton below; runbook in `references/runbook-template.md`; other types in `references/templates.md`.
5. **Verify every command and path:** run it (or a safe dry-run) in a clean shell/environment where possible; check paths exist, env var names match the code, versions match the manifest. Mark unverifiable steps `TODO(verify): <what>`.
6. **Docs-as-code checks:** markdown lint, link check, spelling/prose lint if the repo has them; fix what they flag. Read `references/templates.md` for tool options when the repo has none.
7. **Ship with the change:** docs update in the same PR as the code; update or delete the stale parts you found, and only those.

## Rules
| Type | Reader's state | Shape | Avoid |
|---|---|---|---|
| Tutorial | Learning, new | A guaranteed-success lesson with one path, concrete steps, visible results | Options, edge cases, theory |
| How-to guide | Working, has a goal | Numbered steps to one outcome; prerequisites first; assumes basics | Teaching concepts |
| Reference | Looking something up | Complete, consistent, structured like the thing it describes (CLI flags, API, config keys, env vars) | Narrative, opinions |
| Explanation | Wants to understand | Context, design reasons, trade-offs, alternatives | Step-by-step instructions |

**README skeleton** (keep it short; link out for depth)
```
# <Project>
<One sentence: what it is and who it's for.>
## Quick start      <- copy-paste commands to a working result in < 5 minutes
## Requirements     <- runtimes/tools with versions
## Usage            <- the 2-3 most common tasks
## Configuration    <- table: name | required | default | description (or link to reference)
## Development      <- install, test, lint, run locally
## Deployment / Operations <- link to runbooks
## Contributing / License / Support
```

- Quick start before reference; the first command block must work as written on a clean machine.
- Commands in fenced blocks with a language tag; one command per line; show expected output when it confirms success. Separate commands from output.
- State the platform when commands differ (bash vs PowerShell); don't assume one OS silently.
- Placeholders are explicit and consistent: `<your-api-url>`, `YOUR_TOKEN`; tell the reader where to get the value.
- Link to canonical sources (generated API docs, config schema, CI file) instead of copying them; copied docs go stale.
- Write in second person, present tense, active voice; descriptive headings; meaningful link text; alt text for diagrams; tables for mappings.
- Version-sensitive facts (versions, limits, URLs) belong in one place, referenced elsewhere.
- Diagrams as code (Mermaid, PlantUML) so they diff and review like text.
- Runbooks: each alert links to one; every step has the command, the expected result and what to do if it differs; include rollback and escalation.
- Delete or archive docs that describe removed features; wrong docs are worse than none.
- Keep the README under ~300 lines; move depth into `docs/` by Diataxis type and link it.
- Prefer text over screenshots (they go stale and can't be searched); date any screenshot you keep.

**Doc-update triggers** (a code change with any of these needs a doc change in the same PR)
| Change | Update |
|---|---|
| New/renamed env var, config key, CLI flag or default | Reference table, `.env.example`, README configuration |
| Changed install, build, test or run command | Quick start, CONTRIBUTING, CI badge/docs |
| New external dependency (service, queue, API key) | Requirements, setup guide, runbook prerequisites |
| New alert, job, or operational procedure | Runbook, linked from the alert |
| Public API change | API reference, migration notes (changelog via `ship-release`) |
| Architecture or major design change | Explanation doc; ADR via `plan-architecture` |

## Safety gates
- Never put real secrets, tokens, connection strings, internal hostnames or personal data in docs or examples; use placeholders and point to the secret store.
- Never invent commands, flags, endpoints, env vars or behaviour; unverified = `TODO(verify)`.
- Destructive commands in docs (data deletion, force push, production changes) carry a warning, a confirmation step and a rollback note.
- Ask before publishing docs to an external site or wiki, or changing docs owned by another team.

## Output
- The document(s), plus a short note: reader and type chosen, files changed, commands verified (and how), checks run (lint/link), `TODO(verify)` items remaining.
