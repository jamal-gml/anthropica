# Reviewing a pull request

Review the change against its stated intent. Prioritise what is expensive to fix later: correctness, data, security, contracts. Hand deep security review to `secure-code-review` and test design to `test-validation` when the change warrants it.

## Checklist by pass
**1. Scope**
- Matches the linked issue; no unrelated refactors, drive-by renames or debug leftovers.
- Title type matches the change (a behaviour change titled `refactor` hides it from release notes).

**2. Correctness**
- Edge cases: empty, null, zero, max size, unicode, time zones, DST, leap days, concurrent calls, retries, partial failure.
- Error handling: errors surfaced or logged with context, not swallowed; no broad catch that hides bugs.
- Resource handling: connections, files, locks released; timeouts set on network calls.
- Idempotency of handlers, jobs and migrations that can be retried.
- Off-by-one, wrong comparison operators, inverted conditions, missing `await`.

**3. Tests**
- New behaviour has a test that fails without the change. Bug fixes include a regression test.
- Tests assert outcomes, not implementation details; no sleeps, no order dependence, no real network.
- Coverage of the failure paths, not just the happy path.

**4. Security and data**
- Input validated at trust boundaries; output encoded; parameterised queries.
- Authorisation checked on every new endpoint or action, not just authentication.
- No secrets, tokens or PII in code, logs, error messages or test fixtures.
- New dependencies: maintained, licence acceptable, pinned via lockfile (`secure-supply-chain`).

**5. Compatibility and rollout**
- Public API, CLI, config, event or schema changes: backwards compatible, or marked breaking with migration notes.
- Database migrations follow expand/contract and work with the previous app version (`build-database`, `ship-ci-cd`).
- Feature flag or config gate for risky behaviour; defaults safe.
- Observability: logs/metrics for the new path (`operate-observability`).

**6. Maintainability**
- Names say what things are; functions do one thing; no duplicated logic that already exists in the codebase.
- Comments explain why, not what. Docs and changelog updated if user-visible.

## Comment templates
```
blocking: <problem>. <consequence>. Suggest: <fix>.   (path/file.ts:42)
suggestion: <alternative> would <benefit>. Fine to keep as is.
question: What happens when <case>? I could not find a test for it.
nit: <tiny style point>. Non-blocking.
praise: <specific thing done well>.
```
Suggestion block (author can apply in one click):
````
```suggestion
const timeoutMs = config.timeoutMs ?? 5000;
```
````

## Posting with gh
- Whole review: `gh pr review <n> --request-changes --body-file review.md` (or `--approve`, `--comment`).
- Inline comments need the API: `gh api repos/{owner}/{repo}/pulls/<n>/reviews -f event=COMMENT -f body=... -F 'comments[][path]=src/a.ts' -F 'comments[][line]=42' -F 'comments[][body]=...'`.
- Ask before posting anything under the user's account.

## Verdict rules
- **Request changes:** at least one `blocking` item.
- **Approve:** no blocking items; suggestions and nits may remain. Approving with nits is normal.
- **Comment:** questions that must be answered before you can judge.
- Do not approve what you did not read. If you reviewed only part, say which part.
- Re-review only the delta after updates: `gh pr diff <n>` against the last reviewed commit, or the "changes since last review" view.

## Etiquette
- Critique the code, not the author. Be specific; link docs or prior art.
- Batch comments into one review rather than many notifications.
- Respond to the author's replies; resolve your own threads once satisfied.
- Turnaround target: first review within one working day.
