---
name: devops-ci-cd
description: Design CI/CD pipelines and deployment flows (stages, environments, promotion gates, deployment strategies, release automation) on any CI platform. Use when planning or restructuring a pipeline or deployment process; for GitHub workflow files specifically, also use github-actions.
---

# DevOps and CI/CD

1. Inspect existing workflows, branch strategy, environments and deployment targets.
2. Identify trigger conditions, permissions, secrets and artifact flow.
3. Use least-privilege workflow permissions.
4. Pin third-party actions to verified immutable commit SHAs where policy requires.
5. Separate build, test, security scan and deployment stages.
6. Protect production deployment with environment approvals and explicit gates.
7. Validate workflow syntax and run available local checks.
8. Document required repository secrets by name only, never their values.

Do not deploy, push or change repository settings without explicit approval.
