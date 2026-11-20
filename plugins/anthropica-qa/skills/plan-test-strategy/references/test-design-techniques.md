# Test-design techniques with worked examples

Use these to keep a suite small but able to find defects. Combine them: partition first, add boundaries, then cover rules and states.

## Equivalence partitioning (EP)
Split each input into classes the system should treat the same; test one value per class. Always include invalid classes.

Example - "age must be 18 to 65 inclusive, whole number":
| Class | Example | Expected |
|---|---|---|
| Valid 18-65 | 40 | accepted |
| Below range | 10 | rejected with message |
| Above range | 80 | rejected |
| Non-integer | 30.5 | rejected |
| Non-numeric | "abc" | rejected |
| Empty / null / missing field | "" | "required" error |

Rules:
- One invalid class per case. Combining two invalid values hides which one the system rejected.
- Partition outputs too (e.g. each discount tier, each error code), not only inputs.
- Hidden partitions often come from the implementation: storage type limits, currency precision, locale formats, pagination size.

## Boundary value analysis (BVA)
Defects cluster at edges. For each ordered partition test the boundary and its neighbours.
- 2-value BVA: min, max and the first invalid on each side (17, 18, 65, 66).
- 3-value BVA (High risk): min-1, min, min+1, max-1, max, max+1.
- Lengths: 0, 1, max-1, max, max+1 characters; check multi-byte characters count the same in UI, API and database column.
- Collections: 0, 1, many, page size, page size + 1, the configured maximum.
- Money: smallest unit (0.01), rounding half cases (0.005), negative, zero, currency with 0 or 3 decimals.
- Dates and times: start/end of day, month, year; 29 Feb; DST switch days; UTC vs local midnight; expiry exactly at "now".
- Numeric types: 0, -1, max int, overflow, NaN/Infinity where floats are accepted.

## Decision tables
Use when outcomes depend on combinations of conditions. List conditions as rows, one column per rule, then one test per column.

Example - free shipping:
| Condition / Rule | R1 | R2 | R3 | R4 |
|---|---|---|---|---|
| Member | Y | Y | N | N |
| Order >= 50 | Y | N | Y | N |
| **Free shipping** | Y | Y | Y | N |

- Collapse columns whose outcome does not depend on a condition (mark "-") to reduce cases, but keep at least one case per distinct outcome.
- Missing or contradictory columns are requirement gaps: raise them as questions.
- Cause-effect graphing is the same idea for larger rule sets; generate the table, do not test the cartesian product.

## State-transition testing
Use for anything with a status: orders, subscriptions, accounts, documents, payments.

1. Draw states and events. Example order: Draft -> Submitted -> Paid -> Shipped -> Delivered; Submitted/Paid -> Cancelled; Paid -> Refunded.
2. Build the state table (state x event -> next state or "invalid").
3. Coverage levels:
   - All states visited.
   - **0-switch:** every valid transition once (minimum for Medium risk).
   - **1-switch:** every valid pair of consecutive transitions (High risk, e.g. Paid -> Cancelled -> refund issued).
   - **Invalid transitions:** each event in each state where it is not allowed (e.g. Ship a Cancelled order) must be rejected and leave state unchanged.
4. Check guards (conditions on transitions), side effects (emails, ledger entries, stock) and idempotency (same event twice, concurrent events).

## Pairwise (all-pairs) testing
Most interaction defects involve two parameters. Pairwise covers every pair of values with far fewer cases than the full product.

Example: 3 browsers x 3 OS x 4 locales x 3 roles = 108 full combinations; pairwise needs about 12-16.
- Generate with a tool rather than by hand: PICT, allpairspy, or the generator built into your test-management tool.
- Declare constraints (e.g. Safari only on macOS/iOS) so impossible rows are excluded.
- Add explicit cases for known risky triples and for the most-used production combination.
- Use 3-wise only for High-risk configuration spaces; it grows quickly.

## Scenario and use-case testing
For each user goal: main success flow, alternate flows (other valid paths), exception flows (each failure point). Include interruptions: back button, refresh, timeout, closed tab, network loss, payment declined, concurrent edit by another user.

## CRUD and data-lifecycle matrix
For each entity: Create, Read, Update, Delete x each role, plus: duplicate create, update of a deleted item, delete with dependents, soft vs hard delete, audit trail, export and data-subject deletion. Check that data written by one feature reads correctly in every other feature that uses it.

## Permission matrix
Rows = roles (anonymous, user, owner, admin, other tenant, suspended account); columns = sensitive actions. Every cell has an expected allow/deny and a case. Test via API directly, not only by hidden UI buttons.

## Error guessing checklist
Empty, whitespace-only, leading/trailing spaces, unicode, emoji, combining characters, RTL text, HTML/script-like text, SQL-like quotes, very long strings, zero-width characters, pasted formatted text, duplicate submit, slow network, timeouts, offline, retry after partial failure, expired session, clock skew, different timezone, leap day, high-latency third party, third-party error responses (4xx, 5xx, malformed body).

## Choosing quickly
1. List inputs and outputs -> EP.
2. Ordered partitions -> add BVA.
3. Rules combine conditions -> decision table.
4. A status field exists -> state transitions.
5. Configuration matrix -> pairwise.
6. Always finish with error guessing and a permission check.
