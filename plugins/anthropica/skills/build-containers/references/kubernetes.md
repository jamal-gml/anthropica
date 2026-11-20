# Kubernetes workload reference

## Annotated Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  labels: { app.kubernetes.io/name: api, app.kubernetes.io/version: "1.8.2" }
spec:
  # replicas omitted when an HPA owns scaling
  revisionHistoryLimit: 5
  strategy:
    type: RollingUpdate
    rollingUpdate: { maxUnavailable: 0, maxSurge: 25% }
  selector:
    matchLabels: { app.kubernetes.io/name: api }      # immutable after creation
  template:
    metadata:
      labels: { app.kubernetes.io/name: api }
    spec:
      serviceAccountName: api
      automountServiceAccountToken: false
      terminationGracePeriodSeconds: 45               # > preStop + drain
      securityContext:
        runAsNonRoot: true
        runAsUser: 10001
        fsGroup: 10001
        seccompProfile: { type: RuntimeDefault }
      topologySpreadConstraints:
        - maxSkew: 1
          topologyKey: topology.kubernetes.io/zone
          whenUnsatisfiable: ScheduleAnyway            # DoNotSchedule if zone HA is mandatory
          labelSelector: { matchLabels: { app.kubernetes.io/name: api } }
      containers:
        - name: api
          image: registry.example/api@sha256:...
          ports: [{ name: http, containerPort: 8080 }]
          envFrom: [{ configMapRef: { name: api-config } }]
          volumeMounts:
            - { name: tmp, mountPath: /tmp }
            - { name: secrets, mountPath: /var/run/secrets/app, readOnly: true }
          resources:
            requests: { cpu: 250m, memory: 512Mi }
            limits:   { memory: 512Mi }                 # CPU limit per policy
          startupProbe:
            httpGet: { path: /healthz, port: http }
            periodSeconds: 5
            failureThreshold: 30                        # 150s boot budget
          readinessProbe:
            httpGet: { path: /ready, port: http }
            periodSeconds: 5
            failureThreshold: 3
          livenessProbe:
            httpGet: { path: /healthz, port: http }     # process-local only
            periodSeconds: 10
            timeoutSeconds: 2
            failureThreshold: 3
          lifecycle:
            preStop: { sleep: { seconds: 10 } }         # K8s >= 1.30; else exec ["sleep","10"]
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities: { drop: ["ALL"] }
      volumes:
        - { name: tmp, emptyDir: { sizeLimit: 256Mi } }
        - name: secrets
          secret: { secretName: api-secrets }           # synced by External Secrets / CSI
```

## PodDisruptionBudget
```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata: { name: api }
spec:
  maxUnavailable: 1
  selector: { matchLabels: { app.kubernetes.io/name: api } }
```
- `minAvailable: 100%` or `maxUnavailable: 0` blocks node upgrades forever. With 1 replica, any PDB that protects it blocks drains; run 2+.

## HorizontalPodAutoscaler
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata: { name: api }
spec:
  scaleTargetRef: { apiVersion: apps/v1, kind: Deployment, name: api }
  minReplicas: 3
  maxReplicas: 20
  metrics:
    - type: Resource
      resource: { name: cpu, target: { type: Utilization, averageUtilization: 70 } }
  behavior:
    scaleUp:   { stabilizationWindowSeconds: 0 }
    scaleDown: { stabilizationWindowSeconds: 300, policies: [{ type: Percent, value: 50, periodSeconds: 60 }] }
```
- Utilization is relative to requests. Queue workers: scale on queue depth/lag (KEDA or external metrics), not CPU.
- `maxReplicas` must fit cluster/node-pool capacity and downstream limits (DB connections = replicas x pool size).

## Graceful shutdown sequence
1. Pod marked Terminating; endpoint removal starts propagating (async, seconds).
2. `preStop` runs (sleep keeps serving while load balancers catch up).
3. SIGTERM to PID 1: stop accepting, finish in-flight, close pools, exit 0.
4. After `terminationGracePeriodSeconds` total: SIGKILL (exit 137).
Workers: stop polling, finish/requeue current message within the grace period; make handlers idempotent.

## Jobs and CronJobs
- `backoffLimit`, `activeDeadlineSeconds`, `ttlSecondsAfterFinished`; CronJob `concurrencyPolicy: Forbid` for non-reentrant jobs, `startingDeadlineSeconds` to avoid missed-run storms.

## Config and secrets
- ConfigMap changes do not restart pods. Use a checksum annotation (`checksum/config: {{ include ... | sha256sum }}` in Helm) or Kustomize `configMapGenerator` hash suffix to roll on change.
- Secrets in env vars do not update; mounted secrets update eventually (not with `subPath`).

## NetworkPolicy default deny
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: { name: default-deny }
spec:
  podSelector: {}
  policyTypes: [Ingress, Egress]
```
Then allow ingress from the ingress controller namespace and egress to DNS (UDP/TCP 53 to kube-dns) plus required backends. Forgetting DNS egress breaks everything.

## Helm practice
- `helm upgrade --install <rel> <chart> -n <ns> --atomic --wait --timeout 10m -f values-<env>.yaml` (`--atomic` rolls back on failure).
- `helm diff upgrade` (plugin) before upgrade. `helm get values <rel>` to see what is live.
- Charts: `values.schema.json` for validation, no logic-heavy templates, `required` for mandatory values.

## Kustomize practice
- `base/` + `overlays/<env>/`; patches small and targeted; `images:` transformer to set tag/digest per env.

## Rollout commands
```bash
kubectl rollout status deploy/api -n <ns> --timeout=5m
kubectl rollout history deploy/api -n <ns>
kubectl rollout undo deploy/api -n <ns> --to-revision=<N>      # approval required
kubectl rollout pause|resume deploy/api -n <ns>
kubectl argo rollouts get rollout api -n <ns> --watch           # if Argo Rollouts canary
```
Progressive delivery (Argo Rollouts, Flagger): canary steps 5% -> 25% -> 50% -> 100% with analysis on error rate and p99 latency; abort automatically on failure.

## Deeper debugging
```bash
kubectl get pod <pod> -n <ns> -o yaml | less          # full status, lastState, reasons
kubectl get endpointslices -n <ns> -l kubernetes.io/service-name=<svc>   # Ready endpoints?
kubectl auth can-i <verb> <resource> --as=system:serviceaccount:<ns>:<sa>
kubectl describe node <node>                          # pressure, allocatable vs requested
kubectl run tmp --rm -it --image=busybox --restart=Never -- sh   # network tests (ask first)
```
- Service has no endpoints: selector labels do not match pod labels, or no pod is Ready.
- DNS failures: NetworkPolicy blocking 53, CoreDNS saturation, `ndots:5` causing slow external lookups (use FQDN with trailing dot).
