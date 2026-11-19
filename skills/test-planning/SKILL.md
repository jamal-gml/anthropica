---
name: test-planning
description: Create risk-based test plans, test cases, exploratory test charters and high-quality bug reports for features and releases. Use when planning QA for a feature or release, writing test cases from requirements, organising regression testing, or reporting a defect.
---

# Test Planning (QA)

## Test plan
1. **Scope:** features and changes in scope and explicitly out of scope; the requirements or acceptance criteria being tested.
2. **Risk analysis:** rank areas by likelihood × impact (payments, authentication, data integrity and new integrations rank high). Test effort follows risk.
3. **Approach per level:** unit and integration (developer-owned), API, end-to-end UI, exploratory, non-functional (performance, accessibility, security, compatibility).
4. **Environments and data:** which environment, test accounts and roles, seeded data, feature-flag states. Never use real customer personal data.
5. **Coverage matrix:** browsers, devices and OS versions, locales, user roles.
6. **Entry and exit criteria:** for example "all critical and high cases pass, no open SEV1/SEV2 bugs, regression suite green".
7. **Schedule, owners and reporting cadence.**

## Test cases
```
ID / Title:        <behaviour under test, not "test 1">
Preconditions:     <state, role, data, flags>
Steps:             1. ... 2. ...
Expected result:   <observable outcome per step or at the end>
Priority:          P1–P3     Type: functional | negative | boundary | regression
Traceability:      <requirement or ticket ID>
```
- Cover the happy path, negative paths (invalid input, permission denied), boundaries (empty, max length, zero, limits), state transitions and concurrency where relevant.
- Use equivalence partitioning and boundary-value analysis to keep the suite small but meaningful.
- Mark cases worth automating (stable, repeated, high value) and keep exploratory testing for new or changing areas.

## Exploratory charter
"Explore <area> with <resources/approach> to discover <risk or information>." Timebox it (60–90 min) and record notes, questions and bugs.

## Bug report
```
Title:        <what is wrong, where> (specific, searchable)
Environment:  build/version, URL, browser/device, account role
Steps:        minimal numbered steps to reproduce (reproduce it twice)
Expected:     ...
Actual:       ...  (exact error text; screenshot, video, logs, request ID)
Severity:     impact on users   Priority: urgency to fix
Frequency:    always / intermittent (x of y)
```
One bug per report. Check for duplicates first. No blame, just facts.

## Output
The plan, cases or report in the formats above, plus a risk summary and a coverage-gaps list.
