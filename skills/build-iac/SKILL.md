---
name: build-iac
description: Write, review and apply IaC (Terraform/OpenTofu, Bicep, Pulumi, CloudFormation/CDK) with saved plans, locked remote state, pinned versions, policy checks, safe state surgery and drift control. Use when changing cloud resources, IAM or networking, reading a plan for replace/destroy, importing or moving state, or fixing drift. Not for pipelines (use `ship-ci-cd`) or Kubernetes manifests (use `build-containers`).
---

# Infrastructure as Code

Change cloud infrastructure only through reviewed, reproducible, policy-checked code with a saved plan.

## Procedure
1. **Discover.** Tool and version (`terraform version`, `tofu version`, `.terraform-version`, `bicepconfig.json`, `Pulumi.yaml`, `cdk.json`), providers, backend, environment layout, CI that applies it. Read existing modules, naming, tagging and variable conventions before writing.
2. **Confirm target.** Name the environment, account/subscription/project and region the change hits. Check the active credentials (`aws sts get-caller-identity`, `az account show`, `gcloud config list`). Never assume production is safe to touch.
3. **Write.** Reuse modules; follow the Rules below. Keep one logical change per PR.
4. **Static checks** (cheap first, fail fast):
   ```bash
   terraform fmt -check -recursive && terraform init -backend=false && terraform validate
   tflint --recursive                       # if configured
   trivy config . || checkov -d .           # misconfig scan (tfsec is now part of trivy)
   conftest test tfplan.json -p policy/     # OPA policies, after step 5
   ```
5. **Plan to a file** and render it for review and policy:
   ```bash
   terraform plan -out=tfplan -lock-timeout=5m
   terraform show -json tfplan > tfplan.json
   terraform show tfplan | grep -E '# .* (must be replaced|will be destroyed)'
   ```
6. **Read the plan** (see Rules: plan reading). List every destroy, replace and IAM/network/data-store change with the attribute that forces it.
7. **Gate.** Present the plan summary and wait for explicit approval (see Safety gates).
8. **Apply the saved plan only:** `terraform apply tfplan`. Never `apply -auto-approve` without a plan file outside CI. If the plan is stale (`Saved plan is stale`), re-plan and re-review.
9. **Verify.** Re-run `terraform plan -detailed-exitcode`; exit 0 = no drift. Smoke-test the resource. Record outputs that changed.

## Rules

**Plan reading**
| Symbol / text | Meaning | Action |
|---|---|---|
| `+` create | new resource | check naming, tags, public exposure |
| `~` update in-place | safe-ish | check IAM/SG/policy diffs line by line |
| `-/+` or `+/-` / "forces replacement" | destroy and recreate | find the forcing attribute; for data stores, stop and ask |
| `-` destroy | deletion | require reason; check `prevent_destroy` / deletion protection |
| `(known after apply)` on IDs used elsewhere | downstream churn | expect dependent replacements |
| `moved` / `imported` | state-only | confirm zero real changes accompany it |
- Replacement of databases, buckets, disks, key vaults/KMS keys, DNS zones or load balancers with fixed IPs = data loss or outage risk. Prefer `create_before_destroy`, `lifecycle.ignore_changes` (with a comment why), or a migration plan.
- An unexpected large diff usually means wrong workspace, wrong var-file, provider upgrade or drift. Stop and explain before applying.

**State**
- Remote backend with locking and encryption (S3 + native lockfile `use_lockfile = true` or DynamoDB, azurerm blob lease, GCS, Terraform Cloud). State is never committed; `*.tfstate*` and `.terraform/` are in `.gitignore`.
- One state per environment per component; keep blast radius small (network, data, app layers separate). Read across with `terraform_remote_state` or data sources, sparingly.
- State surgery is declarative and reviewed: `import {}` blocks to adopt, `moved {}` to rename/refactor, `removed {}` (`lifecycle { destroy = false }`) to forget without deleting. CLI `state mv/rm`, `import`, `taint`, `force-unlock` only with approval and after `terraform state pull > backup-<date>.tfstate`. Never hand-edit state JSON.
- `force-unlock` only after confirming no other run holds the lock (check CI runs, the lock's `Who`/`Created`).

**Versions and modules**
- Pin `required_version` and every provider in `required_providers` with `~>` on minor; commit `.terraform.lock.hcl`. Add platforms CI needs: `terraform providers lock -platform=linux_amd64 -platform=darwin_arm64`.
- Module sources pinned to a tag or registry version (`?ref=v1.4.2`, `version = "~> 1.4"`); never a branch. Upgrade providers/modules in their own PR with the plan attached.
- Modules: one purpose, typed variables with `validation` blocks and descriptions, outputs for what callers need, no provider blocks inside reusable modules, no hard-coded regions/accounts.

**Environments**
| Use | When |
|---|---|
| Directory per env (`envs/prod/`) + shared modules | environments differ in shape, separate credentials/backends, most teams |
| Workspaces | identical stacks differing only by variables, same backend and credentials |
| Terragrunt / stacks | many accounts/regions with DRY backend config |
- Never let one root module apply to prod and non-prod with only a variable flip and shared credentials.

**Security and hygiene**
- No secrets in `.tf`, `*.tfvars`, outputs or defaults. State stores every attribute in plain text, including generated passwords: prefer resources that write straight to a secret manager, `ephemeral` resources/write-only args where supported, mark variables/outputs `sensitive = true`, and restrict state bucket access.
- Least-privilege IAM: specific actions and ARNs; justify every `*`. Private by default: no `0.0.0.0/0` ingress, public buckets or public databases without explicit approval.
- Encryption at rest and in transit, access logging, backups and deletion protection on stateful resources; `prevent_destroy` on critical ones.
- Tags via provider `default_tags` (or equivalent): `owner`, `environment`, `cost-center`, `managed-by`, `repo`. Policy-check that they exist.
- Drift: scheduled `plan -detailed-exitcode` (exit 2 = drift) that alerts the owner. Reconcile by changing code or importing, never by re-clicking in the console.

## Safety gates
- Ask before any `apply`, `destroy`, `import`, `state mv/rm`, `taint`/`-replace`, `force-unlock`, workspace delete or backend migration. State the environment, counts (add/change/destroy) and every destroy/replace first.
- Production applies need the user's explicit "yes" in this conversation and go through the pipeline where one exists.
- Stop and ask if the plan destroys or replaces any data store, DNS, identity/IAM root, network core (VPC/VNet, peering, gateways) or KMS key, or if counts differ from what the change should produce.
- Never disable policy checks, deletion protection or state locking to get a change through.

## References
- Read `references/terraform.md` when the repo uses Terraform or OpenTofu and you need state surgery recipes, refactoring with `moved`, import generation, or CI apply patterns.
- Read `references/other-tools.md` when the repo uses Bicep/ARM, Pulumi or CloudFormation/CDK.

## Output
```
Target:      <env> / <account|subscription|project> / <region>, tool+version
Change:      <what and why>
Plan:        +<n> ~<n> -<n>  (replace: <n>)
Destroy/replace: <resource> - <forcing attribute> - <data/outage risk> - <mitigation>
High risk:   IAM / network / data-store diffs, blast radius
Checks:      fmt/validate/tflint/trivy|checkov/conftest - pass|fail (findings)
Apply:       saved plan <file/run id>, approved by <who>, via <pipeline|local>
Rollback:    <revert PR + re-apply | restore from backup | n/a> 
Follow-ups:  drift, tags, pinned versions, docs
```
