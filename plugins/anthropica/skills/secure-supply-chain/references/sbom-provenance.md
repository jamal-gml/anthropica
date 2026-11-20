# SBOM, signing and provenance

## SBOM
- Formats: CycloneDX (security/VEX oriented) or SPDX (licence/compliance oriented). Pick what consumers ask for; JSON output.
- Generate from the built artifact, not only the source tree, so OS packages and vendored files are included:
  - `syft <image-or-dir> -o cyclonedx-json=sbom.cdx.json` (or `-o spdx-json=sbom.spdx.json`)
  - `npm sbom --sbom-format cyclonedx --omit dev`
  - `cdxgen -o sbom.json .` for multi-language repos
  - `trivy image --format cyclonedx --output sbom.cdx.json <ref>`
- Attach the SBOM to the release and attest it (below). Scan the SBOM later for new CVEs: `grype sbom:sbom.cdx.json`.
- VEX (OpenVEX or CycloneDX VEX) to state "not affected" with justification for advisories you triaged as unreachable.

## Signing (Sigstore)
- Keyless signing in CI using the workflow's OIDC identity (needs `id-token: write`):
  - Image: `cosign sign --yes <registry>/<image>@sha256:<digest>` (sign the digest, never a tag).
  - Blob: `cosign sign-blob --yes --bundle artifact.sigstore.json <file>`.
- Verify by identity, not just "signed":
  ```bash
  cosign verify <image>@sha256:<digest> \
    --certificate-identity-regexp '^https://github.com/<owner>/<repo>/.github/workflows/release.yml@refs/tags/v' \
    --certificate-oidc-issuer https://token.actions.githubusercontent.com
  cosign verify-blob <file> --bundle artifact.sigstore.json --certificate-identity ... --certificate-oidc-issuer ...
  ```
- Key-based signing only when keyless is impossible; keep the key in a KMS/HSM (`cosign sign --key <kms-uri>`), never in the repo.
- Git: signed tags for releases (`git tag -s` or SSH signing; or gitsign keyless).

## Provenance and SLSA (v1.0 Build track)
| Level | Requirement | Typical way to get it |
|---|---|---|
| Build L1 | Provenance exists describing how the artifact was built | any generator, even unsigned |
| Build L2 | Hosted build platform generates and signs provenance | hosted CI + platform attestations / SLSA generator |
| Build L3 | Hardened build: isolated runs, signing keys inaccessible to user-defined steps, non-forgeable provenance | reusable/isolated workflow that builds and attests (e.g. SLSA GitHub generator, reusable workflow + attestation) |

- GitHub artifact attestations: `actions/attest-build-provenance` (subject path or image digest) with `id-token: write` and `attestations: write`; SBOM via `actions/attest-sbom`.
  - Verify: `gh attestation verify <file-or-oci://image> --repo <owner>/<repo>` (add `--signer-workflow` to pin the builder workflow).
- npm: publish with provenance from CI (`npm publish --provenance --access public`, `id-token: write`), or use trusted publishing (OIDC, no long-lived npm token). Consumers: `npm audit signatures`.
- PyPI: trusted publishing via OIDC (`pypa/gh-action-pypi-publish`, `id-token: write`), which also generates attestations; remove API tokens afterwards.
- Container registries: push signatures and attestations as OCI referrers next to the image digest.

## Consumer-side verification
- Deploy pipelines and admission controllers verify signature identity and provenance (builder, source repo, ref) before running an image: Kubernetes policy engines (Sigstore policy-controller, Kyverno `verifyImages`, OPA Gatekeeper with external data).
- Pin the expected builder identity; reject unsigned or wrong-identity artifacts by default, with a documented break-glass.
- Record digests deployed per environment for audit and rollback (`ship-release`).

## Checklist for a release pipeline
1. Build from a protected tag in a clean, ephemeral runner with pinned toolchain.
2. Produce artifact digests and an SBOM.
3. Sign and attest provenance + SBOM in the same workflow (no artifact hand-off to less-trusted jobs).
4. Publish with OIDC trusted publishing; no long-lived registry tokens.
5. Verify signatures/attestations in the deploy step and at admission.
