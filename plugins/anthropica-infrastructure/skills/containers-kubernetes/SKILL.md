---
name: containers-kubernetes
description: Build and review container images and Kubernetes workloads to production standard (minimal non-root images, multi-stage builds, resource limits, probes, security contexts, Helm and Kustomize, safe rollouts). Use when editing Dockerfiles, compose files, Kubernetes manifests or Helm charts, or debugging pods.
---

# Containers and Kubernetes

## Dockerfiles
- Multi-stage builds: build tools in the builder stage only. Use a small, pinned base image (`distroless`, `-slim`, `alpine`), with a digest where policy requires it.
- Run as a non-root user (`USER 10001`). No secrets in layers, build args or `ENV`; use build secrets (`--secret`).
- Order layers for caching: dependency manifests and install first, source code last. Maintain a `.dockerignore`.
- One process per container. Use exec-form `ENTRYPOINT`/`CMD` so signals reach the app. Add `HEALTHCHECK` for non-Kubernetes runtimes.
- Scan images (`trivy image`, `docker scout`) and fix critical and high findings, or document why they don't apply.

## Kubernetes workloads
- `resources.requests` and `limits` on every container (memory limit always; CPU limit per team policy).
- `readinessProbe` gates traffic, `livenessProbe` restarts deadlocks, `startupProbe` covers slow starts. Do not point liveness at dependencies.
- `securityContext`: `runAsNonRoot: true`, `allowPrivilegeEscalation: false`, `readOnlyRootFilesystem: true`, drop all capabilities.
- Config in ConfigMaps, secrets from a secret manager (External Secrets, CSI driver, sealed secrets), never plain values committed in manifests.
- High availability: `replicas >= 2`, a `PodDisruptionBudget`, topology spread or anti-affinity, and graceful shutdown (`lifecycle.preStop`, `terminationGracePeriodSeconds`).
- Image tags are immutable versions or digests, never `latest`.
- NetworkPolicies default-deny where the cluster supports them. ServiceAccounts with minimal RBAC.

## Helm and Kustomize
- Keep environment differences in values or overlays; no copy-pasted manifests per environment.
- `helm lint`, `helm template | kubeconform`, `kustomize build | kubeconform` before applying.
- Prefer GitOps (Argo CD, Flux) over manual `kubectl apply` for shared clusters.

## Rollouts
- Rolling updates with `maxUnavailable: 0` for user-facing services. Canary or blue-green for risky changes.
- Watch rollout status and error rates. `kubectl rollout undo` is the rollback path; know it before deploying.

## Debugging a pod
```bash
kubectl get pods -n <ns> -o wide
kubectl describe pod <pod> -n <ns>        # events: scheduling, image pull, probe failures, OOMKilled
kubectl logs <pod> -n <ns> --previous     # crash-looping container
kubectl get events -n <ns> --sort-by=.lastTimestamp
```
Read-only commands are fine. `delete`, `scale`, `exec` and `edit` on shared clusters need explicit approval, and production changes go through the pipeline.

## Output
Images and manifests changed, security and resource settings, validation results, and the rollout and rollback plan.
