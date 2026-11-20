# Documentation templates and docs-as-code

## How-to guide
````markdown
# How to <achieve outcome>
<One line: when you'd need this.>

## Before you start
- <Prerequisite with version/access>

## Steps
1. <Action>.
   ```bash
   <command>
   ```
   You should see: `<key output line>`
2. ...

## Verify
<Command or check that proves success.>

## Troubleshooting
| Problem | Cause | Fix |

## Related
<Links to reference/explanation.>
````

## Tutorial
```
# <Build/learn X> in <time>
What you'll build: <concrete result, screenshot if useful>
What you need: <minimal prerequisites>
## Step 1: <...>   <- every step ends with a visible result
...
## What you learned / Next steps
```
- One path, no options; tested end to end on a clean environment before publishing.

## Reference (config / CLI / env vars)
```
# <Component> configuration reference
| Key / flag / variable | Type | Required | Default | Description | Example |
|---|---|---|---|---|---|
```
- Generate from source (schema, CLI help, OpenAPI) where possible; hand-written reference drifts.
- Alphabetical or grouped by feature, consistently; document every option, including deprecated ones (marked, with replacement).

## Explanation / architecture note
```
# <Topic>: how it works and why
## Context        <- problem and constraints
## How it works   <- components, data flow, diagram (Mermaid)
## Design choices <- decision, alternatives, trade-offs (link the ADR in `plan-architecture` if one exists)
## Limitations and future work
```

## CONTRIBUTING
Setup, branch and commit conventions (link `build-commits`), test/lint commands, PR process (link `review-pull-request`), code of conduct, how to report security issues privately.

## Docs-as-code checks
| Check | Typical tools | Command example |
|---|---|---|
| Markdown lint | markdownlint-cli2 | `npx markdownlint-cli2 "**/*.md" "#node_modules"` |
| Link check | lychee, markdown-link-check | `lychee --no-progress "**/*.md"` |
| Prose/style lint | Vale (with a style guide), cspell for spelling | `vale docs/`, `npx cspell "**/*.md"` |
| Command verification | Run in a clean container/VM; doc-test tools for code samples | `docker run --rm -it <base-image>` then follow the quick start |
| Diagrams | Mermaid/PlantUML rendered in CI | render step in docs build |
| Docs site build | The repo's static-site generator, with warnings as errors | `<build> --strict` |

- Add checks to CI on paths `**/*.md` and `docs/**`; treat broken internal links as failures, external link failures as warnings (flaky).
- Keep docs next to the code they describe; one owner per doc (CODEOWNERS).

## Docs review checklist
- [ ] Reader and doc type clear from the first paragraph
- [ ] Quick start works on a clean machine, copy-paste as written
- [ ] Commands, paths, env vars, versions match the repo today
- [ ] No secrets or real internal values; placeholders explained
- [ ] Links resolve; no duplicated content that will drift
- [ ] Headings descriptive; scannable; tables for mappings
- [ ] Destructive steps warned and reversible
- [ ] Stale sections removed; `TODO(verify)` items listed in the PR
