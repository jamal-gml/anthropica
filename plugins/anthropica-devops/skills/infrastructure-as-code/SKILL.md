---
name: infrastructure-as-code
description: Write, review and apply infrastructure as code (Terraform/OpenTofu, Bicep, Pulumi, CloudFormation) safely, with plan-before-apply, remote state, modules, policy checks and drift control. Use when changing cloud resources, networking, IAM or IaC modules, or reviewing a plan.
---

# Infrastructure as Code

## Discovery
1. Identify the tool and version, providers, state backend, workspaces or stacks, and the environment layout.
2. Read existing modules and naming and tagging conventions before writing new code.
3. Confirm which environment and account or subscription a change targets. Never assume production is safe to touch.

## Writing changes
- Reuse existing modules. Keep modules small, versioned and input-validated.
- Pin tool and provider versions (`required_version`, `required_providers`, lockfile committed).
- No hard-coded secrets, account IDs or IPs; use variables, data sources and a secret manager.
- Tag every resource with owner, environment, cost centre and the managing repository.
- Least-privilege IAM: specific actions and resources, no wildcards without justification.
- Private by default: no public storage buckets, open security groups (`0.0.0.0/0`) or public databases without explicit approval.
- Enable encryption at rest and in transit, logging, and backups where the resource supports them.
- Protect critical resources (`prevent_destroy`, deletion protection, resource locks).

## Validate before apply
```bash
terraform fmt -check && terraform validate
tflint                      # if configured
checkov -d . / tfsec / trivy config .   # policy and security scan, if available
terraform plan -out=tfplan
```
- Read the plan in full. Call out every **destroy** and **replace** and explain why it happens.
- Changes to state, IAM, networking or data stores are high risk: summarise the blast radius.

## Apply
- Apply only a reviewed saved plan, through the pipeline where one exists.
- Production applies require the user's explicit approval in the current conversation.
- Never run `state rm`, `import`, `taint`, `force-unlock` or `destroy` without explicit approval and a stated reason.

## Drift and state
- Detect drift with a scheduled `plan`. Reconcile by updating code, not by clicking in the console.
- Never edit state files by hand or commit them. Use remote state with locking and encryption.

## Output
Target environment, resources added, changed and destroyed (from the plan), risk summary, validation results, and the apply and rollback approach.
