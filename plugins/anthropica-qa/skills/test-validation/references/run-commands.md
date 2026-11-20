# Run commands per stack

Always prefer the repo's own scripts and the exact CI invocation (flags, env vars, services). Use these to narrow a run, repeat it, randomise order or measure coverage. Check the tool version in the lockfile before relying on a newer flag.

## Detect the stack
| File | Stack / runner |
|---|---|
| `package.json` (`scripts.test`, devDependencies) | Jest, Vitest, Mocha, node:test, Playwright, Cypress |
| `pyproject.toml`, `pytest.ini`, `tox.ini`, `noxfile.py`, `manage.py` | pytest, unittest, Django |
| `go.mod` | go test |
| `pom.xml` / `build.gradle(.kts)` | Maven Surefire/Failsafe, Gradle (JUnit, Kotest) |
| `*.csproj`, `*.sln` | dotnet test (xUnit, NUnit, MSTest) |
| `Cargo.toml` | cargo test / nextest |
| `Gemfile`, `.rspec` | RSpec, Minitest |
| `composer.json`, `phpunit.xml` | PHPUnit, Pest |
| `mix.exs` | ExUnit |
| `Package.swift` | swift test |
| `Makefile`, `justfile`, `Taskfile.yml` | Wrapper targets: read them, they often set env vars |

## JavaScript / TypeScript
| Goal | Command |
|---|---|
| Vitest one test | `npx vitest run path/to/file.test.ts -t "test name"` |
| Vitest shuffle | `npx vitest run --sequence.shuffle` |
| Jest one test | `npx jest path/to/file.test.ts -t "test name"` |
| Jest serial / random order | `npx jest --runInBand` / `npx jest --randomize --seed=1234` |
| Jest changed only | `npx jest --changedSince=origin/main` |
| Coverage | `npx vitest run --coverage` / `npx jest --coverage` |
| Playwright one test, repeated | `npx playwright test tests/x.spec.ts:42 --repeat-each=30 --workers=1 --retries=0` |
| Playwright debug artefacts | `npx playwright test --trace on` then `npx playwright show-trace <zip>` |
| Cypress one spec | `npx cypress run --spec cypress/e2e/x.cy.ts` |
| node:test | `node --test --test-name-pattern="name"` |
| Type-check / lint | `npx tsc --noEmit` / `npx eslint .` |
| Mutation | `npx stryker run --mutate "src/changed/file.ts"` |

Set `CI=true` to disable watch mode and interactive reporters.

## Python
| Goal | Command |
|---|---|
| One test | `pytest tests/test_x.py::TestClass::test_method -vv` |
| By expression / stop on first fail | `pytest -k "transfer and not slow" -x` |
| Re-run last failures | `pytest --lf` |
| Repeat (pytest-repeat) | `pytest tests/test_x.py::test_y --count=50 -x` |
| Random order (pytest-randomly) | `pytest --randomly-seed=1234` (auto-enabled once installed); disable with `-p no:randomly` |
| Parallel (pytest-xdist) | `pytest -n auto` (also exposes shared-state bugs) |
| Coverage with branches | `pytest --cov=pkg --cov-branch --cov-report=term-missing --cov-report=xml` |
| Changed-lines coverage | `diff-cover coverage.xml --compare-branch=origin/main` |
| unittest | `python -m unittest pkg.tests.test_mod.TestClass.test_method` |
| Django | `python manage.py test app.tests.TestClass.test_method --parallel 1` |
| Lint / types | `ruff check .` / `mypy .` or `pyright` |
| Mutation | `mutmut run` (scope via its config to changed modules) |

## Go
| Goal | Command |
|---|---|
| One test, no cache | `go test ./pkg/... -run '^TestName$' -count=1 -v` |
| Subtest | `go test ./pkg -run 'TestName/case_name' -count=1` |
| Race detector | `go test -race ./...` (run in CI too) |
| Repeat / shuffle | `go test -run '^TestName$' -count=50 -shuffle=on ./pkg` |
| Coverage | `go test -coverprofile=cover.out ./... && go tool cover -func=cover.out` |
| Vet / lint | `go vet ./...` / `golangci-lint run` |

`-count=1` is required to bypass the test cache when inputs outside Go files changed.

## JVM
| Goal | Command |
|---|---|
| Maven one method | `mvn -pl module -Dtest='ClassTest#method' -Dsurefire.failIfNoSpecifiedTests=false test` |
| Maven integration test | `mvn -pl module -Dit.test=ClassIT verify` |
| Gradle one method | `./gradlew :module:test --tests 'com.example.ClassTest.method' --rerun` |
| Coverage | JaCoCo report task (`mvn verify` or `./gradlew jacocoTestReport`) |
| Mutation | `mvn org.pitest:pitest-maven:mutationCoverage -DtargetClasses=com.example.changed.*` |

## .NET
| Goal | Command |
|---|---|
| One test | `dotnet test --filter "FullyQualifiedName~Namespace.Class.Method"` |
| Verbose output | `dotnet test --logger "console;verbosity=detailed"` |
| Coverage | `dotnet test --collect:"XPlat Code Coverage"` |
| Mutation | `dotnet stryker` |

## Rust
| Goal | Command |
|---|---|
| One test | `cargo test -p crate_name test_name -- --exact --nocapture` |
| Serial | `cargo test -- --test-threads=1` |
| nextest (retries visible as flaky) | `cargo nextest run -E 'test(test_name)'` |
| Coverage | `cargo llvm-cov --html` |
| Lint / mutation | `cargo clippy --all-targets -- -D warnings` / `cargo mutants --file src/changed.rs` |

## Ruby, PHP, Elixir, Swift
| Goal | Command |
|---|---|
| RSpec one example | `bundle exec rspec spec/models/x_spec.rb:42` |
| RSpec random order / find polluter | `bundle exec rspec --order rand --seed 1234` / `--bisect` |
| Minitest (Rails) | `bin/rails test test/models/x_test.rb:12 --seed 1234` |
| PHPUnit | `vendor/bin/phpunit --filter testName tests/XTest.php`; random: `--order-by=random` |
| ExUnit | `mix test test/x_test.exs:42 --seed 0`; repeat: `--repeat-until-failure 50` (recent versions) |
| Swift | `swift test --filter ModuleTests.ClassTests/testMethod` |

## Generic repeat loop
```bash
for i in $(seq 1 30); do <test command> || { echo "failed on run $i"; break; }; done
```
```powershell
1..30 | ForEach-Object { <test command>; if ($LASTEXITCODE -ne 0) { "failed on run $_"; break } }
```

## Reporting
Record the exact command, working directory, relevant env vars (names only, never secret values), tool version and the summary line (passed/failed/skipped/duration). Quote the first failure's assertion message verbatim.
