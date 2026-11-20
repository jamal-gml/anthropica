# OIDC from GitHub Actions to clouds and registries

The job requests a short-lived ID token (`permissions: id-token: write`); the cloud exchanges it for temporary credentials if the token's claims match a trust policy. No stored keys.

## Token claims to constrain on
| Claim | Example | Use |
|---|---|---|
| `aud` | `sts.amazonaws.com`, `api://AzureADTokenExchange` | provider-specific audience |
| `sub` | `repo:org/repo:environment:production` | primary binding |
| | `repo:org/repo:ref:refs/heads/main` | branch-bound |
| | `repo:org/repo:ref:refs/tags/v1.2.3` | tag-bound |
| | `repo:org/repo:pull_request` | PR runs (avoid for anything privileged) |
| `job_workflow_ref` | `org/repo/.github/workflows/deploy.yml@refs/heads/main` | bind to a specific reusable workflow |
| `repository_owner_id`, `repository_id` | numeric | resist rename/squatting |

Rules:
- Bind production roles to `environment:production` (the `sub` uses the environment when the job has one), so environment protection rules gate credential issuance.
- Never use wildcards like `repo:org/*` or `repo:org/repo:*` for write-capable roles.
- One role/identity per environment with least privilege; separate read-only plan roles from apply/deploy roles.

## AWS
- Create IAM OIDC provider `token.actions.githubusercontent.com`, audience `sts.amazonaws.com`.
- Role trust policy conditions: `StringEquals` on `token.actions.githubusercontent.com:aud` and `token.actions.githubusercontent.com:sub` (use `StringLike` only for deliberate patterns).
```yaml
permissions: { id-token: write, contents: read }
steps:
  - uses: aws-actions/configure-aws-credentials@<sha> # vX
    with:
      role-to-assume: arn:aws:iam::<account-id>:role/<deploy-role>
      aws-region: <region>
      role-session-name: gha-${{ github.run_id }}
```
- Error `Not authorized to perform sts:AssumeRoleWithWebIdentity`: `sub` mismatch (check environment vs branch form) or missing `id-token: write`.

## Azure
- Entra app registration or user-assigned managed identity with a **federated credential**: issuer `https://token.actions.githubusercontent.com`, subject e.g. `repo:org/repo:environment:production`, audience `api://AzureADTokenExchange`.
```yaml
- uses: azure/login@<sha> # vX
  with:
    client-id: ${{ vars.AZURE_CLIENT_ID }}
    tenant-id: ${{ vars.AZURE_TENANT_ID }}
    subscription-id: ${{ vars.AZURE_SUBSCRIPTION_ID }}
```
- IDs are not secrets; store as variables. Grant RBAC at the narrowest scope (resource group, not subscription).
- Error `AADSTS70021: No matching federated identity record`: subject string differs exactly (case, environment name, ref form).

## Google Cloud
- Workload Identity Pool + OIDC provider with issuer `https://token.actions.githubusercontent.com`, attribute mapping (`google.subject=assertion.sub`, `attribute.repository=assertion.repository`) and an **attribute condition** (e.g. `assertion.repository_owner_id == '<id>'`), which is mandatory to avoid accepting any GitHub repo.
- Prefer direct resource access for the principal, or impersonate a service account with `roles/iam.workloadIdentityUser`.
```yaml
- uses: google-github-actions/auth@<sha> # vX
  with:
    workload_identity_provider: projects/<num>/locations/global/workloadIdentityPools/<pool>/providers/<provider>
    service_account: <sa>@<project>.iam.gserviceaccount.com
```

## Registries and package publishing
- GitHub Container Registry: `GITHUB_TOKEN` with `packages: write` (no OIDC needed).
- Cloud registries: log in after the cloud OIDC step (`aws ecr get-login-password`, `az acr login`, `gcloud auth configure-docker`).
- npm and PyPI trusted publishing: register the repo + workflow file (+ environment) as a trusted publisher on the registry; the job needs `id-token: write`; no tokens stored.

## HashiCorp Vault and others
Use the JWT auth method with `bound_audiences` and `bound_claims` on `sub`/`repository`; same principles.

## Debug claims
Print the claims (never the token) in a test run:
```bash
curl -sH "Authorization: bearer $ACTIONS_ID_TOKEN_REQUEST_TOKEN" "$ACTIONS_ID_TOKEN_REQUEST_URL&audience=<aud>" \
  | jq -r .value | cut -d. -f2 | base64 -d 2>/dev/null | jq '{sub,aud,repository,environment,job_workflow_ref}'
```
Remove the step after debugging.
