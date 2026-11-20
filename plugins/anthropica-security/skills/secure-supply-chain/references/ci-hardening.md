# CI/CD hardening checklist

Run `zizmor .github/workflows` (security lints) and `actionlint` first; then check by hand. Examples use GitHub Actions; the principles apply to any CI.

## Token and permissions
- Top-level `permissions: {}` or `contents: read`; grant writes per job only where needed (`contents: write` for release, `id-token: write` for OIDC, `packages: write` for registry push, `attestations: write` for attestations).
- Repository/org default `GITHUB_TOKEN` permission set to read-only (see `secure-repo-governance`).
- `actions/checkout` with `persist-credentials: false` unless a later step must push.
- No personal access tokens for automation; use a GitHub App with minimal scopes or the job token.

## Dangerous triggers
- `pull_request_target` and `workflow_run` run with base-repo secrets and write token. Never check out or execute PR head code (`ref: ${{ github.event.pull_request.head.sha }}`) in them, including `npm install`, `make`, or scripts from the PR. If you must, split into an unprivileged `pull_request` build that uploads an artifact and a privileged job that only reads data from it, treating it as untrusted.
- `issue_comment` / label-triggered workflows: check the actor's permission (`author_association` is not enough on its own for forks) before running privileged steps.
- Fork PRs: secrets are not passed to `pull_request` from forks by default - keep it that way; require approval for first-time contributors' workflows.

## Script injection
- Never interpolate untrusted context directly into `run:` (`${{ github.event.pull_request.title }}`, `head_ref`, issue body, commit messages, branch names). Pass via `env:` and quote the shell variable:
  ```yaml
  env:
    TITLE: ${{ github.event.pull_request.title }}
  run: echo "$TITLE"
  ```
- Same for `actions/github-script` inputs: read from `context` or env, not string templates.
- Writing untrusted data to `$GITHUB_ENV` or `$GITHUB_OUTPUT` can inject variables (`LD_PRELOAD`, `NODE_OPTIONS`) into later steps.

## Third-party actions and tools
- Pin by full commit SHA with a version comment; verify the SHA belongs to the upstream repo (not a fork commit reachable via the same name). Tools: `pinact run`, `ratchet pin`.
- Prefer actions from verified creators or well-maintained repos; read composite actions and their transitive `uses:`; Docker-based actions should reference images by digest.
- Replace `curl | bash` installers with pinned releases verified by checksum or signature.

## Cloud credentials
- Use OIDC federation, not long-lived cloud keys in secrets. Restrict the cloud-side trust policy by `sub` claim to repo + branch/tag or environment (e.g. `repo:<owner>/<repo>:environment:production`), audience, and issuer. Wildcards on `sub` across repos are Critical.
- Separate roles per environment; deploy roles cannot read unrelated data.

## Secrets handling
- Environment-scoped secrets for production, gated by required reviewers and branch/tag rules.
- Do not echo secrets, write them to files that are uploaded as artifacts, or pass them as command-line args visible in process lists. Mask derived values with `::add-mask::`.
- Rotate any secret that a compromised action or workflow version could have read.

## Runners
- GitHub-hosted or ephemeral self-hosted runners only for untrusted code; never attach persistent self-hosted runners to public repos.
- Separate runner groups for privileged (deploy) and unprivileged jobs; restrict which repos and workflows can use each group.
- Restrict egress where possible (egress-allowlist agent or network policy) and audit it for build jobs that publish artifacts.

## Caches and artifacts
- Caches are writable by any workflow on the branch and readable by PRs: never cache secrets; do not restore caches from untrusted refs into release jobs; key release caches separately or build releases without cache.
- Artifacts passed between jobs are untrusted input to the consuming job; verify checksums or attestations before deploying.

## Release pipeline integrity
- Release only from protected refs (tags matching a ruleset, `main`) via a dedicated workflow with an environment gate.
- Build once, promote the same digest through environments; deploy by digest, not tag.
- Generate provenance in the same job that builds (see `sbom-provenance.md`); reusable workflows give stronger isolation (SLSA Build L3).
- Require two-person review for workflow file changes via CODEOWNERS on `.github/workflows/`.

## Other CI systems
- GitLab: protected variables and branches, `CI_JOB_TOKEN` scope allowlist, `rules` avoiding secrets on fork MRs, pinned `include:` refs.
- Jenkins/others: credentials binding scoped per job, no script approval bypass, agents isolated per trust level, plugin versions pinned and patched.
