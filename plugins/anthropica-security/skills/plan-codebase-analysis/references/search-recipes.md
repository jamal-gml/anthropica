# Search recipes

Use with ripgrep (`rg`) or the Grep tool. Add `-g '!**/node_modules/**' -g '!**/dist/**'` style excludes when the tool does not respect `.gitignore`.

## General patterns
| Goal | Pattern / command |
|---|---|
| Exact message, fixed string | `rg -F 'Order not found'` |
| Symbol definition | `rg -n '\b(def|function|func|class|interface|type|const|let|var)\s+Name\b'` |
| All references, whole word | `rg -nw Name` |
| Message built from parts | search the stable fragment, then the i18n/message key |
| Env var reads | `rg -n 'ORDER_TIMEOUT'` and check `.env.example`, config loaders, deployment manifests |
| SQL table usage | `rg -n -i '\border_items\b'` plus ORM model names |
| Files changed together | `git log --name-only --format= -- path/file \| sort \| uniq -c \| sort -rn \| head` |
| Hot spots (churn) | `git log --since=6.months --name-only --format= \| sort \| uniq -c \| sort -rn \| head -20` |
| Who knows this area | `git shortlog -sn -- path/` |

If a string is not found: it may be generated (i18n file, enum-to-text), come from the server/DB, or be split across lines or template literals. Search a shorter fragment and the translation files.

## Entry points and routing by stack
| Stack | Entry / wiring | Route or handler registration |
|---|---|---|
| Node/Express, Fastify, Koa | `package.json` `main`/`scripts.start`, `src/index|server|app.*` | `router.get(`, `app.post(`, `fastify.route(` |
| NestJS | `src/main.ts`, `*.module.ts` | `@Controller(`, `@Get(`, `@Post(`; providers in modules |
| Next.js | `app/` or `pages/` directory | `app/**/page.tsx`, `app/**/route.ts`, `pages/api/**`, `middleware.ts`, server actions `'use server'` |
| Python Django | `manage.py`, `settings.py` `INSTALLED_APPS` | `urls.py` `path(`, views, `signals.py`, `admin.py`, middleware list |
| Python FastAPI/Flask | `main.py`/`app.py`, `uvicorn`/`gunicorn` command | `@app.get(`, `@router.post(`, `APIRouter(`, `Blueprint(`, `Depends(` |
| Go | `cmd/*/main.go` | `http.HandleFunc(`, `mux.Handle(`, `r.GET(` (gin), `e.GET(` (echo), `chi.Router` |
| Java Spring | `@SpringBootApplication` class | `@RestController`, `@RequestMapping`, `@GetMapping`, `@KafkaListener`, `@Scheduled` |
| .NET | `Program.cs` | `app.MapGet(`, `[ApiController]`, `[HttpGet]`, `services.Add*` for DI |
| Ruby on Rails | `config/application.rb` | `config/routes.rb`, `app/controllers`, callbacks in models, `app/jobs` |
| PHP Laravel | `public/index.php`, `bootstrap/app.php` | `routes/web.php`, `routes/api.php`, service providers |
| Rust | `src/main.rs`, `src/bin/` | axum `Router::new().route(`, actix `#[get(` |
| Frontend SPA | `index.html`, `src/main.*` | router config (`createBrowserRouter`, `routes`), lazy imports |

Background work: search for `cron`, `schedule`, `@Scheduled`, `celery`, `sidekiq`, `bull`, `queue`, `consumer`, `subscribe`, `worker`, `handler` and job registries. CI and deployment files often reveal processes that do not appear in code (separate worker commands).

## Hidden callers
| Mechanism | How to find callers |
|---|---|
| DI container | search for the interface/token name and the binding (`bind`, `provide`, `AddScoped`, `@Injectable`) |
| Events / pub-sub | search for the event name string or constant on both `emit/publish` and `on/subscribe` sides |
| Decorators / annotations | search for the decorator name; registration often happens at import time |
| ORM hooks | `before_save`, `pre_save`, `@PrePersist`, `beforeCreate`, triggers in migrations |
| DB triggers / procedures | search migrations and `*.sql` for `CREATE TRIGGER`, `CREATE FUNCTION` |
| Reflection / string dispatch | search the string key; look for maps of name -> handler |
| Templates | search template files for the helper or partial name |
| Infra-level | API gateway routes, rewrites in reverse-proxy config, serverless function definitions |

## Git archaeology
```
git log --oneline -20 -- path/to/file          # recent intent for a file
git log -S'exact string' --oneline             # commits adding/removing a string
git log -G'regex' --oneline -- path/            # commits whose diff matches a regex
git log -L :functionName:path/to/file           # history of one function
git blame -w -C -C -L 120,160 path/to/file      # ignore whitespace, follow moves
git show <sha> --stat                           # what else changed with it
git log --follow -- path/to/file                # history across renames
```
Check `.git-blame-ignore-revs` and pass `--ignore-revs-file` if it exists, so bulk-format commits do not hide the real author.

## Tests as spec
- Find tests by name convention: `*_test.go`, `test_*.py`, `*.spec.ts`, `*.test.tsx`, `*Test.java`, `*Tests.cs`, `spec/**/*_spec.rb`.
- Fixtures and factories show valid data shapes; snapshot files show expected output.
- Run only the slice: `pytest path::test_name -q`, `go test ./pkg/... -run TestName`, `npx vitest run file`, `npx jest file -t name`, `mvn -Dtest=Class#method test`, `dotnet test --filter Name`.
