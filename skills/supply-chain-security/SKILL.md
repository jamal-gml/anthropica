---
name: supply-chain-security
description: Assess and harden software supply-chain security (dependency vulnerabilities, lockfiles, pinning, SBOMs, artifact signing and provenance, CI/CD hardening, secret exposure). Use when auditing dependencies, triaging Dependabot or scanner alerts, adding a new dependency, or hardening build and release pipelines.
---

# Supply-Chain Security

## Dependencies
- Commit lockfiles and install from them in CI (`npm ci`, `pip install --require-hashes` or `uv sync --frozen`, `go mod verify`).
- Before adding a dependency, check maintenance activity, popularity, licence compatibility, known vulnerabilities and typosquatting (exact package name). Prefer the standard library or an existing dependency.
- Scan: `npm audit`, `pip-audit`, `osv-scanner`, `trivy fs .`, `govulncheck`, or the repository's configured scanner.
- Keep updates flowing through Dependabot or Renovate. Group low-risk updates and review major versions individually.

## Triaging an alert
1. Is the vulnerable package actually used, and is the vulnerable function reachable from our code?
2. Is it a runtime dependency or dev/build-only? Is the exposure internet-facing?
3. Fix by upgrading to the patched version through the update bot where possible. If no fix exists, mitigate or document the risk with an owner and review date.
4. Never dismiss an alert without a recorded reason.

## Build and CI
- Pin GitHub Actions and container base images to digests or commit SHAs.
- Least-privilege CI tokens (`permissions:` per job). OIDC for cloud access; no long-lived keys.
- Protect release workflows: environment approvals, restricted branches and tags, no secrets exposed to pull requests from forks.
- Reproducible builds where feasible: pinned toolchains, no network fetches of unpinned scripts (`curl | bash`).

## Artifacts
- Generate an SBOM (CycloneDX or SPDX) for releases (`syft`, or the build tool's plugin).
- Sign artifacts and images (Sigstore `cosign`) and publish build provenance (SLSA, GitHub artifact attestations).
- Consumers verify signatures and provenance before deploying.

## Secrets
- Secret scanning with push protection on. Scan history (`gitleaks`, `trufflehog`) when auditing.
- A leaked secret is compromised even after it's removed from the code: rotate it first, then clean up and review access logs.

## Output
Findings table (component, issue, severity, reachable or exposed, fix, owner), the changes made or proposed, and residual risks. Never print secret values.
