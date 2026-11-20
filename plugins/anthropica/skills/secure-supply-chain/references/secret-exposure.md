# Leaked-secret response

A secret that reached a commit, log, artifact, issue, chat or public package is compromised, even if deleted seconds later. Removal from history does not un-leak it; rotation does.

## Order of operations
1. **Confirm without spreading.** Identify secret type, owning system and scope from the scanner output (location, rule, redacted fingerprint). Do not print or copy the value. Check whether it is live only via a non-destructive method the owner approves (e.g. provider's "token info" endpoint or the scanner's verification flag).
2. **Contain - rotate or revoke first.** Ask the owner to revoke/rotate in the issuing system (cloud IAM key, API token, DB password, signing key, OAuth client secret, webhook secret). Update consumers through the secret manager, not by committing the new value. Rotation is outward-facing: state the action and get approval.
3. **Investigate use.** Review the provider's audit logs for the exposure window (first commit time to revocation): unusual IPs, regions, API calls, new credentials or resources created, data access. Escalate to `operate-incident` if any suspicious use appears or the scope is production/customer data.
4. **Clean up the source.** Remove from the current tree, replace with a secret-manager reference or env var. Rewriting history (`git filter-repo --replace-text`, BFG) needs coordination: it requires force-push to shared branches, invalidates clones and forks, and does not remove cached copies (forks, PR refs, CI logs, package registries, search caches). Ask before doing it; often rotation + leaving history is acceptable.
5. **Purge secondary copies.** CI logs, build artifacts, container image layers (rebuild without the secret; `docker history` exposes `ARG`/`ENV`), published package versions (deprecate/unpublish per registry policy), issue/PR comments, chat, tickets.
6. **Prevent recurrence.** Enable push protection, add pre-commit scanning (`gitleaks protect --staged` or `gitleaks git --pre-commit --staged` on v8.19+), `.gitignore` for `.env*` and key files, move to OIDC/short-lived credentials so there is nothing long-lived to leak.
7. **Record.** What leaked (type, not value), where, window, rotation time, log review result, follow-ups.

## Scanning commands
- Working tree + history: `gitleaks git --redact -v` (older versions: `gitleaks detect --redact`), `trufflehog git file://. --only-verified --no-update`.
- Directory or artifact: `gitleaks dir <path> --redact`, `trufflehog filesystem <path>`.
- Container image: `trivy image --scanners secret <ref>`.
- Org-wide: use the platform's secret scanning alerts; validity checks where supported.

## Triage of scanner hits
- True positive: real credential, any environment. Rotate even if "only dev" - dev keys often reach shared resources.
- Test fixture: clearly fake value, documented as such; add a narrowly scoped allowlist entry (path + rule) with a comment, never a global rule disable.
- Revoked/expired: confirm revocation in the issuing system before closing; record evidence.
- Dismissing a push-protection block requires a reason; "will fix later" is not one.

## Never
- Paste the secret into a ticket, PR, chat or this report.
- Test a leaked third-party credential against the provider beyond the owner-approved validity check.
- Rely on making a public repo private as remediation.
