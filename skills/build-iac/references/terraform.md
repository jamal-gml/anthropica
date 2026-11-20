# Terraform / OpenTofu recipes

OpenTofu is CLI-compatible: swap `terraform` for `tofu`. Differences that matter: OpenTofu has native state encryption (`encryption {}` block) and its own registry; lockfile hashes differ between the two registries, so do not mix binaries on one state.

## Backends (with locking)
```hcl
terraform {
  required_version = "~> 1.9"
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.60" }
  }
  backend "s3" {
    bucket       = "<state-bucket>"
    key          = "network/prod/terraform.tfstate"
    region       = "<region>"
    encrypt      = true
    use_lockfile = true          # S3-native locking (TF >= 1.10); else dynamodb_table
  }
}
```
- azurerm: `resource_group_name`, `storage_account_name`, `container_name`, `key`; locking via blob lease (automatic). Use `use_azuread_auth = true`.
- gcs: `bucket`, `prefix`; locking automatic. Enable object versioning on the bucket.
- Enable versioning on every state bucket/container so a bad write can be rolled back.
- Backend blocks cannot use variables; pass per-env values with `terraform init -backend-config=envs/prod.backend.hcl`.
- Migrating backends: `terraform init -migrate-state` after a `state pull` backup. Ask first.

## Import existing resources
```hcl
import {
  to = aws_s3_bucket.logs
  id = "<bucket-name>"
}
```
```bash
terraform plan -generate-config-out=generated.tf   # writes HCL for import blocks without config
```
- Clean the generated HCL (remove computed/default attributes), re-plan until it shows `1 to import, 0 to add, 0 to change, 0 to destroy`.
- `for_each` on import blocks works (TF >= 1.7) for bulk adoption.

## Refactor without destroy
```hcl
moved {
  from = aws_instance.web
  to   = module.web.aws_instance.this
}
moved {
  from = aws_subnet.private[0]
  to   = aws_subnet.private["eu-a"]     # count -> for_each
}
```
- Keep `moved` blocks until every state that uses the module has applied them (for shared modules: at least one major version).
- Renaming a resource or wrapping it in a module without `moved` = destroy + create. Always check the plan says `has moved to`.

## Stop managing without deleting
```hcl
removed {
  from = aws_s3_bucket.legacy
  lifecycle { destroy = false }
}
```
Prefer this over `terraform state rm` because it is reviewed in a PR and works in CI.

## CLI state surgery (approval required)
```bash
terraform state pull > backup-$(date +%Y%m%d%H%M).tfstate   # always first
terraform state list | grep <pattern>
terraform state show <addr>
terraform state mv <from> <to>
terraform state rm <addr>
terraform apply -replace=<addr>          # replaces 'taint'
terraform force-unlock <LOCK_ID>         # only after confirming no run holds it
```
Restoring: `terraform state push backup.tfstate` (serial must be higher, or `-force`; ask first).

## Targeting
`-target` is for recovery only (break a cycle, unblock a failed apply). It skips dependencies; always follow with a full plan showing no changes.

## Lifecycle patterns
| Need | Pattern |
|---|---|
| Never delete | `lifecycle { prevent_destroy = true }` + cloud deletion protection |
| Zero-downtime replace | `create_before_destroy = true` (name collisions: use `name_prefix`) |
| Externally managed field (autoscaler desired count, tags set by another tool) | `ignore_changes = [desired_count]` with a comment |
| Replace when another thing changes | `replace_triggered_by = [terraform_data.x]` |
| Validate inputs | `variable { validation { condition = ... error_message = ... } }` |
| Assert after apply | `check {}` blocks / `postcondition` |

## Common failure modes
| Symptom | Cause | Fix |
|---|---|---|
| `Error acquiring the state lock` | concurrent run or crashed run | wait / find the run; `force-unlock` with approval |
| Plan wants to replace everything | wrong workspace, var-file or backend key | `terraform workspace show`, check `-var-file`, backend key |
| `Inconsistent dependency lock file` | provider changed without lock update | `terraform init -upgrade` in a dedicated PR |
| `Provider produced inconsistent result after apply` | provider bug / eventual consistency | re-plan; pin provider; report upstream |
| Perpetual diff on one attribute | API normalises value (JSON policy order, case) | `jsonencode`/`aws_iam_policy_document`, normalise, or `ignore_changes` |
| `count`/`for_each` "cannot be determined until apply" | key depends on unknown value | key on known inputs (names), not generated IDs |
| Cycle error | mutual references | split resource (e.g. SG + separate rule resources) |

## CI pattern
1. PR: `fmt -check`, `validate`, `tflint`, `trivy config`, `plan -out` per changed root, upload `tfplan` + post summary to PR, `conftest` on `tfplan.json`.
2. Merge: apply the exact plan artifact from the PR run (or re-plan and fail if it differs), behind an environment approval for production.
3. Auth via OIDC federation to the cloud, no long-lived keys. Separate plan (read-only) and apply roles.
4. Nightly: `plan -detailed-exitcode` per root; exit 2 opens a drift ticket.
Pipeline mechanics live in `ship-ci-cd` and `ship-github-actions`.

## Testing
- `terraform test` (`*.tftest.hcl`) with `command = plan` for fast assertions on module logic; `command = apply` in a sandbox account for integration.
- Terratest for complex end-to-end checks. Always destroy test infrastructure in a `defer`/cleanup step.

## Useful policy checks (OPA/Conftest on plan JSON)
- Deny `resource_changes[_].change.actions` containing `"delete"` for types in a protected list unless an override label is set.
- Require tags on every taggable resource.
- Deny security group rules with `0.0.0.0/0` on ports other than 80/443.
- Deny unencrypted storage, public ACLs, IAM `"Action": "*"`.
