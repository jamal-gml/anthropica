# Frontend performance and security

## Core Web Vitals targets (75th percentile of real users)
| Metric | Good | Typical causes when bad |
|---|---|---|
| LCP (largest contentful paint) | <= 2.5 s | slow server/TTFB, render-blocking CSS/JS, unoptimised or lazy-loaded hero image, client-only rendering |
| INP (interaction to next paint) | <= 200 ms | long tasks on main thread, heavy re-renders, large hydration, synchronous work in handlers |
| CLS (cumulative layout shift) | <= 0.1 | images/embeds without dimensions, injected banners, late web fonts, skeletons that differ in size |

Measure, do not guess: Lighthouse/DevTools Performance panel locally, field data (RUM, CrUX) for real users. Deeper investigation: `test-performance`.

## Bundle and loading
- Check the cost before adding a dependency (bundle analyser or a package size tool); prefer native APIs (`Intl`, `fetch`, `URL`, `structuredClone`) and small single-purpose modules.
- Import only what you use (named imports from ESM packages; avoid whole-library imports of utility and icon libraries).
- Code-split at route level; lazy-load heavy, below-the-fold or rarely used components (editors, charts, maps) with dynamic `import()`.
- Set a performance budget (e.g. initial JS gzipped per route) and check it in CI if the project has one.
- Avoid shipping large JSON or translations for all locales up front; load per route/locale.
- Third-party scripts (analytics, chat, ads): load `async`/`defer` or after interaction; each one is a performance and security liability.

## Images and media
- Explicit `width` and `height` (or `aspect-ratio`) on every image and embed.
- Modern formats (AVIF/WebP) with `srcset`/`sizes`, or the framework image component.
- Hero/LCP image: never lazy-load; add `fetchpriority="high"`; preload if discovered late.
- Below the fold: `loading="lazy"` and `decoding="async"`.
- SVG icons inline or as a sprite; avoid icon fonts.

## Fonts
- Self-host or preconnect; `font-display: swap` or `optional`; subset; preload only the 1-2 critical files; match fallback metrics (`size-adjust`) to cut CLS.

## Rendering
- Avoid unnecessary re-renders: keep state close to where it is used, split contexts, memoise only after profiling shows a problem.
- Virtualise long lists (hundreds+ rows).
- Debounce search inputs (~200-300 ms) and cancel in-flight requests.
- Break long tasks (> 50 ms) with yielding (`scheduler.yield`, `setTimeout`), web workers for CPU-heavy work.
- Animate `transform` and `opacity`, not layout properties.
- SSR/SSG for content pages and SEO; streaming and server components where the framework supports them; hydrate only interactive parts.

## Caching
- Hashed asset filenames with long-lived `Cache-Control: public, max-age=31536000, immutable`; HTML `no-cache`.
- Service workers only with a clear update strategy; a stale service worker can pin users to old code.

## Security
### XSS
- Frameworks escape text by default. Danger zones: `dangerouslySetInnerHTML`, `v-html`, `[innerHTML]`, `{@html}`, `innerHTML`, `outerHTML`, `insertAdjacentHTML`, `document.write`, `eval`, `new Function`, `setTimeout(string)`.
- If HTML from users or a CMS must render: sanitise with DOMPurify (or the framework's sanitiser) on the client at render time, with a restrictive allowlist; server-side sanitising alone is not enough if content can be changed elsewhere.
- URLs: allow only `http:`, `https:`, `mailto:` (and relative) in `href`/`src`; parse with `new URL()` and check `protocol`.
- `target="_blank"` links to external sites: add `rel="noopener noreferrer"` (modern browsers imply noopener, but be explicit for older ones).
- Never build HTML, CSS or script strings from user input; never pass user input to `postMessage` targets without checking `event.origin` on receipt.

### Content Security Policy
- Prefer a nonce- or hash-based CSP with `strict-dynamic`; avoid `unsafe-inline` and `unsafe-eval` for scripts.
- Set `object-src 'none'`, `base-uri 'none'` (or `'self'`), `frame-ancestors` to prevent clickjacking.
- Roll out with `Content-Security-Policy-Report-Only` first and watch reports.
- Consider Trusted Types (`require-trusted-types-for 'script'`) to block DOM XSS sinks.

### Tokens and sessions
| Storage | Risk | Use |
|---|---|---|
| `HttpOnly; Secure; SameSite=Lax/Strict` cookie | not readable by JS; needs CSRF protection for state-changing requests if SameSite is not enough | default for session/refresh tokens |
| In-memory variable | lost on reload; safe from persistent theft | short-lived access tokens |
| `localStorage` / `sessionStorage` | readable by any XSS, persists | avoid for tokens |
- With cookie auth: CSRF tokens or SameSite plus Origin/Referer checks on the server.
- OAuth in SPAs: authorization code flow with PKCE; never the implicit flow; consider a backend-for-frontend that keeps tokens server-side.
- Clear client caches and stored user data on logout.

### Secrets and data exposure
- Anything in the bundle is public, including env vars prefixed for client use. API keys in the client must be publishable keys restricted by domain/scope.
- Source maps in production: upload privately to error tracking or restrict access if they reveal sensitive code.
- Do not log tokens or personal data to the console or client-side error tracking; scrub before sending.
- Client-side permission checks only hide UI; the server must enforce authorisation.

### Dependencies
- Lockfile committed; review new transitive dependencies; audit regularly (`secure-supply-chain`).
- Subresource Integrity (`integrity` + `crossorigin`) for scripts and styles loaded from third-party CDNs.
