# Pipeline speed and monorepos

## Find the time
- Take p50 and p90 per job over the last 50 runs; optimise the critical path (longest chain of `needs`), not the sum.
- Separate queue time (waiting for runners) from run time; queue time is a capacity problem.
- Track flake rate per test; a 2% flaky suite fails a 50-test-file pipeline often.

## Levers, in order of usual payoff
| Lever | How | Typical gain |
|---|---|---|
| Dependency caching | cache keyed on lockfile hash; restore-keys fallback | 30-90 s per job |
| Parallel stages | lint, unit, build, scan as sibling jobs | critical path shrinks to the longest one |
| Test sharding | split by timing data into N shards (`--shard`, `pytest-split`, `knapsack`) | near-linear until setup dominates |
| Affected-only builds | run only changed projects and their dependents | large in monorepos |
| Build cache | remote cache (Bazel, Nx, Turborepo, Gradle, BuildKit `cache-to`) | rebuilds skip unchanged targets |
| Container layer order | copy lockfile and install deps before source copy; multi-stage | image builds from minutes to seconds |
| Shallow clones | depth 1 except for versioning/changelog jobs | seconds to minutes on large repos |
| Right-size runners | bigger runners for CPU-bound builds/tests | 2-4x on heavy jobs |
| Fail fast | cheap checks first; cancel superseded runs on the same branch | saves queue capacity |
| Move slow suites | full e2e, perf, fuzzing to merge queue, nightly or pre-release | PR feedback under 10 min |

Do not "speed up" by skipping security scans, lowering coverage gates or retrying flakes automatically.

## Monorepo strategy
- **Detect affected projects** with the build graph (Nx `affected`, Turborepo `--filter=...[origin/main]`, Bazel query, Pants `--changed-since`, Gradle/Maven module diffs) rather than raw path globs where possible; graphs catch shared-library dependents.
- **Path filters** are acceptable for loosely coupled folders; include shared paths (`libs/`, root lockfile, CI config) in every filter.
- **Required checks:** skipped workflows never report. Use one always-running aggregator job as the required status (`ship-github-actions` has the pattern).
- **Per-project pipelines** share templates; versions and releases per package (`ship-release`).
- **Base for comparison:** PRs diff against the merge base with the target branch; pushes to main diff against the previous successful main run, not `HEAD~1` (squash merges and multi-commit pushes break `HEAD~1`).
- Run the full suite on a schedule or before release to catch graph mistakes.

## Merge queue
- Use when main breaks because PRs pass individually but conflict semantically.
- The queue tests the PR rebased on the queue head; required checks must run on the merge-queue event too.
- Keep queue checks under ~15 min or throughput collapses.

## Caching gotchas
- Caches are not artifacts: never rely on a cache for correctness; a miss must still build.
- Key on lockfiles and tool versions; include OS and architecture.
- PR caches should not write into the default branch scope (cache poisoning).
- Prune: oversized caches slow restore more than they save.

## Reliability
- Pin tool and runtime versions (`.tool-versions`, `.nvmrc`, `go.mod` toolchain, container images by digest).
- Hermetic tests: no shared external services; spin dependencies in containers.
- Timeouts on every job and network call.
- Quarantine flaky tests with an owner and deadline; report flake rate weekly.
