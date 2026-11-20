# Templates for proposed changes

Adapt names and counts to the team's standard. Show these to the user and apply only after approval.

## Branch ruleset (default branch)
Save as `ruleset-default.json`, then `gh api -X POST repos/<owner>/<repo>/rulesets --input ruleset-default.json`. Start with `"enforcement": "evaluate"` where the plan supports it, then switch to `active`. Rollback: `gh api -X DELETE repos/<owner>/<repo>/rulesets/<id>` or set `"enforcement": "disabled"`.
```json
{
  "name": "default-branch-protection",
  "target": "branch",
  "enforcement": "active",
  "conditions": { "ref_name": { "include": ["~DEFAULT_BRANCH"], "exclude": [] } },
  "bypass_actors": [],
  "rules": [
    { "type": "deletion" },
    { "type": "non_fast_forward" },
    { "type": "required_linear_history" },
    { "type": "pull_request", "parameters": {
        "required_approving_review_count": 1,
        "dismiss_stale_reviews_on_push": true,
        "require_code_owner_review": true,
        "require_last_push_approval": true,
        "required_review_thread_resolution": true,
        "allowed_merge_methods": ["squash"] } },
    { "type": "required_status_checks", "parameters": {
        "strict_required_status_checks_policy": true,
        "required_status_checks": [
          { "context": "build", "integration_id": 15368 },
          { "context": "lint", "integration_id": 15368 } ] } }
  ]
}
```
- `integration_id` 15368 is the GitHub Actions app on github.com; confirm with the check-runs `app.id` from `audit-commands.md`.
- `context` must equal the exact check-run name. Remove `required_linear_history` and adjust `allowed_merge_methods` if merge commits are needed for promotions.
- Break-glass bypass: `{ "actor_id": <team_id>, "actor_type": "Team", "bypass_mode": "pull_request" }`.
- Add `{ "type": "required_signatures" }` only after confirming every committer and bot signs.

## Tag ruleset (releases)
```json
{
  "name": "release-tags",
  "target": "tag",
  "enforcement": "active",
  "conditions": { "ref_name": { "include": ["refs/tags/v*"], "exclude": [] } },
  "bypass_actors": [ { "actor_id": 0, "actor_type": "<Team or Integration of release role>", "bypass_mode": "always" } ],
  "rules": [ { "type": "creation" }, { "type": "update" }, { "type": "deletion" }, { "type": "non_fast_forward" } ]
}
```
Replace the placeholder bypass actor with the release team or app; without one, nobody can create release tags.

## Actions least privilege
```bash
gh api -X PUT repos/<owner>/<repo>/actions/permissions/workflow \
  -f default_workflow_permissions=read -F can_approve_pull_request_reviews=false
gh api -X PUT repos/<owner>/<repo>/actions/permissions \
  -F enabled=true -f allowed_actions=selected
gh api -X PUT repos/<owner>/<repo>/actions/permissions/selected-actions \
  -F github_owned_allowed=true -F verified_allowed=true -f 'patterns_allowed[]=<owner>/*'
```
Before switching to read-only default, find workflows without a `permissions:` block that need writes (`rg --files-without-match "permissions:" .github/workflows`) and add explicit job permissions first, or they will fail.

## Security features
```bash
gh api -X PATCH repos/<owner>/<repo> --input - <<'EOF'
{ "security_and_analysis": {
    "secret_scanning": { "status": "enabled" },
    "secret_scanning_push_protection": { "status": "enabled" },
    "dependabot_security_updates": { "status": "enabled" } } }
EOF
gh api -X PUT repos/<owner>/<repo>/vulnerability-alerts
gh api -X PUT repos/<owner>/<repo>/private-vulnerability-reporting
gh api -X PATCH repos/<owner>/<repo>/code-scanning/default-setup -f state=configured
```
Private repos may need Advanced Security / Secret Protection / Code Security licences; the API returns 403/422 otherwise.

## Production environment
```bash
gh api -X PUT repos/<owner>/<repo>/environments/production --input - <<'EOF'
{ "wait_timer": 0, "prevent_self_review": true, "can_admins_bypass": false,
  "reviewers": [ { "type": "Team", "id": <team_id> } ],
  "deployment_branch_policy": { "protected_branches": false, "custom_branch_policies": true } }
EOF
gh api -X POST repos/<owner>/<repo>/environments/production/deployment-branch-policies -f name='main' -f type=branch
gh api -X POST repos/<owner>/<repo>/environments/production/deployment-branch-policies -f name='v*' -f type=tag
```

## CODEOWNERS (`.github/CODEOWNERS`)
```
# Last matching pattern wins. Owners need write access.
*                       @<org>/<core-team>
/.github/workflows/     @<org>/<platform-team> @<org>/<security-team>
/.github/CODEOWNERS     @<org>/<security-team>
/infra/                 @<org>/<platform-team>
/src/auth/              @<org>/<security-team>
```

## SECURITY.md
```markdown
# Security policy
## Supported versions
| Version | Supported |
|---|---|
| latest major | yes |
## Reporting a vulnerability
Please do not open public issues. Report privately via this repository's "Report a vulnerability" (private vulnerability reporting) or to <security contact>.
Include affected version, reproduction steps and impact. We acknowledge within <N> business days and aim to fix critical issues within <N> days.
We coordinate disclosure and credit reporters unless asked not to.
```

## `.github/dependabot.yml`
```yaml
version: 2
updates:
  - package-ecosystem: github-actions
    directory: /
    schedule: { interval: weekly }
    groups: { actions: { patterns: ["*"] } }
  - package-ecosystem: npm          # repeat per ecosystem and directory (use `directories:` for monorepos)
    directory: /
    schedule: { interval: weekly }
    cooldown: { default-days: 3 }   # delay non-security updates; security updates are not delayed
    groups:
      minor-and-patch: { update-types: ["minor", "patch"] }
    open-pull-requests-limit: 10
  - package-ecosystem: docker
    directory: /
    schedule: { interval: weekly }
```
Every manifest directory must be covered; check with `rg --files -g 'package.json' -g 'go.mod' -g 'pyproject.toml' -g 'Dockerfile' -g '*.csproj'` excluding vendored paths.
