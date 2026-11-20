# Exploratory testing: charters, sessions and debriefs

Exploratory testing is simultaneous learning, test design and execution. Use it for new features, high-risk areas, vague requirements and after scripted tests pass. Structure it with session-based test management so it is plannable and reportable.

## Charter
Format: **Explore** <target> **with** <resources, data, tools, technique> **to discover** <information or risk>.

Good charters are narrow enough for one session and name a risk:
- Explore checkout with expired and declined test cards to discover how payment failures are surfaced and whether the order is left half-created.
- Explore bulk import with 10k-row and malformed CSV files to discover limits, partial-failure handling and error reporting.
- Explore the role switcher with two browser sessions to discover stale permissions after a role change.
- Explore the mobile signup flow with the network throttled and interrupted to discover duplicate accounts and lost input.

Bad: "Explore the app", "Test checkout" (no target risk, unbounded).

## Session
- Timebox: 60-90 minutes, uninterrupted. Short session 45, long 120.
- One tester (or a pair: driver + navigator).
- Keep running notes; do not stop to write formal cases.
- Track the split: % on charter, % on opportunity (useful detours), % setup. High setup % signals environment or data problems worth fixing.

Session sheet:
```
Charter: <...>
Tester: <...>   Date: <...>   Build: <...>   Duration: <min>
Split: charter <%> / opportunity <%> / setup <%>
Areas and data covered: <...>
Notes: <what was tried, observed, assumed>
Bugs: <IDs, one line each>
Issues / questions: <blockers, unclear requirements, test-data gaps>
Follow-up charters: <...>
```

## Heuristics to generate ideas
**SFDIPOT (product elements):**
| Element | Ask |
|---|---|
| Structure | What is it made of: files, services, config, dependencies? |
| Function | What does each feature do, including error handling and hidden features? |
| Data | What inputs, outputs, defaults, large, empty, corrupt, legacy data? |
| Interfaces | APIs, UI, imports/exports, integrations, CLI? |
| Platform | Browsers, OS, devices, runtime versions, third parties it depends on? |
| Operations | How is it really used: frequent paths, power users, misuse, admin tasks? |
| Time | Concurrency, timeouts, scheduling, timezones, sequences, long-running sessions? |

**Other quick heuristics:**
- Goldilocks: too small, too big, just right.
- CRUD: create, read, update, delete every object; then do them out of order.
- Interrupt: back, refresh, close, sleep the device, lose the network, double-click.
- Follow the data: enter a value and trace it through every screen, export, email and report.
- Starve it: low disk, slow CPU, no permissions, blocked third party.
- Opposite user: novice, expert, keyboard-only, screen-reader user, malicious user.
- Configuration toggles: every setting and flag, then combinations.
- Consistency oracles (when unsure if behaviour is a bug): compare with the product's history, similar features in the product, comparable products, the claims in docs and marketing, user expectations, the spec, and law or standards.

**Tours:**
- Feature tour: touch every control once.
- Money tour: the flows that make or move money.
- Landmark tour: hop between key features in different orders.
- Bad-neighbourhood tour: areas with the most past bugs.
- Saboteur tour: break dependencies, revoke access, corrupt data mid-flow.
- Back-alley tour: least-used features and settings.

## Debrief (10-15 minutes after each session)
Use PROOF:
- **Past:** what happened during the session?
- **Results:** what was achieved, bugs found?
- **Obstacles:** what got in the way (environment, data, access)?
- **Outlook:** what still needs testing; new charters?
- **Feelings:** how confident is the tester about this area?

Feed results back: new charters into the backlog, repeatable checks into automation candidates, requirement gaps to the product owner, and a confidence rating per area into the test report.

## Reporting exploratory coverage
| Area | Sessions | Bugs (S1/S2/S3/S4) | Confidence (low/med/high) | Next charter |
|---|---|---|---|---|
Report coverage by risk area, not number of sessions alone.
