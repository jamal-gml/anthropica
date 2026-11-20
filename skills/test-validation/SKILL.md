---
name: test-validation
description: Choose, write and run the smallest tests that prove a change - regression test first for bugs, repo run commands, flaky-test diagnosis and quarantine, mocking boundaries, coverage and mutation signals. Use after implementing a change, when fixing a bug, or when tests are missing, failing or flaky. Not for release test planning (use `plan-test-strategy`) or load tests (use `test-performance`).
---

# Testing and Validation

Prove the requested behaviour works, with the fewest, most reliable tests, and report exactly what ran.

## Procedure
1. **State the change.** One sentence for the behaviour changed, one per acceptance criterion, and what result would prove it wrong.
2. **Find how this repo tests.** The CI workflow is the source of truth for commands, flags, env vars and services; then `package.json` scripts, `Makefile`, `pyproject.toml`/`tox.ini`/`noxfile.py`, `go.mod`, `pom.xml`/`build.gradle`, `*.csproj`, `Cargo.toml`, `Gemfile`. Read the nearest existing tests for naming, fixtures, factories and helpers and follow them. Read `references/run-commands.md` for single-test, repeat and coverage commands per stack.
3. **Pick the level** with the Rules table: the smallest test that can fail for the right reason.
4. **Bugs: test first.** Write the regression test, run it, and confirm it fails *for the expected reason* (assertion on the buggy behaviour, not an import, setup or typo error). Then fix, then see it pass. Name it after the behaviour, link the issue in a comment if the repo does that.
5. **Features:** one test per acceptance criterion, plus boundaries, invalid input, permission denied and the error path of each external call.
6. **Run narrow to wide:** single test -> file -> package/module -> full suite, then lint, type-check and build exactly as CI does. Do not claim done on a narrow run.
7. **On failure:** read the first failure in full (message, diff, stack), not the summary. Reproduce it alone. Classify: product bug, test bug, or environment. Fix the root cause; never paper over it.
8. **Prove the test bites.** Temporarily revert the fix or invert the key condition; the new test must fail. Restore and re-run. This is a one-off manual mutation test.
9. **Check stability** of new or touched tests: repeat 20-50 times and once in random order (commands in `references/run-commands.md`). Read `references/flaky-tests.md` when any test passes and fails without code changes.
10. **Report** with the Output template.

## Rules

Level selection:
| Change | Smallest proving test |
|---|---|
| Pure logic, parsing, calculation | Unit |
| Query, migration, ORM mapping, serialisation | Integration against a real DB/broker in a container |
| HTTP handler, authz, validation | API/integration test through the router |
| Contract with another service | Consumer-driven contract test |
| UI component behaviour | Component test querying by role/label, not CSS selectors |
| Critical cross-system journey | One E2E test, only if lower levels cannot prove it |
| Config, build, infra | The build/validate/plan command CI runs |

Mocking:
- Mock at boundaries you do not own: third-party HTTP, payment, email/SMS, clock, randomness, OS. Wrap them in your own adapter; mock the adapter in callers and test the adapter separately against a sandbox or recorded responses.
- Do not mock what you own: use real collaborators, in-memory fakes or a containerised real dependency. Mocking your ORM or repository hides SQL and mapping bugs.
- Over-mocking smells: asserting call order of internals, tests that break on a pure refactor, more mock setup than assertions, mocks returning mocks.
- Prefer fakes and stubs (state) over interaction verification, except for "it must send exactly one email/charge" type outcomes.

Test quality:
- Arrange-Act-Assert; one behaviour per test; name says behaviour and condition (`rejects_transfer_above_daily_limit`).
- Deterministic: injected/frozen clock, seeded randomness, no `sleep` (wait for a condition with a timeout), no real network.
- Independent: each test creates its own data and cleans up; no order dependence; no shared mutable globals.
- Assert outcomes, specifically: exact values and error types, not `toBeTruthy`, "no exception", or a broad snapshot of a whole object.
- Snapshots only for stable rendered output; review every snapshot diff line by line.

Coverage and mutation:
- Coverage is a signal, not a goal. Look at **changed-lines coverage**: every uncovered changed branch is a question to answer ("is it tested elsewhere or dead?"), not a number to raise.
- 100% line coverage with weak assertions proves nothing. For critical logic run mutation testing (Stryker for JS/TS/C#, PIT for JVM, mutmut for Python, cargo-mutants for Rust, Gremlins for Go) on the changed files only; each surviving mutant is a missing assertion.
- Do not lower coverage thresholds or add exclusions to pass CI.

Never, to get green: delete, skip, `xfail`/`@Ignore`/`.only` a test; loosen assertions; change expected values to match actual without proving actual is correct; raise timeouts or add retries blindly; bulk-update snapshots or golden files; catch and swallow errors in tests. If a test is obsolete because the requirement changed, say so and get confirmation before changing it.

Gotchas:
- A test that passes before the fix is not a regression test.
- Watch modes and caches lie: use the CI-style one-shot run (`CI=true`, `-count=1`, `--no-cache` where relevant).
- Tests that pass locally but fail in CI usually differ in timezone, locale, parallelism, file-system case sensitivity, env vars or service versions; compare them.
- Async tests that never `await` the assertion pass vacuously; make sure the test can fail.

## Safety gates
- Before running suites that migrate, truncate or seed databases, check the connection env vars point to a local or disposable test resource, never shared or production.
- Never run tests with production credentials or against live payment, email or SMS providers.
- Ask before deleting, skipping or quarantining tests, bulk-updating snapshots or baselines, changing CI test configuration or coverage thresholds, or adding test dependencies.
- Scrub secrets and personal data from fixtures, recorded HTTP cassettes and logs before committing.

## Output
```
Change: <behaviour>   Acceptance: <criteria>
Tests added/changed:
- <path>::<test name> - proves <behaviour>; failed before fix: yes (<reason>) | n/a
Commands run (exact) and results:
- <command> -> <N passed, M failed, K skipped, duration>
- lint / type-check / build -> <result>
Bite check: <what was reverted/inverted> -> test failed as expected
Stability: <test> x<N> runs, random order -> <all passed | x of N failed>
Coverage of changed lines: <%> ; uncovered: <file:line - reason>
Not tested / residual risk: <...>   Pre-existing failures (unrelated): <test - evidence>
```
