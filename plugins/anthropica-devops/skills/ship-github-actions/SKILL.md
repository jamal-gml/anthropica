---
name: ship-github-actions
description: Author, harden and debug GitHub Actions workflows (least-privilege token, SHA-pinned actions, OIDC to clouds, script-injection safety, caching, matrices, concurrency, environments, reusable workflows). Use when editing .github/workflows or action.yml, or when a run fails or is slow. Not for platform-neutral pipeline design (use `ship-ci-cd`) or rulesets (use `secure-repo-governance`).
---

# GitHub Actions

Write workflows that are least-privilege, reproducible, fast and safe against untrusted input.

## Procedure
1. **Inventory:** `ls .github/workflows`, `gh workflow list`, existing reusable workflows and composite actions (`.github/actions/*/action.yml`). Extend before adding.
2. **Know the gates:** required check names (`gh api repos/{owner}/{repo}/rules/branches/<default>`), environments (`gh api repos/{owner}/{repo}/environments`). Renaming a job or workflow can break a required check.
3. **Write or edit** following the checklist below.
4. **Pin actions:** resolve each tag to a commit SHA: `gh api repos/<owner>/<repo>/commits/<tag> -q .sha` (works for lightweight and annotated tags). Never guess a SHA. Comment the version: `@<sha> # v4.2.2`.
5. **Validate:** `actionlint` (includes shellcheck for `run:`), and `zizmor .github/workflows` for security if available. Fix every finding or justify it.
6. **Test** on a branch via `pull_request` or `workflow_dispatch` (`gh workflow run <file> --ref <branch>`); watch with `gh run watch`.
7. **Report** permissions, secrets and environments by name, and any repo settings that must change.

## Authoring checklist
1. **Triggers explicit:** `push: branches: [main]`, `pull_request`, `workflow_dispatch`, `schedule`, `release`. Add `paths`/`paths-ignore` for monorepos (but see required-check gotcha in `references/patterns.md`).
2. **Permissions:** top-level `permissions: {}` or `contents: read`; grant per job only what it needs (`id-token: write` for OIDC, `pull-requests: write` to comment, `packages: write` to push images, `contents: write` to tag/release, `attestations: write` for provenance).
3. **Pinned actions** to full 40-char SHAs, including first-party `actions/*`. Let the dependency bot update them (`package-ecosystem: github-actions`).
4. **Cloud auth via OIDC**, never stored access keys. Read `references/oidc.md` for AWS, Azure, GCP and registry setups.
5. **Secrets** only via `${{ secrets.X }}` into `env:` of the step that needs them; never echo, never pass on the command line. Scope deploy secrets to environments, not the repo.
6. **Untrusted input never in `run:` directly.** Pass through `env:` and quote: `env: TITLE: ${{ github.event.pull_request.title }}` then `"$TITLE"`.
7. **Concurrency:** CI `group: ${{ github.workflow }}-${{ github.ref }}`, `cancel-in-progress: true`. Deploys: `group: deploy-<env>`, `cancel-in-progress: false` (never cancel a running deploy).
8. **`timeout-minutes`** on every job (default is 360). Typical: lint 10, tests 20-30, deploy 30.
9. **Caching:** setup action's built-in cache (`cache: npm|pip|gradle|maven`, `cache-dependency-path` for monorepos) keyed on lockfiles; `actions/cache` only for other tool caches.
10. **Deploy jobs** use `environment: <name>` with required reviewers, branch/tag restrictions and wait timers configured in settings.
11. **Reuse** with reusable workflows (`workflow_call`) for whole pipelines, composite actions for step sequences.
12. **Shell safety:** `defaults: run: shell: bash` gives `-eo pipefail`; add `set -u` in scripts.

## Rules
- **Dangerous triggers:** `pull_request_target` and `workflow_run` run with base-repo secrets and a write token. Never check out or execute PR head code in them. If you must label/comment, use only event metadata via `env:`.
- **Fork PRs:** no secrets and a read-only token by design; never work around it (e.g. switching to `pull_request_target`).
- **Self-hosted runners:** never for public repos; use ephemeral runners, isolated per trust level.
- **`GITHUB_TOKEN` actions do not trigger new workflow runs** (except `workflow_dispatch`/`repository_dispatch`). Release bots that need CI on their PRs need an app token, not a PAT tied to a person.
- **Artifacts:** `retention-days` set; never upload `.env`, credentials or full workspaces. Artifacts of `pull_request` runs are untrusted when consumed by `workflow_run`.
- **Outputs:** write via `echo "name=value" >> "$GITHUB_OUTPUT"` (`set-output` is removed). Mask derived secrets: `echo "::add-mask::$VALUE"`.
- **Expressions:** `if: ${{ always() }}` runs even on cancel; prefer `if: ${{ !cancelled() }}` for cleanup/report steps. `success()` is implicit.
- **Matrices:** `fail-fast: false` for test matrices so one failure does not hide others; cap with `max-parallel` for rate-limited targets.
- **Required checks with path filters or matrices:** skipped workflows leave required checks "pending". Use an always-running aggregator job (`references/patterns.md`).

Read `references/security.md` when a workflow handles untrusted input, uses `pull_request_target`/`workflow_run`, publishes artifacts, or is being audited.
Read `references/patterns.md` when writing reusable workflows, matrices, caches, monorepo filters, release or deploy jobs.
Read `references/oidc.md` when a job authenticates to a cloud provider or registry.
Read `references/debugging.md` when a run fails, hangs, is skipped unexpectedly or is slow.

## Typical CI job
```yaml
name: ci
on: { pull_request: {}, push: { branches: [main] } }
permissions: { contents: read }
concurrency: { group: "${{ github.workflow }}-${{ github.ref }}", cancel-in-progress: true }
jobs:
  test:
    runs-on: ubuntu-latest
    timeout-minutes: 20
    steps:
      - uses: actions/checkout@<sha> # vX.Y.Z
        with: { persist-credentials: false }
      - uses: actions/setup-node@<sha> # vX.Y.Z
        with: { node-version-file: .nvmrc, cache: npm }
      - run: npm ci
      - run: npm test
```

## Safety gates
- Ask before: pushing workflow changes, running `workflow_dispatch` against deploy workflows, re-running deploy jobs, creating or changing environments, secrets, variables or rulesets.
- Never weaken `permissions:`, remove environment protection, delete or rename required checks, or add `continue-on-error` to hide a failure.
- Production deploy runs need the user's explicit approval with version and SHA stated.

## Output
```
Workflow: <file>   Triggers: <events + filters>
Jobs: <name> (timeout, runner) ...
Permissions: top <..>; per job <job: scopes>
Actions pinned: <n>/<n> to SHA
Secrets/vars used (names only): <..>   Environments: <..>
OIDC: <provider, role/identity name>
Validation: actionlint <pass/fail>, zizmor <pass/fail/not run>
Settings needed: <required checks, environment rules, OIDC trust>
```

---

**Phase:** Ship. **Next:** when this step is done, offer Operate: `/anthropica-devops:operate-incident`, `/anthropica-devops:operate-observability`. Phase map: `/anthropica-devops:workflow`.
