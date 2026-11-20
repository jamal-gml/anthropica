# Per-ecosystem commands

Check tool versions before relying on flags; prefer the repo's existing tooling.

| Ecosystem | Lock / frozen install | Audit | Confusion / registry control |
|---|---|---|---|
| npm | `package-lock.json`; `npm ci` (fails on mismatch) | `npm audit --omit=dev`; `npm audit signatures` (registry signatures and provenance) | `.npmrc`: `@scope:registry=https://<private>`; scoped internal names; `ignore-scripts=true` |
| pnpm | `pnpm-lock.yaml`; `pnpm install --frozen-lockfile` | `pnpm audit --prod` | `.npmrc` scope registry; `minimumReleaseAge` (pnpm 10.16+); `onlyBuiltDependencies` |
| Yarn (berry) | `yarn.lock`; `yarn install --immutable` | `yarn npm audit --environment production` | `npmScopes` in `.yarnrc.yml` |
| Python (pip) | `requirements.txt` with hashes from `pip-compile --generate-hashes`; `pip install --require-hashes --no-deps -r` | `pip-audit -r requirements.txt` (or in the venv) | one `--index-url` pointing at a proxy that mirrors public and serves private; never `--extra-index-url` for internal names |
| Python (uv) | `uv.lock`; `uv sync --locked` (fail if stale) or `--frozen` | `uv pip compile` + `pip-audit`, or `osv-scanner` on `uv.lock` | `[[tool.uv.index]]` with `explicit = true` for private packages |
| Python (Poetry) | `poetry.lock`; `poetry check --lock`, `poetry install --no-root` | `pip-audit` on exported requirements | source priority `explicit` for private sources |
| Go | `go.sum`; `go mod verify`, `GOFLAGS=-mod=readonly` | `govulncheck ./...` (symbol-level reachability) | `GOPRIVATE`, `GONOSUMDB`, `GOPROXY` |
| Rust | `Cargo.lock` (commit it for binaries); `cargo build --locked` | `cargo audit`; `cargo deny check advisories bans licenses sources` | `[registries]` + `cargo deny` sources allowlist |
| Maven | no native lock; pin versions, `maven-enforcer-plugin` (`requireUpperBoundDeps`, ban SNAPSHOT) | `mvn org.owasp:dependency-check-maven:check` or `osv-scanner` | repository mirrors in `settings.xml`; checksum policy `fail` |
| Gradle | `gradle dependencies --write-locks`; dependency verification `--write-verification-metadata sha256` | dependency-check plugin or `osv-scanner` | `repositories` content filtering (`includeGroup`) per repo |
| .NET | `packages.lock.json` with `RestorePackagesWithLockFile`; `dotnet restore --locked-mode` | `dotnet list package --vulnerable --include-transitive` | `nuget.config` `packageSourceMapping` |
| Ruby | `Gemfile.lock`; `BUNDLE_FROZEN=true bundle install` | `bundle-audit check --update` | `source "https://<private>" do ... end` blocks per gem |
| PHP | `composer.lock`; `composer install --no-dev` in prod | `composer audit` | `repositories` with `canonical: true` |
| Containers | base images by `FROM image@sha256:...`; `apt`/`apk` versions pinned where reproducibility matters | `trivy image <ref>`, `grype <ref>`; rebuild regularly for OS patches | private registry mirror; deny unknown registries at admission |
| Terraform | `.terraform.lock.hcl` committed; `terraform init -lockfile=readonly` | `trivy config .`, `checkov -d .` | modules pinned to tag + commit or registry version |
| GitHub Actions | `uses: owner/repo@<40-char SHA> # vX.Y.Z` | `zizmor .github/workflows`, `actionlint`; Dependabot `github-actions` | allowed-actions policy (see `secure-repo-governance`) |

## Cross-ecosystem scanners
- `osv-scanner scan source -r .` - reads most lockfiles, OSV database, supports `--call-analysis` for Go/Rust.
- `trivy fs --scanners vuln,secret,misconfig .` - deps, secrets, IaC in one pass.
- `syft` + `grype` - SBOM then vulnerability match on the SBOM.

## Reachability triage helpers
- Go: `govulncheck` output lists call stacks; "imported but not called" is Low.
- JS/TS: search for the advisory's vulnerable export (`rg -n "from '<pkg>'|require\('<pkg>'\)"`), then the specific function; check whether only a dev tool imports it (`npm ls <pkg>`, `pnpm why <pkg>`).
- Python: `pipdeptree -r -p <pkg>` to find who pulls it; search for the vulnerable API.
- Java: `mvn dependency:tree -Dincludes=<group>:<artifact>`; `gradle dependencyInsight --dependency <name>`.
- .NET: `dotnet nuget why <project> <package>` (SDK 8.0.4xx+).

## Typosquat and package-health checks
- Compare against the canonical name from the project's own docs/repo, not search results.
- Red flags: created recently, single maintainer with no history, install/postinstall scripts, obfuscated code, network calls at install, repository link missing or pointing elsewhere, version jump with new maintainer, download count mismatch with popularity claims.
- Inspect before installing: `npm view <pkg> maintainers repository scripts time --json`, `pip download <pkg> --no-deps` then read the sdist/wheel in a sandbox.
