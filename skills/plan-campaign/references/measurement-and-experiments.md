# Measurement and experiments

## UTM convention
| Parameter | Meaning | Rule | Example |
|---|---|---|---|
| `utm_source` | Who sent the traffic (platform, publisher, list) | lowercase, no spaces | `newsletter`, `partner-name`, `search-engine` |
| `utm_medium` | Channel type, from a fixed vocabulary | only values from the list below | `cpc`, `paid-social`, `email` |
| `utm_campaign` | Campaign ID | `yyyy-qn-<slug>` or `yyyymm-<slug>` | `2026-q2-spring-launch` |
| `utm_content` | Creative/variant/placement | encode the test variant | `hero-video-a`, `cta-button-b` |
| `utm_term` | Paid keyword or audience | optional | `audience-lookalike-1pct` |

- Fixed `utm_medium` vocabulary: `cpc`, `paid-social`, `display`, `video`, `email`, `sms`, `organic-social`, `affiliate`, `referral`, `partner`, `qr`, `print`. Analytics tools group channels by medium; off-list values land in "unassigned/other".
- Everything lowercase, hyphens not spaces or underscores, no personal data (no emails or names in URLs).
- Never tag internal links (site nav, banners to own pages): it restarts the session and overwrites the real source.
- Keep one shared UTM sheet (link, source, medium, campaign, content, owner, date); generate links from it, don't hand-type.
- Test every tagged link before launch: it resolves, survives redirects (some redirects strip query strings), and the session appears with correct source/medium in real-time analytics.
- Use a link shortener only after tagging; QR codes and offline media need their own `utm_medium` and a vanity URL.

## Conversion tracking checklist
- [ ] Each KPI has one source of truth (analytics, CRM, billing) and a definition (what counts, dedup rule, attribution window).
- [ ] Primary conversion events fire once per conversion (check duplicates on thank-you page refresh).
- [ ] Server-side or offline conversion import set up where browser tracking is unreliable.
- [ ] CRM captures first-touch and last-touch source fields on lead creation; hidden form fields populated from UTMs.
- [ ] Consent mode/banner behaviour tested: what is recorded when consent is refused.
- [ ] Revenue/value passed with conversions where ROAS is a KPI.
- [ ] A test conversion run end to end and seen in every reporting tool.

## Attribution caveats
| Model | Over-credits | Use for |
|---|---|---|
| Last click | Brand search, retargeting, email, direct | Short, single-session purchases |
| First click | Awareness channels | Understanding acquisition entry points |
| Linear / position-based | Spreads credit by rule, not evidence | Rough multi-touch view |
| Data-driven (platform) | Opaque; limited to tracked touchpoints | Within-platform optimisation |
| Platform-reported | Every platform claims the same conversion | Never sum across platforms |

- Sum of platform-reported conversions usually exceeds actual conversions. Reconcile against the source of truth.
- View-through conversions inflate display/video results; report them separately with a short window.
- Consent rejection, ad blockers, cross-device journeys and walled gardens mean tracked paths are a sample, not the truth.
- Long sales cycles: report leading indicators weekly, revenue attribution at cycle length (and say so).
- For causal claims ("this channel drove X"), use incrementality, not attribution.

## Incrementality
- **Audience holdout** (email, CRM, retargeting): randomly hold out 10-20% of the eligible audience, send nothing. Lift = (CVR_exposed - CVR_holdout) / CVR_holdout. Incremental conversions = (CVR_exposed - CVR_holdout) x exposed count.
- **Geo test** (paid media, offline): split matched regions into test and control (similar pre-period trends), change spend only in test, compare against forecast from control. Needs several weeks and enough regions; plan it with the media owner.
- **Platform lift studies:** useful if available; still check the design (randomisation unit, minimum audience).
- **Media mix modelling:** for larger, multi-channel budgets with 2+ years of weekly data; complements, does not replace, experiments.
- Report iROAS = incremental revenue / spend, and cost per incremental conversion, next to platform-reported figures.

## Experiment design
1. Hypothesis: "Changing <one variable> for <audience> will increase <primary metric> because <reason>."
2. Primary metric (one) and guardrail metrics (unsubscribe rate, refund rate, lead quality).
3. Minimum detectable effect (MDE): the smallest lift worth acting on.
4. Sample size per arm, two-sided alpha 0.05, power 80%, approximation for conversion rates:
   `n ~= 16 x p x (1 - p) / d^2` where p = baseline rate, d = absolute MDE.
   - Baseline 3%, detect 3.0% -> 3.6% (d = 0.006): n ~= 16 x 0.03 x 0.97 / 0.000036 ~= 12,900 per arm.
   - Baseline 20% open rate, detect 20% -> 22% (d = 0.02): n ~= 16 x 0.2 x 0.8 / 0.0004 = 6,400 per arm.
   Use a proper calculator for multiple arms or non-binary metrics.
5. Duration = (n x arms) / daily eligible traffic, rounded up to whole weeks (minimum one full weekly cycle, ideally two) to absorb day-of-week effects.
6. If the required sample is unreachable, raise the MDE, test a bolder change, pick a higher-funnel metric, or don't test - say so.

**Rules**
- One variable per test; everything else identical (audience, timing, budget).
- Randomise at the unit you analyse (user, not session; region for geo tests).
- Fix sample size and end date before starting; no stopping when it "looks significant". Use sequential methods only if pre-planned.
- Check sample ratio mismatch (arms differ from planned split by more than chance) before reading results.
- Report effect size with confidence interval, not just "won".
- A non-significant result means "no detectable effect at this MDE", not "no effect".
- Email subject-line tests on small lists: pick the winner by clicks or conversions when possible; opens are distorted by privacy features that pre-load images.

## Reporting cadence
- Daily during the first week (tracking and pacing sanity only, no optimisation decisions on tiny samples).
- Weekly: spend vs plan, KPI tree actuals vs targets, test status, decisions taken.
- End of campaign: post-campaign review.
