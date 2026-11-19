---
name: github-actions
description: Author, harden and debug GitHub Actions workflows (least-privilege permissions, SHA-pinned actions, OIDC instead of long-lived secrets, caching, concurrency, environments, reusable workflows). Use when creating or editing files in .github/workflows, or when a workflow run fails.
---

# GitHub Actions

## Authoring checklist
1. **Triggers:** be explicit (`push` with `branches`, `pull_request`, `workflow_dispatch`, `schedule`). Add `paths` filters for monorepos.
2. **Permissions:** set `permissions: {}` or `contents: read` at the top, and grant more per job only where needed (`id-token: write`, `pull-requests: write`, `packages: write`).
3. **Pin actions to full commit SHAs** with the version as a comment:
   `uses: actions/checkout@<40-char-sha> # v4.2.2`. Look the SHA up (`gh api repos/<owner>/<action>/git/ref/tags/<tag>`); never guess one. Let Dependabot keep them current.
4. **Cloud credentials:** use OIDC federation (`id-token: write` plus the provider's login action) instead of stored access keys.
5. **Secrets:** reference with `${{ secrets.NAME }}`; never echo them. Pass them to scripts through `env:`, not inline in `run:`.
6. **Concurrency:** `concurrency: { group: ${{ github.workflow }}-${{ github.ref }}, cancel-in-progress: true }` for CI. Never cancel in-progress deploys.
7. **Timeouts:** set `timeout-minutes` on every job.
8. **Caching:** use the setup action's built-in cache (`cache: npm`, `pip`, `gradle`) keyed on lockfiles.
9. **Deployments:** use `environment:` with protection rules (required reviewers, branch restrictions) for staging and production.
10. **Reuse:** extract repeated logic into reusable workflows (`workflow_call`) or composite actions.

## Security pitfalls
- **Script injection:** never interpolate untrusted input (`github.event.pull_request.title`, issue bodies, branch names) directly in `run:`. Pass it through `env:` and quote it.
- **`pull_request_target` and `workflow_run`:** run with write tokens and secrets. Never check out or execute PR code in them.
- **Forks:** secrets are unavailable to fork PRs by design; do not work around it.
- **Self-hosted runners:** never use them for public repositories.
- **Artifacts:** do not upload files containing secrets. Set `retention-days`.

## Typical CI job
```yaml
permissions:
  contents: read
jobs:
  test:
    runs-on: ubuntu-latest
    timeout-minutes: 15
    steps:
      - uses: actions/checkout@<sha> # vX.Y.Z
      - uses: actions/setup-node@<sha> # vX.Y.Z
        with: { node-version-file: .nvmrc, cache: npm }
      - run: npm ci
      - run: npm test
```

## Debugging a failed run
1. `gh run list --workflow <file>` then `gh run view <id> --log-failed`.
2. Find the first real error, not the last line. Check whether it fails on `main` too.
3. Reproduce locally where possible. Re-run with debug logging (`ACTIONS_STEP_DEBUG`) only if needed.
4. Fix the cause. Never weaken permissions gates, delete required checks or add `continue-on-error` to hide a failure.

## Validation
Run `actionlint` if available. Confirm required status-check names still match branch protection after renaming jobs.

## Output
Summarise triggers, permissions, secrets and environments used, and any required repository settings (by name only, never values).
