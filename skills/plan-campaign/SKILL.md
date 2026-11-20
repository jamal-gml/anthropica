---
name: plan-campaign
description: Plan marketing campaigns and launches - objective and KPI tree, segments, message house, channel mix by funnel stage, budget and pacing, timeline, UTMs, measurement, experiments and post-campaign review. Use when planning a launch, campaign, content calendar or media budget, or defining KPIs, attribution or A/B tests for marketing. Not for writing the copy (use `create-marketing-copy`).
---

# Campaign Planning

Turn a business goal into a measurable, budgeted, testable campaign plan that someone else can execute and evaluate.

## Procedure
1. **Intake.** Get, or ask once for: business goal, product/offer, deadline, budget ceiling, target market, existing baselines (traffic, conversion rates, CAC, list sizes), brand and legal constraints, past campaign results. List anything missing as an assumption rather than guessing silently.
2. **Objective and KPI tree.** One primary objective, SMART (specific, measurable, owner, target, date). Decompose it backwards into driver metrics so every channel KPI rolls up to it:
   ```
   Revenue 250k by 30 Jun
   = closed deals (50) x average deal (5k)
     deals    = SQLs (250) x win rate (20%)
     SQLs     = MQLs (1,000) x MQL->SQL (25%)
     MQLs     = landing sessions (25,000) x CVR (4%)
     sessions = paid (15k) + email (6k) + organic/social (4k)
   ```
   Use the client's real baselines for each ratio; if a rate must improve to hit the target, say so explicitly - that is the plan's main risk.
3. **Audience and segments.** Primary segment (and at most two secondary): job/need, trigger event, objections, awareness stage, where they spend attention, size estimate and source. Exclusions (existing customers, competitors, employees) matter as much as targets.
4. **Positioning and message house.** Positioning statement: for [segment] who [need/trigger], [product] is the [category] that [primary benefit], unlike [main alternative], because [reason to believe]. Then a message house:
   - Roof: one umbrella message (one sentence, benefit-led).
   - Pillars: 2-3 supporting messages, one per key benefit or objection.
   - Foundation: proof per pillar (data, customer evidence, certification, demo). A pillar with no proof is marked `[NEEDS PROOF]`, not filled with invented evidence.
5. **Channel mix by funnel stage.** Pick channels where the segment already is and that serve the objective. Map each to a stage and a KPI (table in Rules). Plan one hero asset and its channel derivatives.
6. **Budget and pacing.** Allocate, set guardrails and the reallocation rule before launch (see Rules).
7. **Timeline.** Back-plan from launch date: tracking QA, creative production, legal/compliance review, audience build, soft launch, launch, optimisation checkpoints, wrap-up. Include owners, approval steps and dependencies (product readiness, pricing, landing page).
8. **Measurement plan.** UTM convention, conversion events, source of truth for each KPI, attribution model and its caveats, incrementality method, reporting cadence. Read `references/measurement-and-experiments.md` when defining UTMs, attribution, holdouts or sample sizes.
9. **Experiments.** One hypothesis per test, one variable changed, primary metric, minimum detectable effect, sample size and run time computed before launch.
10. **Risks and contingencies.** Pre-agree what you cut, double down on or pause at each checkpoint (for example "if CPA > 1.3x target after 2 weeks and 30+ conversions, shift 20% budget to best ad group").
11. **Review.** Schedule the post-campaign review now; use `references/post-campaign-review.md`.

## Rules
**Channel mix by funnel stage**
| Stage | Typical channels | Primary KPI | Diagnostic metrics |
|---|---|---|---|
| Awareness | Video, display, paid social reach, PR, sponsorships, podcasts | Reach at target frequency, branded search lift | CPM, view rate, frequency |
| Consideration | Search (non-brand), content/SEO, webinars, retargeting, social | Engaged sessions, content conversions | CTR, CPC, time engaged, scroll depth |
| Conversion | Brand search, retargeting, email, landing pages, sales outreach | Conversions, CPA/CAC, ROAS | CVR, form drop-off, cost per lead |
| Retention/expansion | Lifecycle email, in-product, community, customer marketing | Activation, repeat rate, churn, NRR | Open/click, feature adoption |

- Owned (site, email, blog), earned (PR, community, partners, reviews), paid (search, social, display, video, affiliates). Paid buys speed; owned/earned compound. Do not plan paid-only for a long-term goal.
- **Budget split default:** ~70% proven channels, ~20% promising, ~10% experiments; adjust with evidence. Reserve 10-15% unallocated for mid-flight reallocation.
- **Minimum viable spend:** a paid channel or ad set that cannot reach enough conversions per week to exit the platform's learning phase (check the platform's guidance; often a few dozen per week) is noise. Consolidate rather than spreading thin.
- **Pacing:** daily target = remaining budget / remaining days. Check weekly; flag deviations beyond +-10%. Avoid major edits to bids/targeting/creative more than about weekly on algorithmic campaigns - each edit can reset learning.
- **Frequency:** watch frequency and creative fatigue (falling CTR at stable reach). Rotate creative before CTR drops ~30% from its first-week level.
- **Targets come from math, not hope.** If the funnel math needs a CVR far above baseline, either raise traffic, lower the target or flag it.
- **Benchmarks:** use the client's own history first. Never present an industry benchmark without a cited source; mark unsourced figures as estimates.
- **Message consistency:** every asset maps to one pillar. Ad, landing page and email for the same audience share the headline promise (message match lowers bounce).
- **Tracking before traffic:** no paid spend until conversion events fire correctly in a test, UTMs resolve, and the CRM records source fields.
- **Consent:** track only with valid consent where law requires; plan for measurement loss from consent rejection and blocked cookies rather than ignoring it. Email only to opted-in or lawfully contactable lists, with unsubscribe.
- **Attribution humility:** platform-reported conversions overlap and over-count; last-click over-credits brand search, retargeting and email. Treat channel ROAS as directional unless backed by an incrementality test.
- **One variable per test,** pre-registered stopping rule, no peeking-driven early stops.

## Safety gates
- Ask before: launching or increasing paid spend, committing budget or signing media/vendor contracts, sending emails or SMS, publishing posts or press releases, contacting press, partners or influencers, uploading customer lists to ad platforms.
- Uploading customer data (custom audiences) needs a lawful basis and consent; confirm before planning it.
- No fabricated stats, benchmarks, testimonials, customer logos or results in the plan or its proof points. Mark unknowns `[ASSUMPTION]` or `[NEEDS DATA]`.
- Regulated sectors (finance, health, gambling, alcohol, children): flag required approvals and disclaimers; route copy through compliance before launch.

## Output
```
# Campaign brief: <name>
Objective: <metric> from <baseline> to <target> by <date>   Owner: <name>
KPI tree: <primary> <- <drivers with baseline -> target>
Audience: <primary segment, trigger, objections>; exclusions: <...>
Positioning: <statement>
Message house: roof | pillars (1-3) | proof per pillar ([NEEDS PROOF] where missing)
Channel plan:
| Channel | Stage | Audience | Asset | Budget | KPI + target | Owner |
Budget: total, split, reserve, pacing rule, reallocation trigger
Timeline: <phase - dates - owner - approval - dependency>
Measurement: UTM convention, events, source of truth, attribution model + caveats, incrementality test
Experiments: hypothesis | variable | metric | MDE | sample/arm | duration
Risks and contingencies: <trigger -> action>
Assumptions and missing data: <list>
Approvals needed before launch: <spend, sends, publishes>
```
Follow with the detailed content calendar and tracking sheet when asked.
