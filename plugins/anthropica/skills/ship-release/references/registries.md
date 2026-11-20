# Publishing to registries

All registries: publish from CI on the tag, with OIDC trusted publishing or a short-lived scoped token, after build and tests pass on that exact commit. Dry-run first where supported. Published versions are effectively permanent.

| Ecosystem | Dry run / check | Publish | Auth best practice | Notes |
|---|---|---|---|---|
| npm | `npm pack --dry-run`; `npm publish --dry-run` | `npm publish --provenance --access public` | trusted publishing (OIDC) from CI | `prepublishOnly` must build; check `files`/`.npmignore` so no secrets or tests ship; dist-tags (`latest`, `next`) for pre-releases: `--tag next` |
| PyPI | `python -m build`; `twine check dist/*` | `pypa/gh-action-pypi-publish` or `twine upload` | trusted publishing | upload sdist and wheels together; TestPyPI for rehearsal; versions can be yanked, not reused |
| Maven Central | `mvn -P release verify` | central publishing plugin / `gradle publish` | portal tokens in CI secrets | requires signed artifacts (GPG), sources and javadoc jars; releases are immutable |
| NuGet | `dotnet pack` | `dotnet nuget push *.nupkg --source <feed>` | trusted publishing or scoped API key with glob | unlist instead of delete |
| crates.io | `cargo publish --dry-run` | `cargo publish` | trusted publishing or scoped token | `cargo yank --version X.Y.Z` prevents new lockfiles from selecting it |
| Go modules | `go mod tidy` clean; `GOPROXY=direct go list -m` | push tag `vX.Y.Z` (submodules: `path/vX.Y.Z`) | n/a | v2+ requires `/v2` module path; versions are cached forever by the proxy; use `retract` in `go.mod` |
| RubyGems | `gem build` | `gem push` | trusted publishing, MFA | `gem yank` |
| Container images | `docker buildx build` without push; scan | `docker buildx build --push` with tags `X.Y.Z`, `X.Y`, `sha-<short>` | OIDC to cloud registry, `GITHUB_TOKEN` for GHCR | deploy by digest; sign and attest (`cosign`, build provenance); multi-arch via `--platform` |
| Helm charts | `helm lint`; `helm package` | `helm push <chart>.tgz oci://<registry>` | OIDC registry login | bump chart `version` for every change, `appVersion` for app changes |
| GitHub/GitLab Releases | `gh release view` | `gh release create ... --verify-tag files...` | CI token with `contents: write` | attach `checksums.txt` (`sha256sum * > checksums.txt`) |

## Pre-publish checklist
- [ ] Version in manifest equals tag; not already published (`npm view <pkg>@<v>`, `pip index versions <pkg>`, registry UI/API).
- [ ] Package contents inspected (`npm pack --dry-run`, `tar tzf`, `unzip -l`): no `.env`, keys, test fixtures with data, source maps if not intended.
- [ ] Licence and metadata fields correct (repository URL, homepage, description).
- [ ] Provenance/signature generated (`secure-supply-chain`).
- [ ] For monorepos: publish order respects inter-package dependencies; all-or-nothing is not guaranteed, so record what got published.

## After publishing
- Verify install from a clean environment: `npm i <pkg>@<v>`, `pip install <pkg>==<v>`, `docker pull <image>@<digest>`.
- Verify provenance: `npm audit signatures`, `gh attestation verify`, `cosign verify`.
- Announce in release notes; update docs that pin versions.
