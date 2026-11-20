# Flaky tests: diagnosis, fixes and quarantine

A flaky test passes and fails on the same code. It erodes trust in the whole suite, so treat it as a bug in either the test or the product. Concurrency flakes are often real product bugs.

## 1. Confirm and measure
- Re-run the single test 30-50 times on the same commit (`references/run-commands.md`). Record the failure rate (e.g. 3/50).
- Run it alone, then with its file, then with the full suite. Passing alone but failing in the suite = order dependence or shared state.
- Run in random order with a fixed seed; when it fails, keep the seed to reproduce.
- Compare CI vs local: parallelism, CPU, timezone, locale, env vars, service versions, file-system case sensitivity, network access.
- Collect artefacts on failure: full logs, timestamps, traces/screenshots/videos (UI), DB state, the seed.

## 2. Classify by symptom
| Symptom | Likely cause | Fix |
|---|---|---|
| Fails under load or in CI only, timeouts | Fixed `sleep`, waiting less than the operation takes, slow CI runners | Wait on a condition (element visible, message received, row exists) with a generous ceiling; use framework auto-waiting; never add bare sleeps |
| Fails only after another test | Shared state: DB rows, globals, singletons, env vars, module-level caches, mocks not restored, files in a shared temp dir | Isolate data per test (transaction rollback, unique IDs, per-test schema), restore mocks/env in teardown, unique temp dirs; find the polluter with bisect (`rspec --bisect`, pytest-randomly seed + halving) |
| Fails around midnight, month end, DST, in another timezone | Real clock, local timezone, date arithmetic | Inject or freeze the clock; set `TZ=UTC` in tests and add explicit timezone cases |
| Fails with certain generated data | Random data without a seed, unique-constraint collisions | Seed generators and log the seed; generate unique values per test |
| Fails on assertion of order | Unordered collections (sets, maps, DB without `ORDER BY`, parallel results) | Sort before comparing or assert as a set; add `ORDER BY` in product code if order matters |
| Port in use, file locked, "too many open files" | Hard-coded ports, leaked resources, parallel workers sharing resources | Bind to port 0, close resources in teardown, per-worker resources |
| Intermittent network errors | Real external calls | Stub the boundary; for integration, use local containers |
| Fails rarely with wrong result under parallel requests | Race condition in product code | Reproduce with the race detector or stress loop; fix the product, keep the test |
| Float comparisons off by a tiny amount | Floating point | Compare with tolerance; use decimal types for money |
| UI element "detached" or "not clickable" | Re-render during interaction, animation, overlay | Query by role at action time, wait for stable state, disable animations in test config |
| Fails only on first run after deploy/start | Cold start, lazy init, cache warm-up, migrations still running | Wait for a readiness signal before starting tests |

## 3. Fix, then prove
- Fix the root cause, then re-run the repeat loop (same count) with zero failures, plus one random-order run of the whole file.
- If the cause was product code (race, missing ordering, timezone bug), report it as a product bug and keep the test as the regression guard.

## 4. Quarantine rules
Quarantine only when the fix cannot land quickly and the flake is blocking others. It is a temporary, visible state, not a delete.
- Ask before quarantining; name an owner and an issue.
- Use the framework's skip/tag mechanism with a reason and issue link, e.g. a `@flaky`/`quarantine` tag that is excluded from the blocking job but still runs in a non-blocking job so failures stay visible.
- Set an expiry (e.g. 14 days). At expiry: fixed, or escalated, or deleted with an explicit decision that the behaviour is covered elsewhere.
- Never quarantine a test that guards security, money, data integrity or a critical journey without an equivalent guard; escalate instead.
- Track the quarantine list size; a growing list is a process problem to report.

## 5. Retries
- Automatic retries in CI hide flakes. If the repo uses them, make sure retried passes are reported as "flaky", not as plain passes, and review that report.
- Never add retries to a specific test as the fix.

## Report template
```
Test: <path::name>
Failure rate: <x/N> locally, <x/N> in CI (links)
Mode: alone pass/fail | in suite pass/fail | random order seed <n>
Root cause: <category + evidence>
Fix: <test or product change>  Verification: <command> x<N> -> 0 failures
Quarantine: none | <tag, issue, owner, expiry>
```
