---
name: release-readiness
description: Prepare a release readiness review covering changes, tests, migrations, security, rollback and deployment approvals. Use before tagging, promoting or deploying a release.
---

# Release Readiness

Review:
1. Scope and change summary
2. Build, tests, lint and type checks
3. Dependency and security findings
4. Configuration and secret requirements
5. Database migrations and compatibility
6. Observability and operational impact
7. Rollback or recovery plan
8. Deployment approvals and ownership

Output a go/no-go recommendation with evidence and explicit blockers.
A readiness review does not authorise deployment.
