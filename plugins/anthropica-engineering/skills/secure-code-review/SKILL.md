---
name: secure-code-review
description: Security-review code or a diff by tracing untrusted input source to sink (injection, IDOR/authZ, authN/session, SSRF, deserialization, path traversal, XSS, CSRF, crypto, secrets, logging, races, mass assignment), proving reachability, rating severity and giving fixes. Use when a change touches auth, user input, queries, file paths, server-side fetches, templates, crypto or secrets, or for a PR security review. Not for design-stage threats (use `plan-threat-model`).
---

# Secure Code Review

Find realistic, reachable vulnerabilities in code and prove each one with a source-to-sink trace. Fewer confirmed findings beat many speculative ones.

## Procedure
1. **Scope the attack surface.** For a diff: `git diff --stat <base>...HEAD`, then list new or changed entry points (routes, handlers, RPC methods, message consumers, CLI args, file parsers, webhooks, scheduled jobs reading external data). For a codebase: enumerate entry points from the router/framework config first. Read the threat model if one exists.
2. **Identify trust context per entry point:** authenticated? which roles? tenant scoped? internet-facing? This decides severity later.
3. **Run cheap automation first, if available and offline-safe:** `semgrep scan --config p/default --config p/owasp-top-ten --error` (or `p/<language>`), `bandit -r src` (Python), `gosec ./...` (Go), `brakeman` (Rails), `eslint-plugin-security`, plus `gitleaks detect --no-banner --redact` for secrets. Treat results as leads, not findings.
4. **Trace source to sink.** For each untrusted source (params, body, headers, cookies, path, files, env from users, DB values written by users, third-party API responses, LLM output), follow data through transforms to dangerous sinks. Record each hop as `file:line`. Note validation, encoding or authorisation encountered and whether it is correct for that sink.
5. **Check every class in the table below** against the changed code. Read `references/vuln-classes.md` for the detailed per-class checklist, sinks and safe patterns.
6. **Check framework defaults and gotchas.** Read `references/frameworks.md` when the code uses Express/Node, Next.js/React, Django, Flask/FastAPI, Rails, Spring, ASP.NET, Go net/http, Laravel/PHP or an ORM.
7. **Prove reachability** for each candidate: an attacker-controlled path from an entry point the attacker can call, with their privilege level, to the sink, with no effective control in between. If you cannot show it, report it as a hardening note, not a vulnerability.
8. **Rate, map and fix** each finding (rules below). Give the minimal fix plus a regression test (unit or integration test that sends the malicious input and asserts rejection).
9. **Verify** fixes by re-reading the patched path and running the tests and scanners again.

### Class checklist (OWASP Top 10 2021 mapping)
| Class | Core question | OWASP / CWE |
|---|---|---|
| Broken access control, IDOR/BOLA | Is ownership/tenant checked server-side on every object read, write, list, export? | A01 / CWE-639, 862, 863 |
| Mass assignment | Can the client set fields like `role`, `isAdmin`, `ownerId`, `price`? | A01 / CWE-915 |
| CSRF | Do cookie-authenticated state changes need a token or SameSite + origin check? | A01 / CWE-352 |
| Crypto misuse | Weak algorithms, static IV/nonce, ECB, non-CSPRNG tokens, `==` on MACs, homemade crypto? | A02 / CWE-327, 330, 208 |
| Injection (SQL, NoSQL, OS, LDAP, template, header) | Is input concatenated into an interpreter instead of parameterised? | A03 / CWE-89, 78, 943, 1336 |
| XSS | Is output encoded for its context; any raw HTML sinks; CSP present? | A03 / CWE-79 |
| Race conditions / business logic | Check-then-act on balances, quotas, coupons, one-time tokens without locking or idempotency? | A04 / CWE-362, 367 |
| Security misconfiguration | Debug mode, permissive CORS with credentials, missing headers, verbose errors? | A05 / CWE-16, 942 |
| AuthN / session | Token validation (alg, aud, iss, exp), session rotation on login, logout revocation, reset-token entropy and expiry, rate limiting, enumeration? | A07 / CWE-287, 384, 307 |
| Insecure deserialization / integrity | Native deserialisers on untrusted data, unsigned webhooks, YAML/XML with unsafe loaders? | A08 / CWE-502, 345 |
| Secrets in code / sensitive logging | Hard-coded keys; tokens, passwords, PII or full requests logged? | A09, A02 / CWE-798, 532 |
| SSRF | Can user input choose host, scheme, port or redirect target of a server-side fetch? | A10 / CWE-918 |
| Path traversal / file handling | User input in file paths, archive extraction, upload names or content types? | A01 / CWE-22, 434 |

## Rules
- **Severity = impact x exploitability**, after existing controls:
  - Critical: unauthenticated RCE, SQLi, auth bypass, cross-tenant data access, secret exposure of production credentials.
  - High: authenticated user reaches another user's data or admin function, stored XSS in a privileged view, SSRF to internal network/metadata, account takeover chains.
  - Medium: reflected XSS needing interaction, CSRF on non-critical actions, user enumeration, missing rate limits on login, sensitive data in logs.
  - Low: defence-in-depth gaps with no demonstrated path (missing header, weak but unexploitable config).
  - Raise one level if pre-auth or internet-facing at scale; lower one level if it needs an admin or an unlikely precondition. Give a CVSS v3.1/4.0 vector only if the team uses CVSS.
- **False-positive discipline:** no finding without `file:line` for both source and sink and a stated attacker privilege. Check that the "missing" validation is not done in middleware, a decorator, a DB constraint, an ORM scope or the gateway before reporting. Parameterised queries, auto-escaping templates and framework CSRF protection are effective unless explicitly bypassed (`raw`, `safe`, `dangerouslySetInnerHTML`, `csrf_exempt`, string-built SQL).
- Authorisation is checked on the server for the specific object, not just "is logged in"; list endpoints and exports need the same tenant filter as detail endpoints.
- Prefer allowlists over denylists; validate type, length, range and format at the boundary, and encode at the sink. Sanitising on input for HTML is not a substitute for output encoding.
- Fixes go to the root cause in a shared layer (query builder, authZ policy, HTTP client wrapper) when the same pattern appears in several places; list every other occurrence you found with `rg`.
- Never claim the code is secure; state what was reviewed, what was not (e.g. infra, dependencies), and residual risk.
- Hand dependency and CI findings to `secure-supply-chain`; repository setting gaps to `secure-repo-governance`; general correctness issues to `review-pull-request`.

## Safety gates
- Do not run exploits, fuzzers or scanners against live or shared environments without written authorisation for that target; use local or ephemeral environments and test accounts.
- Never print, copy or commit secret values found during review: report location, type and first/last 4 characters at most, and recommend rotation (`secure-supply-chain` covers exposure response).
- Proof-of-concept payloads go in private findings or tests, not in public issues or PR comments for unfixed vulnerabilities. Ask before posting review comments or opening issues.

## Output
```markdown
## Security review: <scope> (commit/PR, date)
**Reviewed:** entry points, files | **Not reviewed:** ... | **Tools run:** ...

| ID | Severity | Class (OWASP/CWE) | Location (sink) | Source -> sink trace | Attacker / precondition | Impact | Fix | Test |
|---|---|---|---|---|---|---|---|---|
| S1 | High | A01 IDOR / CWE-639 | `api/invoices.py:88` | `GET /invoices/<id>` path param -> `Invoice.get(id)` no owner filter | any logged-in user | read other tenants' invoices | scope query by `request.user.org_id` | `test_invoice_idor` |

**Hardening notes (not proven exploitable):** ...
**Residual risk / follow-ups:** ...
```

---

**Phase:** Secure. **Next:** last phase; feed findings back into Plan: `/anthropica-engineering:plan-architecture`, `/anthropica-engineering:plan-codebase-analysis`. Phase map: `/anthropica-engineering:workflow`.
