# Per-class checklist: sinks, checks, safe patterns

Search tips use ripgrep; adapt to the language. A hit is a lead; trace it before reporting.

## Injection
- Sinks: string-built SQL (`f"SELECT ... {x}"`, `"..." + x`, `.raw(`, `.extra(`, `query(` with template literals), NoSQL operators from JSON bodies (`{"$ne": null}`, `$where`), `exec`/`system`/`subprocess(..., shell=True)`/`child_process.exec`, `eval`, server-side template rendering of user strings (`render_template_string`, `Template(user)`), LDAP filters, XPath, header values (CRLF).
- Safe: parameterised queries/bind variables, ORM query builders without raw fragments, `execFile`/argument arrays with `shell=False`, allowlisted identifiers for dynamic column/table names, schema validation that rejects objects where strings are expected.
- Search: `rg -n "(execute|query|raw)\(.*(\+|%|\$\{|f\")"`, `rg -n "shell=True|child_process|Runtime.getRuntime"`.

## Broken access control (IDOR/BOLA, function-level)
- Every handler that takes an object ID: is the lookup scoped by owner/tenant (`WHERE id = ? AND org_id = ?`) or followed by a policy check before use?
- List, search, export, bulk and GraphQL resolvers: same filter as the detail route? Nested resolvers and `include`/`expand` params leaking related objects?
- Admin and internal routes: role check on the server, not only hidden in UI; internal endpoints not exposed through the public router or gateway.
- Background jobs and webhooks acting on IDs from messages: is the actor's permission re-checked?
- Search: route definitions without auth decorators/middleware; `findById(req.params`, `get_object_or_404(Model, pk=` without owner filter.

## Mass assignment
- Request body bound directly to models (`Model(**body)`, `Object.assign(entity, req.body)`, `update(req.body)`, Rails `permit!`, Spring binding to entities, ASP.NET binding to EF entities).
- Safe: explicit DTO/schema with allowlisted fields; server sets `ownerId`, `role`, `status`, `price`, `tenantId`.

## Authentication and session
- JWT: library verifies signature with a fixed algorithm list (reject `none`, HS/RS confusion), checks `exp`, `nbf`, `aud`, `iss`; secret strength; no sensitive data in payload; revocation strategy for logout/password change.
- Sessions: rotate ID on login and privilege change, `HttpOnly`, `Secure`, `SameSite=Lax/Strict`, idle and absolute timeouts, server-side invalidation on logout.
- Passwords: adaptive hash (argon2id, scrypt, bcrypt cost >= 10), no truncation surprises (bcrypt 72-byte limit), breached-password check where policy needs it.
- Reset/verification tokens: >= 128 bits from CSPRNG, single use, short expiry, stored hashed, not leaked via Referer or logs; same response whether the account exists or not.
- Brute force: rate limit / lockout on login, MFA verify, reset and OTP endpoints. MFA enforced on all login paths (API tokens, SSO fallback, mobile).
- OAuth/OIDC: `state` and PKCE, exact `redirect_uri` match, verify ID token, do not link accounts by unverified email.

## SSRF
- Sinks: HTTP clients, image/PDF renderers, webhook senders, URL previewers, XML external entities, import-from-URL features, headless browsers.
- Checks: scheme allowlist (https), resolve DNS and block private, loopback, link-local (169.254.0.0/16 incl. cloud metadata), IPv6 equivalents and decimal/octal encodings; re-validate after redirects or disable redirects; pin the resolved IP for the connection (DNS rebinding); egress proxy or network policy as defence in depth.

## Insecure deserialization and parsers
- Sinks: `pickle.loads`, `yaml.load` without SafeLoader, Java `ObjectInputStream`, .NET `BinaryFormatter`/`TypeNameHandling` != None, PHP `unserialize`, Ruby `Marshal.load`, `node-serialize`.
- XML: disable DTDs/external entities (XXE). Archives: zip slip (normalise each entry path and confirm it stays under the target dir), zip bombs (limit total size and entry count).
- Webhooks: verify HMAC over the raw body with constant-time compare and a timestamp window.

## Path traversal and file uploads
- User input in `open`, `readFile`, `send_file`, `Path.join`/`os.path.join` (absolute second arg resets the base), static file servers, template names.
- Safe: map IDs to server-chosen paths; else `realpath`/`resolve` then assert prefix of the allowed base.
- Uploads: size limits, allowlisted extensions and sniffed content type, random server-side names, store outside web root or in object storage, serve with `Content-Disposition: attachment` and `X-Content-Type-Options: nosniff` from a separate domain; SVG/HTML treated as active content; image processing libraries patched.

## XSS and CSP
- Sinks: `innerHTML`, `outerHTML`, `document.write`, `dangerouslySetInnerHTML`, `v-html`, `[innerHTML]`, `bypassSecurityTrust*`, `|safe`, `mark_safe`, `html_safe`, `raw`, `@Html.Raw`, `href`/`src` with `javascript:` URLs, jQuery `.html()`, markdown renderers without sanitising.
- Context matters: HTML body, attribute, URL, JS string and CSS need different encoding. Sanitise rich HTML with a maintained sanitiser (DOMPurify or equivalent).
- CSP: prefer nonce- or hash-based `script-src` with `strict-dynamic`, `object-src 'none'`, `base-uri 'none'`, `frame-ancestors`; flag `unsafe-inline`/`unsafe-eval` and wildcard sources.

## CSRF and CORS
- Cookie-authenticated state-changing requests need framework CSRF tokens or SameSite + Origin/Referer verification. GET must not change state.
- CORS: never reflect arbitrary `Origin` with `Access-Control-Allow-Credentials: true`; `null` origin is not safe; regex origin checks must be anchored.

## Cryptography
- Flag: MD5/SHA1 for security, DES/3DES/RC4, AES-ECB, CBC without MAC, static or reused IV/nonce (fatal for GCM), `Math.random`/`random` for tokens, keys in code, custom crypto, disabled TLS verification (`verify=False`, `rejectUnauthorized: false`, `InsecureSkipVerify`), non-constant-time comparison of secrets.
- Safe: AEAD (AES-GCM, ChaCha20-Poly1305) via a high-level library, keys from a KMS/secret manager, `secrets`/`crypto.randomBytes`/`crypto/rand`, `hmac.compare_digest`/`timingSafeEqual`.

## Secrets and sensitive logging
- Hard-coded credentials, private keys, connection strings, tokens in tests or fixtures that are real, `.env` committed.
- Logs, traces, error reports and analytics: passwords, tokens, `Authorization` headers, cookies, full request bodies, card data, national IDs, health data. Check logger redaction config and exception handlers that dump request context.

## Race conditions and business logic
- Check-then-act on balance, stock, quotas, coupons, invites, one-time tokens, file existence (TOCTOU). Safe: DB transaction with row lock (`SELECT ... FOR UPDATE`), atomic conditional update (`UPDATE ... WHERE balance >= ?`), unique constraints, idempotency keys.
- Negative quantities, integer overflow, currency rounding, workflow steps skippable by calling a later endpoint directly, client-supplied prices or totals.

## Denial of service in code
- Unbounded loops over user-controlled sizes, missing pagination limits, ReDoS (nested quantifiers on user input), decompression without limits, expensive hashing on unauthenticated endpoints without rate limits, GraphQL depth/complexity.

## LLM and agent code
- Model output passed to `eval`, shell, SQL, HTML or tool calls without validation; tools executed with service privileges rather than the user's; retrieved documents not filtered by the user's permissions; prompt injection from fetched content able to trigger outward actions without confirmation.
