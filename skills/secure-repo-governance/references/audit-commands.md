# Read-only audit commands (`gh api`)

Set `R=<owner>/<repo>`, `ORG=<org>`, `B=<branch>`. All commands are GET. A 404 can mean "not configured" or "insufficient permission" - check `gh api repos/$R --jq .permissions` first.

## Repository basics
```bash
gh api repos/$R --jq '{default_branch, visibility, archived, permissions,
  allow_squash_merge, allow_merge_commit, allow_rebase_merge, allow_auto_merge,
  delete_branch_on_merge, squash_merge_commit_title, web_commit_signoff_required,
  security_and_analysis}'
```
`security_and_analysis` shows `advanced_security`, `secret_scanning`, `secret_scanning_push_protection`, `secret_scanning_non_provider_patterns`, `secret_scanning_validity_checks`, `dependabot_security_updates` (fields vary by plan; absent = not available or not visible).

## Branch and tag rules
```bash
gh api repos/$R/rules/branches/$B                      # effective rules from all repo + org rulesets
gh api repos/$R/rulesets --jq '.[] | {id, name, target, enforcement, source_type}'
gh api repos/$R/rulesets/<id>                          # conditions, rules, bypass_actors
gh api orgs/$ORG/rulesets                              # org-level rulesets
gh api repos/$R/branches/$B/protection                 # classic protection (404 = none or no access)
gh api repos/$R/branches/$B/protection/required_signatures
```
Key ruleset rule types: `pull_request` (required_approving_review_count, dismiss_stale_reviews_on_push, require_code_owner_review, require_last_push_approval, required_review_thread_resolution, allowed_merge_methods), `required_status_checks` (contexts + integration_id, strict policy), `non_fast_forward`, `deletion`, `required_linear_history`, `required_signatures`, `merge_queue`, `creation`/`update` (tags), `code_scanning`, `file_path_restriction`, `max_file_size` (push rules).

## Required checks vs real check names
```bash
SHA=$(gh api repos/$R/commits/$B --jq .sha)
gh api repos/$R/commits/$SHA/check-runs --paginate --jq '.check_runs[] | {name, app: .app.slug}'
gh api repos/$R/commits/$SHA/status --jq '.statuses[].context'   # legacy commit statuses
```
Every required context must appear in this list for a typical PR head; mismatches block merges.

## Actions settings
```bash
gh api repos/$R/actions/permissions                     # enabled, allowed_actions, sha_pinning_required (if supported)
gh api repos/$R/actions/permissions/selected-actions    # when allowed_actions = selected
gh api repos/$R/actions/permissions/workflow            # default_workflow_permissions, can_approve_pull_request_reviews
gh api repos/$R/actions/permissions/fork-pr-contributor-approval   # if available on the plan
gh api repos/$R/actions/runners --jq '.runners[] | {name, labels: [.labels[].name]}'
gh api repos/$R/actions/secrets --jq '.secrets[].name'  # names only, values are never returned
gh api repos/$R/actions/variables --jq '.variables[].name'
```

## Security features
```bash
gh api repos/$R/vulnerability-alerts -i | head -1       # 204 = Dependabot alerts on, 404 = off
gh api repos/$R/automated-security-fixes                # {enabled, paused}
gh api repos/$R/private-vulnerability-reporting         # {enabled}
gh api repos/$R/code-scanning/default-setup             # state, languages, query_suite
gh api "repos/$R/code-scanning/alerts?state=open" --paginate --jq '.[] | {number, rule: .rule.id, severity: .rule.security_severity_level}'
gh api "repos/$R/secret-scanning/alerts?state=open&hide_secret=true" --paginate --jq '.[] | {number, secret_type, created_at}'   # never print the secret field
gh api "repos/$R/dependabot/alerts?state=open" --paginate --jq '.[] | {number, pkg: .dependency.package.name, sev: .security_advisory.severity}'
```
With `-f` parameters `gh api` switches to POST; add `-X GET` or put parameters in the URL as above.

## Environments
```bash
gh api repos/$R/environments --jq '.environments[] | {name, protection_rules, deployment_branch_policy}'
gh api repos/$R/environments/<env>/deployment-branch-policies
gh api repos/$R/environments/<env>/secrets --jq '.secrets[].name'
```
Check: `required_reviewers` present, `prevent_self_review: true`, `wait_timer`, branch policy not `null` (null = any branch can deploy), `can_admins_bypass`.

## Files
```bash
for f in CODEOWNERS .github/CODEOWNERS docs/CODEOWNERS SECURITY.md .github/SECURITY.md .github/dependabot.yml .github/pull_request_template.md LICENSE; do
  gh api repos/$R/contents/$f --silent 2>/dev/null && echo "present: $f" || echo "missing: $f"; done
gh api repos/$R/codeowners/errors --jq '.errors[] | {line, message}'
```
`SECURITY.md` may also be inherited from the org `.github` repository.

## Access
```bash
gh api repos/$R/collaborators -f affiliation=direct -X GET --paginate --jq '.[] | {login, role_name}'
gh api repos/$R/collaborators -f affiliation=outside -X GET --paginate --jq '.[].login'
gh api repos/$R/teams --jq '.[] | {slug, permission}'
gh api repos/$R/keys --jq '.[] | {title, read_only, created_at}'          # deploy keys
gh api repos/$R/hooks --jq '.[] | {id, active, events, url: .config.url, insecure_ssl: .config.insecure_ssl}'
# Installed apps with access to this repo: check org installations below
```

## Organisation (needs org owner or `read:org`/`admin:org`)
```bash
gh api orgs/$ORG --jq '{two_factor_requirement_enabled, default_repository_permission,
  members_can_create_public_repositories, members_can_create_private_repositories,
  members_can_fork_private_repositories, web_commit_signoff_required,
  members_can_delete_repositories, members_can_change_repo_visibility}'
gh api orgs/$ORG/actions/permissions
gh api orgs/$ORG/actions/permissions/workflow
gh api orgs/$ORG/installations --jq '.installations[] | {app_slug, repository_selection, permissions}'
gh api orgs/$ORG/outside_collaborators --paginate --jq '.[].login'
gh api orgs/$ORG/code-security/configurations --jq '.[] | {id, name, enforcement}'
gh api orgs/$ORG/rulesets
```
Fine-grained PAT approval policy and audit-log streaming are visible in the org settings UI (and audit-log API on Enterprise); record them as "manual check" if not readable.

## Bulk audit across an org
```bash
gh repo list $ORG --limit 1000 --no-archived --json nameWithOwner --jq '.[].nameWithOwner' |
  while read r; do echo "== $r"; gh api repos/$r/rules/branches/$(gh api repos/$r --jq .default_branch) --jq 'map(.type) | join(",")'; done
```
Respect rate limits (`gh api rate_limit`); for large orgs prefer GraphQL or the platform's security overview.
