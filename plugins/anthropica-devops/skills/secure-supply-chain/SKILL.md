---
name: secure-supply-chain
description: Harden the software supply chain - lockfiles, dependency audits with reachability triage, pinning actions by SHA and images by digest, update bots, typosquatting and dependency confusion, SBOM, signing and provenance, CI hardening, leaked-secret response. Use when triaging Dependabot or scanner alerts, adding a dependency, hardening build or release pipelines, or after a secret leak. Not for repo settings (use `secure-repo-governance`).
---

# Supply-Chain Security

Make what you build come only from what you meant to depend on, built by a pipeline that cannot be hijacked, with evidence a consumer can verify.

## Procedure
1. **Inventory.** List ecosystems from manifests and lockfiles (`package.json`/lock, `pyproject.toml`/`uv.lock`/`requirements*.txt`, `go.mod`/`go.sum`, `Cargo.lock`, `pom.xml`/`build.gradle*`, `*.csproj`/`packages.lock.json`, `Gemfile.lock`, `composer.lock`), Dockerfiles, `.github/workflows/*`, IaC modules, install scripts (`curl | sh`), and registry config (`.npmrc`, `pip.conf`, `nuget.config`, `settings.xml`, `GOPRIVATE`).
2. **Lockfiles and reproducible installs.** Each ecosystem has a committed lockfile and CI installs in frozen mode (`npm ci`, `pnpm install --frozen-lockfile`, `uv sync --locked`, `pip install --require-hashes`, `cargo build --locked`, `dotnet restore --locked-mode`, `bundle install` with `BUNDLE_FROZEN=true`, `go mod verify` with `-mod=readonly`). Read `references/ecosystems.md` for per-ecosystem lock, audit and registry commands.
3. **Audit dependencies.** Run the native tool plus one cross-ecosystem scanner: `osv-scanner scan source -r .` (v1: `osv-scanner -r .`) or `trivy fs --scanners vuln,secret,misconfig .`; `govulncheck ./...` for Go (has call-graph reachability); `trivy image <ref>` or `grype <ref>` for containers.
4. **Triage each alert** (rules below): reachable? runtime or dev-only? exposed? fix available? Upgrade through the update bot; otherwise mitigate or accept with owner and review date.
5. **Vet new dependencies** before adding: exact name (typosquats differ by one char, hyphen/underscore or scope), publisher and repo match, maintenance and release history, install scripts, transitive footprint, licence, OpenSSF Scorecard (`scorecard --repo=github.com/<owner>/<repo>`). Prefer the standard library or an existing dependency.
6. **Dependency confusion.** Internal package names are scoped/namespaced and resolved only from the private registry (npm `@scope:registry=`, NuGet package source mapping, a single `--index-url` proxy for pip, never `--extra-index-url` for internal names, `GOPRIVATE`/`GONOSUMDB`). Reserve internal names on public registries if the ecosystem has no scopes.
7. **Pin.** Actions by full 40-char commit SHA with a version comment; base images by `@sha256:` digest; tool downloads by version plus checksum; Terraform providers via `.terraform.lock.hcl`. Let Dependabot/Renovate bump the pins.
8. **Keep updates flowing.** Dependabot or Renovate covers every ecosystem including `github-actions` and `docker`; group minor/patch, review majors individually; set a release-age delay (Dependabot `cooldown`, Renovate `minimumReleaseAge`) of a few days for non-security updates to dodge freshly hijacked releases.
9. **Harden CI/CD.** Read `references/ci-hardening.md` when reviewing workflows or release pipelines (token permissions, `pull_request_target`, script injection, OIDC, runners, caches). Workflow authoring itself belongs to `ship-github-actions`.
10. **SBOM, signing, provenance** for anything you release. Read `references/sbom-provenance.md` when publishing packages, images or binaries.
11. **Secrets.** Scan working tree and history: `gitleaks git --redact` (older: `gitleaks detect --redact`) or `trufflehog git file://. --only-verified`. On any real exposure follow `references/secret-exposure.md` - rotate first.
12. **Report** and propose changes as diffs; dependency bumps go through the update bot where it supports the ecosystem.

## Rules
- **Alert triage order:** (1) Is the package in the shipped artifact (runtime) or only dev/test/build? (2) Is the vulnerable function/config reachable from our code (`govulncheck`, call search with `rg`, vendor advisory "affected functions")? (3) Is the entry point exposed (internet, authenticated, internal, offline CLI)? (4) Is it in CISA KEV or has EPSS >= 0.1 / a public exploit?
- **Severity after triage:** Critical = reachable + exposed + known exploited or RCE; High = reachable and exposed; Medium = reachable but internal/authenticated, or unproven reachability in a runtime dep; Low = dev-only or not reachable. Record the advisory's CVSS but report the triaged severity.
- **Fix windows (default, adjust to team policy):** Critical 72h, High 7 days, Medium 30 days, Low next routine update. No dismissal without a recorded reason (`not reachable: <evidence>`, `dev-only`, `fix not available, mitigated by <control>`, review date).
- Prefer the smallest upgrade that fixes the advisory; use overrides/resolutions (`overrides`, `resolutions`, `constraints.txt`, `replace`) for transitive fixes and add a comment with the advisory ID and a removal condition.
- Lifecycle scripts run arbitrary code at install: use `npm ci --ignore-scripts` (or `ignore-scripts=true`) in CI where builds allow; pnpm 10+ blocks them by default - allowlist with `onlyBuiltDependencies`.
- A floating tag (`@v4`, `:latest`, `>=`) is a trust decision on every future publisher action; pin anything that runs with secrets or ships to users.
- `pull_request_target` or `workflow_run` that checks out PR head code with secrets or write tokens is Critical.
- Vendor-neutral tools first; use the platform's native features (Dependabot, secret scanning, attestations) where the repo already relies on them.

### Common failure modes
- Lockfile committed but CI runs `npm install`/`pip install -r` without hashes, so the lock is ignored.
- Scanner run only on the root manifest; monorepo subdirectories, Dockerfiles and base images unscanned.
- Alerts closed as "not used" with no evidence; the package is pulled in transitively by a runtime dependency.
- Actions pinned by SHA but composite actions or reusable workflows they call still use floating tags.
- Images signed by tag instead of digest, or signatures never verified at deploy time.
- Leaked secret "fixed" by deleting the commit, never rotated.
- Internal package name also published (or claimable) on the public registry and the resolver prefers the higher version.

## Safety gates
- Never print, paste or commit secret values; report type, location and a redacted fingerprint only. Use `--redact` flags.
- Rotating credentials, rewriting git history, force-pushing, publishing packages, deleting releases and changing registry or org settings are outward-facing: state the exact action and ask first.
- Do not install or execute a suspicious package to "test" it outside an isolated, network-restricted sandbox.
- Do not open public issues for unpatched vulnerabilities in third-party packages; use the project's private reporting route.

## Output
```markdown
## Supply-chain review: <repo/scope> (date)
**Ecosystems:** ... | **Tools run:** ... | **Not covered:** ...

| ID | Component / file | Issue | Advisory / evidence | Runtime? | Reachable / exposed | Triaged severity | Fix (version, pin, config) | Owner | Due |
|---|---|---|---|---|---|---|---|---|---|

**Pinning and CI gaps:** action/image refs not pinned, token permissions, risky triggers
**Release integrity:** SBOM [y/n], signing [y/n], provenance [SLSA level], consumer verification [y/n]
**Secrets:** findings (redacted), rotation status
**Proposed changes:** diffs or update-bot actions, awaiting approval
**Accepted risks:** ID, reason, owner, review date
```

---

**Phase:** Secure. **Next:** when this step is done, offer Ship: `/anthropica-devops:ship-ci-cd`, `/anthropica-devops:ship-github-actions`, `/anthropica-devops:ship-readiness`, `/anthropica-devops:ship-release`. Phase map: `/anthropica-devops:workflow`.
