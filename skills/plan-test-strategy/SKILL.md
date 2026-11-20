---
name: plan-test-strategy
description: Plan risk-based testing - risk scoring, test levels, test-design techniques, exploratory charters, entry/exit criteria, test data, traceability and bug reports. Use when planning QA for a feature or release, deriving test cases from requirements, scoping regression, or writing a defect report. Not for writing or running tests (use `test-validation`) or load tests (use `test-performance`).
---

# Test Strategy and Planning

Decide what to test, at which level, how deeply and when testing is done, with effort proportional to risk.

## Procedure
1. **Gather the test basis.** Requirements, acceptance criteria, designs, API contracts, the diff or change list, past incidents and bug clusters in the area. Turn every ambiguity into a written question; a requirement nobody can test is a requirement defect.
2. **Set scope.** In scope, out of scope (with the reason) and assumptions. Include the blast radius of the change: shared components, migrations, config, feature flags, integrations, scheduled jobs.
3. **Score risk** per area: Impact (1-3) x Likelihood (1-3). Impact drivers: money, auth and permissions, data loss or corruption, privacy and legal, number of users, no workaround. Likelihood drivers: new or complex code, large diff, new technology or integration, bug history, unclear requirements, concurrency, time pressure. Map the score to depth with the Rules table.
4. **Choose levels.** Put each check at the lowest level that can catch the failure (Rules table). Pyramid (many unit, fewer integration, few E2E) for logic-heavy services; trophy (integration-heavy, static analysis at the base) for frontends and glue code; consumer-driven contract tests between independently deployed services.
5. **Design cases.** Derive at least one case per requirement, then add negative, boundary, state, permission and error-guessing cases using the technique table. Read `references/test-design-techniques.md` when applying partitioning, boundary values, decision tables, state transitions or pairwise.
6. **Plan exploratory sessions** for new, high-risk or poorly specified areas as timeboxed charters. Read `references/exploratory-testing.md` when writing charters, running sessions or debriefing.
7. **Decide non-functional coverage** per area: performance (`test-performance`), security (`plan-threat-model`, `secure-code-review`), accessibility (`review-ui-ux`), compatibility, resilience, upgrade/migration and rollback, observability. Write "N/A because ..." instead of skipping silently.
8. **Environments and data.** Environment per level and its parity gaps with production, accounts per role and tenant, synthetic seed data, flag states, third-party sandboxes or stubs, clock and timezone control, data reset between runs.
9. **Entry/exit criteria and traceability.** Build the matrix requirement -> cases -> result -> bugs. Every in-scope requirement has at least one case; every High-risk area has negative and boundary cases.
10. **Select regression.** Changed areas + their dependents + smoke of critical journeys. Full regression only for High-risk releases or when impact analysis is unreliable.
11. **Report** with the Output template; log defects with the bug template during execution.

## Rules

Risk score -> depth:
| Score | Depth |
|---|---|
| 6-9 High | All relevant techniques, negative + boundary + state cases, exploratory charter, automated regression, named sign-off owner |
| 3-4 Medium | Requirements-based cases, key boundaries and negative paths, automate if repeated |
| 1-2 Low | Happy path and smoke, rely on existing automation |

Level selection:
| Failure to catch | Lowest useful level |
|---|---|
| Calculation, parsing, branching rules | Unit |
| SQL, ORM mapping, migrations, serialisation, queue or cache wiring | Integration against the real dependency (container) |
| Request/response shape between services | Contract |
| Authorisation per endpoint and role | API / integration |
| Multi-step journey across UI and backend | E2E (critical journeys only) |
| Layout and visual regressions | Visual snapshot |
| Unknown unknowns, usability, odd combinations | Exploratory |

Technique selection:
| Input or behaviour | Technique |
|---|---|
| Ranges, lengths, quantities, dates | Equivalence partitions + boundaries (min-1, min, max, max+1, empty, null) |
| Combinations of conditions with different outcomes | Decision table, one case per rule column |
| Lifecycle or status field | State transition: every valid transition, invalid transitions, guards |
| Many independent parameters (browser x OS x locale x role) | Pairwise instead of the full cartesian product |
| Workflows | Scenario: main, alternate and exception flows |
| Anything | Error guessing: whitespace, unicode/emoji, RTL, very long input, double submit, back/refresh mid-flow, two tabs, session expiry, slow or offline network, timezone/DST, leap day, rounding |

Gotchas:
- E2E is the slowest and flakiest level. Keep it to critical journeys; never repeat in E2E what unit or integration tests already prove.
- Pass rate and test count are not coverage. Report requirement coverage and risk coverage.
- Test the permission matrix explicitly: each role x each sensitive action, plus another tenant's or user's resource (IDOR).
- Test upgrade and migration from realistic existing data and the rollback path, not only a fresh install.
- Feature flags: test on, off and the switch while users are mid-session.
- Severity (impact) and priority (urgency) are independent: a logo typo is S4/P1; a crash in a rarely used export may be S2/P3.

Default entry criteria: build deployed to the test environment, smoke passes, test data ready, requirements stable enough to test, blockers from the previous cycle fixed.
Default exit criteria (tighten for High risk): all High-risk cases executed and passed; no open S1/S2; every S3 triaged with an owner; 100% traceability for in-scope requirements; regression green; residual risks accepted in writing by the owner.

Severity scale (priority P1-P4 is set by the product owner from severity and business context):
| Severity | Meaning |
|---|---|
| S1 Critical | Data loss or corruption, security breach, outage, core journey blocked with no workaround |
| S2 Major | Core function broken with a painful workaround, wrong money or data shown |
| S3 Minor | Non-core function broken, or easy workaround |
| S4 Trivial | Cosmetic, copy |

## Safety gates
- Never use real customer personal data outside production. Use synthetic or irreversibly masked data; ask before copying any production data.
- Do not test against production, or with live payment, email or SMS providers, without explicit approval. Use sandboxes and allowlisted recipients.
- Confirm with the environment owner before destructive, data-reset or load tests on shared environments.
- Waiving exit criteria (shipping with open S1/S2) is a named human decision. Record it; never make it.
- Redact secrets, tokens and personal data from screenshots, logs and HAR files attached to bug reports.

## Output
```
# Test plan: <feature/release>   Build: <version/commit>   Owner: <name>
Scope: in <...> | out <... because ...> | assumptions <...>
Risks:
| Area | Impact | Likelihood | Score | Depth | Levels / techniques |
Approach: unit / integration / contract / API / E2E / exploratory / non-functional (or N/A because ...)
Environments and data: <env, accounts per role, data, flags, stubs>
Entry criteria: <...>   Exit criteria: <...>
Traceability: | Req ID | Case IDs | Result | Bugs |
Open questions and coverage gaps: <...>
```
Test case:
```
ID: TC-<area>-<n>   Title: <behaviour + condition, e.g. "Rejects transfer above daily limit">
Req: <ID>  Risk: High|Med|Low  Level: <unit/API/E2E>  Technique: <BVA/decision table/...>
Preconditions: <role, data, flags>
Steps: 1. ...  Expected: <observable outcome>
Automate: yes|no (<reason>)
```
Bug report (one defect per report; search for duplicates first; facts, no blame):
```
Title: <Component>: <what fails> when <condition>
Environment: build/commit, URL, browser/OS/device, role, flags
Steps: 1. ... (minimal, reproduced twice from a clean state)
Expected: ... (cite the requirement)   Actual: ... (exact error text, request ID, timestamp)
Frequency: always | x of y   Severity: S1-S4   Priority: P1-P4 (owner sets)
Evidence: screenshot/video, logs, HAR (redacted)   Notes: first bad build, workaround, related bugs
```
