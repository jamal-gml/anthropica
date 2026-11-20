# Framework-specific gotchas

Check the framework version first: defaults changed across major versions.

## Node.js / Express
- No CSRF, security headers or body size limits by default: look for `helmet`, CSRF middleware (or SameSite + origin check), `express.json({ limit })`.
- `req.query` can be arrays/objects (`?id[]=1`, `?a[$ne]=x`) - type-check before passing to Mongo/Mongoose (NoSQL injection) or string functions.
- `res.sendFile(userPath)` without `root` option; `path.join(base, userInput)` with absolute input.
- Prototype pollution: deep-merge of user JSON (`lodash.merge`, custom merges) into objects; `__proto__`, `constructor.prototype` keys.
- `trust proxy` misconfigured makes `req.ip` and `req.secure` client-controlled.
- Error handler leaking stack traces when `NODE_ENV` is not `production`.

## Next.js / React
- `dangerouslySetInnerHTML`, `href={userUrl}` (javascript: URLs).
- Server Actions and Route Handlers are public HTTP endpoints: each needs its own authN/authZ; do not rely on the page that calls them being protected. Middleware-only auth can be bypassed on paths the matcher misses.
- Secrets exposed through `NEXT_PUBLIC_*` env vars or props serialised to the client (`getServerSideProps` returning whole DB objects).
- `next/image` remote patterns too broad (SSRF/cost); open redirects via `redirect(searchParams.next)`.

## Django
- Safe by default: ORM params, template auto-escape, CSRF middleware. Flag `raw()`, `extra()`, `RawSQL`, `cursor.execute` with formatting, `|safe`, `mark_safe`, `@csrf_exempt`, `autoescape off`.
- `DEBUG = True`, `ALLOWED_HOSTS = ['*']`, `SECRET_KEY` in code, `SESSION_COOKIE_SECURE`/`CSRF_COOKIE_SECURE` false, `SECURE_PROXY_SSL_HEADER` set without a proxy that strips it.
- DRF: `permission_classes` default (`AllowAny` unless set globally), `get_queryset` not filtered by user, serializers with `fields = '__all__'` (mass assignment and over-exposure), `read_only_fields` missing for owner/role.

## Flask / FastAPI
- Flask: `render_template_string` with user input (SSTI), `send_file` with user paths, `SECRET_KEY` weak, no CSRF without Flask-WTF, debug mode exposes the Werkzeug console (RCE).
- FastAPI: auth via `Depends` must be on every router or included at router level; Pydantic models used for both input and output leak fields (use separate response models); `response_model` missing; CORS `allow_origins=["*"]` with `allow_credentials=True`.

## Ruby on Rails
- Strong parameters: `permit!` or permitting `role`/`admin`/`*_id`. `html_safe`, `raw`, `render inline:`, `send_file params[...]`, `constantize`/`send` on user input, `find_by_sql`/`where("... #{x}")`.
- `protect_from_forgery` disabled; `skip_before_action :authenticate_user!` on sensitive controllers.
- Run `brakeman -q` for a quick pass.

## Spring (Java/Kotlin)
- Spring Security `permitAll()` matchers broader than intended, matcher ordering, `mvcMatchers` vs `antMatchers` path differences (trailing slash, suffix).
- CSRF disabled globally for a cookie-session app. Actuator endpoints exposed (`management.endpoints.web.exposure.include=*`), especially `/env`, `/heapdump`.
- SpEL evaluation of user input; JPQL/native queries built with concatenation; data binding directly to JPA entities (mass assignment).
- Jackson default typing enabled (`activateDefaultTyping`, `@JsonTypeInfo(use = CLASS)`) - deserialization gadgets.

## ASP.NET Core
- `[AllowAnonymous]` on controllers or actions overriding policy; missing `[Authorize]` with no fallback policy.
- `@Html.Raw`, `FromSqlRaw` with interpolation (use `FromSqlInterpolated` or parameters), model binding to EF entities, `TypeNameHandling` in Newtonsoft.Json.
- Antiforgery not validated for cookie-auth APIs; CORS `AllowAnyOrigin` with credentials (blocked at runtime but indicates intent).

## Go (net/http, gin, echo)
- `text/template` instead of `html/template` for HTML; `fmt.Sprintf` into SQL; `filepath.Join` with user input then no prefix check; `http.ServeFile` with request paths.
- `InsecureSkipVerify: true`; default `http.Client` with no timeout (DoS); `math/rand` for tokens; missing `http.MaxBytesReader`.
- Goroutine data races on shared maps (run tests with `-race`).

## PHP / Laravel
- `DB::raw`, `whereRaw` with interpolation, `{!! !!}` Blade output, `$fillable`/`$guarded = []` mass assignment, `unserialize` on user data, `APP_DEBUG=true` in production, file uploads stored under `public/`.

## ORMs in general
- Raw query escape hatches; dynamic `order by`/column names from user input (cannot be parameterised - allowlist them); lazy-loaded relations serialised into responses; soft-deleted records visible through relations.

## Mobile
- Android: exported activities/receivers/providers, `WebView.setJavaScriptEnabled` with `addJavascriptInterface` loading remote content, cleartext traffic allowed, secrets in `strings.xml`/`BuildConfig`.
- iOS: ATS exceptions, secrets in plist, data in `UserDefaults` instead of Keychain, URL scheme handlers trusting input.
