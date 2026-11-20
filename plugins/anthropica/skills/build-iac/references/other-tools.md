# Bicep/ARM, Pulumi, CloudFormation/CDK

Same principles as Terraform: preview to a reviewable diff, gate, apply exactly what was reviewed, never hand-edit state.

## Bicep / ARM
| Step | Command |
|---|---|
| Lint/build | `az bicep build --file main.bicep` (warnings from `bicepconfig.json` linter rules) |
| Preview | `az deployment group what-if -g <rg> -f main.bicep -p @prod.bicepparam` |
| Deploy | `az deployment group create -g <rg> -f main.bicep -p @prod.bicepparam --confirm-with-what-if` |
| Scopes | `group`, `sub`, `mg`, `tenant` - match the template's `targetScope` |
- No state file: Azure is the state. `what-if` can be noisy (shows `Ignore`/`NoChange` and false positives on some resource types); read `Delete` and `Modify` lines carefully.
- Deployment mode: default `Incremental`. `Complete` mode deletes resources in the group that are not in the template - ask before using it; prefer deployment stacks (`az stack group create --action-on-unmanage detachAll|deleteAll`) for managed deletion.
- Pin module versions from registries (`br:<registry>/bicep/modules/x:1.2.0`) and API versions on every resource.
- Secrets: `@secure()` params, `getSecret()` from Key Vault references in `.bicepparam`; never output secrets.
- Locks: `Microsoft.Authorization/locks` `CanNotDelete` on production data resources.
- Drift: periodic `what-if`; Azure Policy for guardrails (deny/audit/modify).

## Pulumi
| Step | Command |
|---|---|
| Select stack | `pulumi stack select <org>/<project>/<stack>` (verify before anything else) |
| Preview | `pulumi preview --diff` (save with `--save-plan plan.json`) |
| Apply reviewed plan | `pulumi up --plan plan.json` |
| Drift | `pulumi refresh --preview-only` / `pulumi preview --refresh` |
| Import | `pulumi import <type> <name> <id>` (generates code) |
- State backend: Pulumi Cloud or self-managed (`pulumi login s3://...`, `azblob://`, `gs://`) with locking. Never commit `.pulumi/`.
- Secrets: `pulumi config set --secret`; encrypted in state with the stack's secrets provider (use a cloud KMS provider, not passphrase, for shared stacks).
- Renames: use `aliases` resource option, otherwise rename = replace.
- Protect: `protect: true` resource option; `retainOnDelete: true` to forget without deleting.
- Pin package versions in the language lockfile (package-lock, poetry.lock, go.sum).
- Replacement markers in preview: `+-` replace, `++` create-before-replace, `--` delete-before-replace.
- State surgery: `pulumi state delete|rename|unprotect` and `pulumi stack export/import` - backup with `pulumi stack export > backup.json` first; ask before running.

## CloudFormation
| Step | Command |
|---|---|
| Lint | `cfn-lint template.yaml`; policy with `cfn-guard validate` |
| Preview | `aws cloudformation create-change-set --change-set-name <n> ...` then `describe-change-set` |
| Apply | `aws cloudformation execute-change-set` |
| Drift | `detect-stack-drift` then `describe-stack-resource-drifts` |
- Always deploy via change sets; read `Replacement: True|Conditional` per resource.
- `DeletionPolicy: Retain|Snapshot` and `UpdateReplacePolicy` on every stateful resource. Stack termination protection on production stacks.
- Stack policies to block replace/delete of critical resources during updates.
- Failed update -> automatic rollback; `UPDATE_ROLLBACK_FAILED` needs `continue-update-rollback` (optionally `--resources-to-skip`) - ask first.
- Import existing resources with an `IMPORT` change set (resource must have `DeletionPolicy` set in the template).
- Secrets: dynamic references `{{resolve:secretsmanager:...}}` / `{{resolve:ssm-secure:...}}`; `NoEcho` params are still visible in some APIs, avoid them for secrets.

## AWS CDK
| Step | Command |
|---|---|
| Synth | `cdk synth` (fails on construct errors; run cdk-nag as an Aspect) |
| Preview | `cdk diff <stack>` |
| Deploy | `cdk deploy <stack> --require-approval broadening` (never `never` locally) |
- Logical IDs come from the construct path: moving or renaming a construct = replacement. Keep IDs stable or use `overrideLogicalId`; check `cdk diff` for `[-]`/`[+]` pairs on stateful resources.
- `removalPolicy: RETAIN` on data resources (default for many but not all).
- Pin `aws-cdk-lib` and the CLI to the same version; commit the lockfile and `cdk.context.json` (otherwise lookups drift between machines).
- Bootstrap per account/region once (`cdk bootstrap`); use the bootstrap roles with least privilege and a trust policy limited to the pipeline.
- Testing: assertions with `Template.fromStack(...)` and snapshot tests of the synthesized template.
