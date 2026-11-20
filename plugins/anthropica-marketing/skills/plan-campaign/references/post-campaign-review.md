# Post-campaign review

Run within two weeks of the end date (or after one sales-cycle length for lead campaigns). Blameless: focus on decisions and evidence, not people.

## Template
```
# Post-campaign review: <campaign>  (<start> - <end>)
Prepared by: <name>   Data as of: <date>   Sources: <analytics, CRM, ad platforms, billing>

## 1. Result vs objective
Objective: <metric, target, date>
Actual: <value> (<% of target>)   Verdict: hit / partial / missed

## 2. KPI tree actuals
| Metric | Baseline | Target | Actual | Delta vs target | Note |
(every driver from the plan; flag the ratio that most explains the gap)

## 3. Spend and efficiency
| Channel | Planned | Actual | Conversions (source of truth) | CPA/CAC | ROAS | Incremental? |
Total spend vs budget; reallocations made and why.

## 4. Incrementality and tests
| Test | Hypothesis | Result (effect, CI) | Decision |
Holdout/geo results; how they compare with platform-reported numbers.

## 5. What worked / what did not
- Evidence-backed only; cite the metric for each point.

## 6. Surprises and data quality issues
- Tracking gaps, attribution conflicts, consent loss, broken links, late assets.

## 7. Decisions and actions
| Action | Owner | Due |
(keep / scale / stop channels; creative to reuse; process fixes)

## 8. Learnings to carry forward
- Updated baselines for the next plan (CVRs, CPAs, open/click rates).
```

## Rules
- Use the source-of-truth numbers (CRM, billing) for outcomes; platform numbers only for in-platform diagnostics.
- Separate correlation from causation: only incrementality tests support "caused" language.
- Compare against the pre-agreed targets and contingency triggers, not against moved goalposts.
- Record updated baselines; they feed the next `plan-campaign` run.
- Note lagging effects (pipeline still open, brand lift) and when to re-check them.
