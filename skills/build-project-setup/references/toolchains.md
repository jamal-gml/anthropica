# Toolchains per ecosystem

Defaults for new projects. If the repo or team already uses a tool, keep it. Ask before adding any of these to an existing project.

## Quick table
| Ecosystem | Version pin | Package manager / lockfile | Format | Lint | Types | Test |
|---|---|---|---|---|---|---|
| JS/TS (Node) | `.nvmrc` + `engines`, `packageManager` field | pnpm `pnpm-lock.yaml` / npm `package-lock.json` | Prettier or Biome | ESLint (flat config) or Biome | `tsc --noEmit` (`strict: true`) | Vitest or Jest; Playwright for e2e |
| Python | `.python-version`, `requires-python` | uv `uv.lock` or Poetry `poetry.lock` | Ruff format | Ruff | mypy or pyright | pytest |
| Go | `go` + `toolchain` in `go.mod` | modules `go.sum` | `gofmt`/`goimports` | `go vet`, golangci-lint | compiler | `go test ./...` (`-race` in CI) |
| Rust | `rust-toolchain.toml` | Cargo `Cargo.lock` (commit for binaries) | `cargo fmt` | `cargo clippy -- -D warnings` | compiler | `cargo test` |
| Java/Kotlin | toolchain in Gradle/Maven | Gradle wrapper + lockfile/version catalog, Maven wrapper | Spotless (google-java-format/ktlint) | Checkstyle/PMD/detekt, Error Prone | compiler | JUnit 5 |
| .NET | `global.json` | NuGet `packages.lock.json` (`RestorePackagesWithLockFile`) | `dotnet format` | analyzers + `TreatWarningsAsErrors` in CI | compiler, `<Nullable>enable</Nullable>` | xUnit/NUnit |
| Ruby | `.ruby-version` | Bundler `Gemfile.lock` | RuboCop (or Standard) | RuboCop | Sorbet/RBS optional | RSpec/Minitest |
| PHP | `composer.json` `require.php` | Composer `composer.lock` | PHP-CS-Fixer / Pint | PHPStan/Psalm | PHPStan level | PHPUnit/Pest |

Use `.tool-versions` (asdf/mise) when the repo mixes several runtimes.

## Init commands (run in an empty directory)
- Node: `pnpm init` (or `npm init -y`), then add `typescript`, `@types/node`, `tsx`/`vitest` as dev deps on approval; `npx tsc --init`.
- Python: `uv init <name>` (creates `pyproject.toml`), `uv add --dev ruff pytest mypy`.
- Go: `go mod init <module-path>`.
- Rust: `cargo new <name>` (or `--lib`).
- .NET: `dotnet new sln`, `dotnet new <template> -o src/<Name>`, `dotnet new gitignore`, `dotnet new editorconfig`.
- Java: Spring Initializr or `gradle init`; commit the wrapper (`gradlew`, `gradle/wrapper/`).

Framework generators (`create-next-app`, `create-vite`, `django-admin startproject`, `rails new`) overwrite files: only run them in an empty directory.

## .editorconfig baseline
```ini
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 2

[*.{py,rs,java,kt,cs,php}]
indent_size = 4

[{Makefile,*.go}]
indent_style = tab

[*.md]
trim_trailing_whitespace = false
```

## .gitattributes baseline
```
* text=auto eol=lf
*.{cmd,bat} text eol=crlf
*.ps1 text eol=crlf
*.{png,jpg,jpeg,gif,ico,pdf,zip,gz,woff,woff2} binary
```
Line-ending normalisation on an existing repo rewrites many files: ask first and do it in its own commit.

## .gitignore essentials
Start from the standard language template, then make sure these are covered: dependency dirs (`node_modules/`, `.venv/`, `vendor/` if not vendored), build output (`dist/`, `build/`, `target/`, `bin/`, `obj/`), caches (`.pytest_cache/`, `.ruff_cache/`, `.next/`, `coverage/`), env files (`.env`, `.env.*`, then `!.env.example`), OS/editor (`.DS_Store`, `Thumbs.db`, `.idea/`, `.vscode/*` with `!.vscode/extensions.json` if shared), `CLAUDE.local.md`, `.claude/settings.local.json`.

## Pre-commit hooks
| Option | Use when |
|---|---|
| `pre-commit` framework (`.pre-commit-config.yaml`) | Python or polyglot repos; pin each hook `rev:` |
| husky + lint-staged | JS/TS repos already on npm/pnpm |
| lefthook | polyglot, fast, single binary |

Hooks run only fast checks on staged files (format, lint, secret scan such as gitleaks). Tests and typecheck of the whole project belong in CI. Never document `--no-verify` as a workaround.

## Scripts to standardise
Expose the same verbs in every repo so humans, CI and agents share them: `setup`/`install`, `dev`, `build`, `test`, `test:unit`, `lint`, `format`, `format:check`, `typecheck`, `check` (all of the fast gates). Use `package.json` scripts, a `Makefile`/`justfile`, or the ecosystem task runner.

## README minimum
Title and one-line purpose; prerequisites with versions; setup; run locally; test; configuration (env var names, link to `.env.example`); project layout pointers (short); how to contribute or release (link). Licence: ask the user, do not choose one.

## CI minimum
Install with the frozen lockfile (`npm ci`, `pnpm install --frozen-lockfile`, `uv sync --locked`, `go mod download` + `go mod verify`, `cargo build --locked`, `dotnet restore --locked-mode`), then format check, lint, typecheck, tests. Workflow authoring: `ship-github-actions`. Dependency update bots and pinning policy: `secure-supply-chain`.
