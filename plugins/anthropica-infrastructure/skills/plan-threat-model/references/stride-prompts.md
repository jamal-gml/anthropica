# STRIDE prompts per component type

Use these as question lists against the numbered DFD elements. Keep only threats that apply to this system; cite the element ID.

## Web front end / browser
- S: Session fixation, token theft via XSS, login CSRF, OAuth redirect manipulation (open `redirect_uri`, missing `state`/PKCE).
- T: DOM injection, unsafe `postMessage` handlers, third-party scripts without SRI, cache poisoning of static assets.
- I: Tokens in `localStorage`, sensitive data in URLs (logged by proxies, Referer), verbose errors, source maps in prod.
- E: Client-side-only authorisation (hidden buttons), feature flags evaluated client-side.

## API / service
- S: Missing auth on new routes, trusting `X-Forwarded-For`/`X-User-Id` headers set by clients, weak service-to-service auth.
- T: Mass assignment, parameter pollution, unsigned webhooks, replay (no nonce/timestamp), race conditions on balances or quotas.
- R: No audit log for admin, money or permission changes; logs a user can edit or delete.
- I: IDOR/BOLA, over-broad responses (whole DB row serialised), error stack traces, timing differences on user lookup.
- D: Unbounded pagination, regex DoS, large uploads, GraphQL depth/complexity, expensive search, no per-tenant rate limits.
- E: Function-level authZ gaps (user calls admin route), role from JWT claims not re-checked, SSRF to metadata endpoints.

## Identity and authentication
- Password reset token reuse/length/expiry, account enumeration via reset or signup, MFA bypass on alternate flows (API, mobile, SSO fallback), session not revoked on password change, JWT `alg` confusion or missing `aud`/`iss` checks, long-lived refresh tokens without rotation, SSO account linking by unverified email.

## Data stores (SQL, NoSQL, cache, object storage, search)
- T: Injection, writable by more principals than needed, no integrity on backups.
- I: Public buckets, presigned URLs with long expiry, unencrypted snapshots, cross-tenant reads in shared tables without row-level filtering, PII in search indexes and caches with no TTL.
- D: Hot partitions, connection exhaustion, storage growth from unbounded user content.
- Backups, replicas, analytics exports and test-data copies are separate data stores: model them.

## Queues, events, streams, webhooks
- S: Any producer can publish; consumer trusts message origin.
- T: Message modification, poison messages, out-of-order or duplicate delivery (idempotency).
- R: No correlation IDs across async hops.
- D: Retry storms, dead-letter queue never drained, unbounded fan-out.
- Outbound webhooks: SSRF via user-supplied URL, secret leakage in retries/logs, no signing so receivers cannot verify.

## File handling
- Upload type confusion (polyglots, SVG with script), path traversal on save, zip slip/zip bombs, image parser CVEs, serving user files from the app origin, malware distribution, metadata (EXIF GPS) leakage.

## CI/CD and build
- S/E: Workflows on `pull_request_target` checking out fork code, self-hosted runners shared across trust levels, OIDC trust policy not restricted to repo/branch/environment.
- T: Unpinned actions/images/scripts, cache poisoning, artifact swap between build and deploy, unsigned releases.
- I: Secrets in logs, secrets available to PR builds, long-lived cloud keys.
- R: Deploys without traceable approval. See `secure-supply-chain` and `secure-repo-governance`.

## Cloud and infrastructure
- Over-broad IAM roles (wildcards), public network exposure (0.0.0.0/0, public IPs on DBs), instance metadata reachable via SSRF (require IMDSv2 or equivalent), cross-account trust without external ID/conditions, unencrypted volumes, disabled audit logging, management plane reachable from workloads, KMS key policies granting decrypt widely.

## Multi-tenancy
- Tenant ID from client input rather than session, missing tenant filter in one query path (reports, exports, search, background jobs), shared caches keyed without tenant, per-tenant encryption keys absent where contractually required, noisy-neighbour DoS.

## Mobile and desktop clients
- Hard-coded API keys, certificate pinning absent for high-risk apps, insecure local storage, deep-link hijacking, exported components, debug builds in store, jailbreak/root assumptions.

## LLM, agents and AI features
- S/T: Direct and indirect prompt injection (from documents, web pages, emails, tool results) steering tool calls.
- E: Agent tools with more privilege than the requesting user (confused deputy); no human approval for destructive or outward actions.
- I: System prompt or other users' data leaked via retrieval (RAG index without per-user ACL), training-data or log retention of sensitive input, data exfiltration via rendered links/images.
- D: Token/cost exhaustion, recursive tool loops.
- T: Model output rendered as HTML/SQL/shell without encoding; poisoned embeddings or fine-tuning data.
- Treat all model output as untrusted input crossing a boundary.

## Third-party integrations
- Over-scoped OAuth grants, vendor holding long-lived credentials, vendor outage as DoS, vendor breach exposing shared data, callback endpoints without signature verification, data residency of the vendor.

## Operations and people
- Support/admin impersonation tools without audit or approval, break-glass accounts without alerting, observability pipelines capturing PII or tokens, offboarding gaps, shared accounts.
