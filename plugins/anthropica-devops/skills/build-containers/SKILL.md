---
name: build-containers
description: Build and review container images and Kubernetes workloads (pinned non-root images, scanning, resources, probes, securityContext, PDB/HPA, graceful shutdown, Helm/Kustomize, rollout/rollback) and debug pods. Use when editing Dockerfiles, compose, manifests or Helm charts, or on CrashLoopBackOff, ImagePullBackOff, OOMKilled or Pending. Not for image signing/SBOM (use `secure-supply-chain`).
---

# Containers and Kubernetes

Ship small, reproducible, non-root images and workloads that schedule, start, scale and stop cleanly.

## Procedure
1. **Discover.** Runtime/orchestrator (Kubernetes version, managed flavour, compose, serverless containers), how manifests are rendered (raw, Kustomize, Helm, GitOps app), registry, existing base images and admission policies (Pod Security Standards level, Kyverno/Gatekeeper).
2. **Image.** Apply the Dockerfile rules; read `references/dockerfile.md` for patterns per language. Build with BuildKit: `docker buildx build --platform linux/amd64,linux/arm64 -t <reg>/<app>:<ver> .` as needed.
3. **Scan and inspect.** `trivy image --severity HIGH,CRITICAL --ignore-unfixed <img>`; `docker history <img>` and `dive <img>` for size/leaked files; `docker inspect -f '{{.Config.User}}' <img>` must not be empty or `root`.
4. **Workload.** Apply the Kubernetes rules; read `references/kubernetes.md` for a full annotated Deployment, PDB, HPA and probe examples.
5. **Validate offline.**
   ```bash
   helm lint <chart> && helm template <rel> <chart> -f values-<env>.yaml | kubeconform -strict -summary -
   kustomize build overlays/<env> | kubeconform -strict -summary -
   kubectl apply --dry-run=server -f <rendered>.yaml      # admission/policy check, read-only
   kubectl diff -f <rendered>.yaml                         # what will change in the cluster
   ```
   Also `trivy config .` / `kube-score` / `polaris` if available.
6. **Roll out** (after approval) through GitOps or the pipeline, then watch: `kubectl rollout status deploy/<app> -n <ns> --timeout=5m`, error rate and latency (`operate-observability`).
7. **Roll back** if SLIs degrade: revert the Git change (GitOps) or `kubectl rollout undo deploy/<app> -n <ns> [--to-revision=N]`; `helm rollback <rel> <rev>`. Know the previous revision (`kubectl rollout history`, `helm history`) before deploying.

## Rules

**Images**
- Multi-stage: compilers and dev deps only in the builder. Runtime base: distroless/`-slim`/`chainguard`/`alpine` (alpine = musl; watch DNS and glibc-linked binaries). Pin base by digest (`FROM image:tag@sha256:...`) and let a bot bump it.
- Non-root numeric user (`USER 10001:10001`) so `runAsNonRoot` can verify it. Files owned by root and read-only; writable dirs explicit.
- Cache order: lockfile copy + install, then source. `.dockerignore` excludes `.git`, `node_modules`, build output, `.env*`, keys, tests.
- No secrets in layers, `ARG` or `ENV` (visible in `docker history`); use `RUN --mount=type=secret`. Deleting a file in a later layer does not remove it.
- Exec-form `ENTRYPOINT ["/app"]`; the app is PID 1 and must handle SIGTERM, or use `tini`/`--init`. Shell form swallows signals.
- `HEALTHCHECK` only matters for Docker/compose/ECS; Kubernetes ignores it.
- Tags immutable (semver or git SHA); deploy by digest where possible. Never `latest`.
- Scan gate: no fixable CRITICAL; HIGH fixed or documented with expiry in `.trivyignore`.

**Probes** (wrong probes cause outages)
| Probe | Question | Must not | Failure effect |
|---|---|---|---|
| startupProbe | has it finished booting? | be too short for cold start | kills container; liveness/readiness wait for it |
| readinessProbe | can it take traffic now? | be identical to liveness by default | removed from Service endpoints, no restart |
| livenessProbe | is the process wedged beyond repair? | check DB/downstreams | container restarted |
- Liveness on dependencies turns a DB blip into a fleet-wide restart storm. Readiness may check critical local deps, cautiously.
- Startup budget = `failureThreshold x periodSeconds` >= worst-case boot (e.g. 30 x 5s = 150s). Liveness: `periodSeconds: 10, failureThreshold: 3` typical; `timeoutSeconds` >= p99 of the endpoint.

**Resources**
- Requests on every container (scheduler uses them). Memory limit = request for predictable QoS; set heap/runtime to ~75% of limit (`-XX:MaxRAMPercentage=75`, `GOMEMLIMIT`, `--max-old-space-size`).
- CPU limits cause CFS throttling even when nodes are idle; for latency-sensitive services prefer CPU request without limit unless policy requires one. Check `container_cpu_cfs_throttled_periods_total`.
- Right-size from real usage (p95 over a week), not guesses.

**Security** (Pod Security Standards "restricted")
- `runAsNonRoot: true`, `allowPrivilegeEscalation: false`, `readOnlyRootFilesystem: true` (+ `emptyDir` for `/tmp`), `capabilities.drop: [ALL]`, `seccompProfile.type: RuntimeDefault`.
- `automountServiceAccountToken: false` unless the pod calls the API; dedicated ServiceAccount, minimal Role. NetworkPolicy default-deny + explicit allows.
- Secrets from a secret manager (External Secrets, CSI driver) or sealed/SOPS-encrypted; never plain `Secret` YAML in Git. Env-var secrets leak into crash dumps; prefer mounted files.

**Availability and shutdown**
- `replicas >= 2` (3 across zones for critical), PDB `maxUnavailable: 1` (never a PDB that allows zero disruptions with 1 replica - it blocks node drains), `topologySpreadConstraints` on `topology.kubernetes.io/zone`.
- Rolling update `maxUnavailable: 0, maxSurge: 25%` for user-facing services.
- Shutdown: endpoints removal races SIGTERM. Add `preStop: sleep 5-15` (or `sleep` action on K8s >= 1.30), app drains in-flight work on SIGTERM, `terminationGracePeriodSeconds` > preStop + max drain time.
- HPA needs requests set; target CPU 60-75%; set `behavior.scaleDown.stabilizationWindowSeconds` (300 default) to avoid flapping. Do not set `replicas` in the manifest when HPA owns it (GitOps will fight it).

**Helm/Kustomize/GitOps**
- Environment differences in values files or overlays only. Pin chart versions. Render and diff before every apply. Prefer Argo CD/Flux on shared clusters; never `kubectl edit` live objects that GitOps manages.

## Troubleshooting
| Status | First command | Usual causes |
|---|---|---|
| Pending | `kubectl describe pod` (Events) | insufficient CPU/memory requests, taints/affinity, PVC unbound, quota |
| ImagePullBackOff / ErrImagePull | describe -> Events | wrong tag/digest, private registry without `imagePullSecrets`, arch mismatch, rate limit |
| CrashLoopBackOff | `kubectl logs <pod> --previous` | app error, missing config/secret, bad command, liveness too aggressive |
| OOMKilled (exit 137) | `kubectl get pod -o jsonpath='{.status.containerStatuses[*].lastState}'` | limit too low, heap not sized to limit, leak |
| Exit 143 / slow stop | events, app logs | SIGTERM not handled, grace period too short |
| Running, not Ready | describe -> readiness failures | wrong port/path, dependency check in readiness, slow warm-up |
| CreateContainerConfigError | describe | missing ConfigMap/Secret key |
| Evicted | `kubectl get events` | node memory/disk pressure, no requests set |
More: `kubectl get events -n <ns> --sort-by=.lastTimestamp`, `kubectl top pod`, `kubectl debug -it <pod> --image=busybox --target=<container>` (ephemeral container, ask first in production).

## Safety gates
- Read-only (`get`, `describe`, `logs`, `top`, `diff`, `--dry-run=server`) is fine. Ask before `apply`, `delete`, `scale` (especially to 0 or down), `rollout restart/undo`, `exec`, `debug`, `edit`, `drain`/`cordon`, `helm upgrade/rollback/uninstall`, and before pushing images to shared registries.
- Confirm the context first: `kubectl config current-context` and namespace. Production changes go through the pipeline/GitOps with the user's explicit "yes".
- Never delete PVCs, namespaces or CRDs without explicit confirmation of data loss.

## Output
```
Image:       <ref@digest>, size <MB>, base <image@digest>, user <uid>, scan C/H: <n>/<n>
Manifests:   <files/charts changed>
Settings:    requests/limits, probes (startup/readiness/liveness), securityContext, PDB, HPA, spread
Validation:  helm lint / kubeconform / dry-run=server / kubectl diff - results
Rollout:     strategy, watch commands, success criteria (SLIs)
Rollback:    <git revert | rollout undo --to-revision=N | helm rollback rel rev>
Open risks:  <...>
```

---

**Phase:** Build. **Next:** when this step is done, offer Review: `/anthropica-devops:review-pull-request`. Phase map: `/anthropica-devops:workflow`.
