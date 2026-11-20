# Workflow patterns

## Reusable workflow (`workflow_call`)
```yaml
# .github/workflows/deploy.yml
on:
  workflow_call:
    inputs:
      environment: { type: string, required: true }
      image-digest: { type: string, required: true }
    secrets: {}            # prefer OIDC; list only what is truly needed
permissions: { contents: read }
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: ${{ inputs.environment }}
    permissions: { contents: read, id-token: write }
    concurrency: { group: "deploy-${{ inputs.environment }}", cancel-in-progress: false }
    timeout-minutes: 30
    steps: [...]
```
Caller: `uses: ./.github/workflows/deploy.yml` (same repo) or `uses: org/repo/.github/workflows/deploy.yml@<sha>`, with `with:` and `secrets: inherit` only when the callee is trusted and needs them.
- Called workflow permissions cannot exceed the caller's. Nesting depth and reusable workflows per run are capped; check current docs before deep nesting.
- `env:` of the caller is not passed to the callee; pass values as inputs.
- Reusable workflow vs composite action: whole jobs with runners/environments -> reusable workflow; shared step sequence inside a job -> composite action (`.github/actions/<name>/action.yml`, `runs: using: composite`, every `run` needs `shell:`).

## Matrix
```yaml
strategy:
  fail-fast: false
  max-parallel: 4
  matrix:
    os: [ubuntu-latest, windows-latest]
    node: [20, 22]
    include: [{ os: ubuntu-latest, node: 22, coverage: true }]
    exclude: [{ os: windows-latest, node: 20 }]
```
Dynamic matrix: a first job writes JSON to `$GITHUB_OUTPUT`; the next uses `matrix: ${{ fromJSON(needs.plan.outputs.matrix) }}`. Max 256 jobs per matrix.

## Caching
| Need | Use |
|---|---|
| Package manager deps | setup action `cache:` input (npm, yarn, pnpm, pip, poetry, gradle, maven, go) |
| Other tool caches | `actions/cache` with `key: ${{ runner.os }}-x-${{ hashFiles('**/lockfile') }}` and `restore-keys:` prefix |
| Docker layers | `docker/build-push-action` with `cache-from/cache-to: type=gha,mode=max` |
- Caches are scoped by branch: PRs read the base branch cache; they cannot poison `main`'s cache. Total 10 GB per repo, least-recently-used eviction.
- Never cache secrets, build outputs that should be rebuilt, or `node_modules` across lockfile changes without the hash in the key.

## Artifacts between jobs
`actions/upload-artifact` / `download-artifact` (v4: names must be unique per run; set `retention-days`). For deploys, pass the immutable reference (image digest, artifact checksum) via job `outputs`, not a rebuild.

## Monorepo filters and required checks
- `on.pull_request.paths` skips the whole workflow; a required check from it then stays "Expected" forever and blocks merge.
- Fix: run the workflow always, decide inside with a changes job (`dorny/paths-filter` pinned, or `git diff --name-only`), and gate jobs with `if: needs.changes.outputs.api == 'true'`. Skipped jobs (via `if`) report success.
- Aggregate for branch protection with one always-running job:
```yaml
  ci-ok:
    if: ${{ always() }}
    needs: [lint, test, build]
    runs-on: ubuntu-latest
    steps:
      - if: ${{ contains(needs.*.result, 'failure') || contains(needs.*.result, 'cancelled') }}
        run: exit 1
```
Make `ci-ok` the single required check.

## Deploy jobs
- `environment: { name: production, url: ${{ steps.deploy.outputs.url }} }`; protection rules: required reviewers (not the actor), deployment branch/tag policy, wait timer.
- Trigger production from tags or manual dispatch of a known version, not from every push.
- Deploy the digest built earlier; verify with a smoke test step; record the deployment.

## Release jobs
- Tag/release creation needs `contents: write`. Package publishing via OIDC trusted publishing needs `id-token: write` (npm `--provenance`, PyPI `pypa/gh-action-pypi-publish`).
- Release bots (release-please, changesets) need an app token if their PRs must trigger CI.

## Speed
- Run lint, unit tests and build in parallel jobs; serialise only real dependencies via `needs:`.
- Shallow checkout (default `fetch-depth: 1`); use `fetch-depth: 0` only for changelog/versioning jobs.
- Larger runners or test sharding for suites over ~10 minutes.
- Cancel superseded CI runs with `concurrency`.

## Scheduled workflows
- Cron is UTC, may be delayed under load, and is disabled after 60 days of repo inactivity in public repos. Avoid minute 0 of the hour.
