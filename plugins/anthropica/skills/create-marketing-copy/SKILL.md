---
name: create-marketing-copy
description: Write and edit conversion-focused marketing copy (landing pages, emails, ads, social, announcements, push/SMS) that is on-brand, accurate, accessible and testable. Use when drafting, rewriting or A/B-testing customer-facing marketing text, headlines, subject lines or CTAs. Not for search-led articles (use `create-seo-content`) or campaign strategy (use `plan-campaign`).
---

# Marketing Copy

Produce copy that moves one specific reader to one specific action, with every claim supportable.

## Procedure
1. **Brief intake.** Get, or ask once for, everything below; proceed with stated assumptions if the user wants speed.
   | Field | Question |
   |---|---|
   | Audience | Who, their trigger, main objection, awareness stage (unaware / problem-aware / solution-aware / product-aware / most-aware) |
   | Goal | The single action (sign up, book, buy, reply, read) and what happens after the click |
   | Offer | What exactly, price/terms, deadline (real?), risk reversal (trial, guarantee) |
   | Proof | Data, customer quotes (with permission), case studies, awards, certifications - with sources |
   | Channel | Format, placement, character limits, device (mobile first?), localisation |
   | Voice | Brand voice guide, banned words, examples of approved copy |
   | Constraints | Legal/regulatory rules, required disclaimers, sector restrictions |
2. **Message.** Write the one-sentence promise (benefit + for whom + proof). If you can't, the brief is incomplete - ask.
3. **Pick the framework** from the table in Rules by awareness stage and length.
4. **Draft** the body, then the headline (headlines are easier once the body is known), then the CTA.
5. **Fit the channel.** Check limits and truncation points in `references/channel-limits.md`. Front-load meaning before the truncation point.
6. **Edit pass.** Cut 20-30%: filler, hedges, adverbs, repeated ideas. Read aloud. Check message match with the ad or email that leads to the page.
7. **Claims pass.** Tag every factual claim; mark anything without a source `[VERIFY: what is needed]`.
8. **Accessibility and inclusion pass** (Rules).
9. **Variants.** Produce 2-3 variants of the highest-impact element (headline, subject line, hook or CTA), each testing one stated hypothesis.

Read `references/formats.md` when writing a landing page, email, ad set, social post, announcement, push or SMS - it has the structure for each.

## Rules
**Framework choice**
| Framework | Structure | Best for |
|---|---|---|
| AIDA | Attention, Interest, Desire, Action | Longer pages and emails for cold or solution-aware readers |
| PAS | Problem, Agitate, Solution | Problem-aware readers; short ads, cold email, pain-driven offers |
| BAB | Before, After, Bridge | Transformation stories, social posts, case-study teasers |
| FAB | Feature, Advantage, Benefit | Product pages and feature announcements for product-aware readers |
| 4Us (headline check) | Useful, Urgent, Unique, Ultra-specific | Score headlines/subject lines 1-4 on each; rewrite anything scoring low on Useful or Ultra-specific |

- **Awareness stage drives the opening:** unaware -> story or surprising fact; problem-aware -> name the problem in their words; solution-aware -> why this approach beats alternatives; product-aware -> proof and offer; most-aware -> the deal and the CTA.
- Lead with the reader's outcome; features are evidence for benefits. Use the customer's vocabulary (from reviews, support tickets, sales calls) over internal jargon.
- Specific beats vague: numbers, timeframes, named use cases. Cut empty superlatives ("innovative", "cutting-edge", "seamless", "world-class", "revolutionary", "unlock", "elevate").
- One idea and one primary CTA per piece. Secondary CTAs only for a different stage (for example "See pricing" next to "Start trial").
- **CTA rules:** verb + outcome, 2-5 words ("Start my free trial", "Get the checklist"), says what happens next, never "Submit" or "Click here". First-person wording often helps; test it. Put a friction reducer next to it only if true ("No card required", "Cancel anytime").
- Short sentences (average under ~20 words), short paragraphs (1-3 lines on mobile), active voice, scannable subheads that work alone.
- Voice: follow the guide; if none exists, state the voice used in 3 adjectives plus "not" counterparts (for example "confident, not arrogant").
- **Message match:** headline on the landing page repeats the promise of the ad/email that sent the click.
- **Urgency and scarcity only when real** (actual deadline, actual limited stock). No countdown timers that reset.
- **Accessibility and inclusion:** plain language (aim for roughly grade 7-9 reading level for general audiences); link and button text meaningful out of context; no meaning carried by colour or emoji alone; emoji sparingly, never replacing words; no ALL CAPS sentences; CamelCase hashtags; alt text for images that carry meaning; gender-neutral defaults ("you", "they", "team" not "guys"); avoid ableist idioms ("crazy", "lame", "tone-deaf"); avoid idioms and wordplay if the copy will be translated.
- **Variants:** change one element per variant, label the hypothesis (for example "B: loss framing vs A: gain framing"). Hand test design and sample size to `plan-campaign`.

**Claims substantiation**
| Claim type | Needs |
|---|---|
| Numbers, stats, "X% faster" | Source, date, method; same conditions as stated |
| Comparative ("better than", "cheaper than") | Current, documented, like-for-like comparison |
| Superlative ("best", "#1", "leading") | Independent evidence, or rephrase as opinion/puffery the brief permits |
| "Free", price, savings | Full terms visible; no hidden conditions; taxes/fees per local rules |
| Testimonials, reviews, endorsements | Real, permissioned, unedited in meaning, typical results or disclosed; material connections (paid, gifted, employee) disclosed |
| Health, financial, environmental, safety | Specialist/legal review; specific, not absolute ("carbon neutral" needs proof and scope) |

## Safety gates
- Never invent stats, quotes, testimonials, reviews, customer names/logos, awards, prices or deadlines. Use `[VERIFY]` placeholders.
- Ask before publishing, scheduling, or sending anything (email, SMS, push, social, ads); deliver drafts only unless explicitly told to send.
- Email/SMS: assume consent rules apply - include sender identity, physical address where required, and a working unsubscribe/opt-out.
- Regulated sectors (finance, health, gambling, alcohol, children, political): flag mandatory disclaimers and that compliance review is required before use.
- Do not imitate a competitor's or real person's identity, or imply endorsements that don't exist.

## Output
```
Assumptions: audience <...>, stage <...>, goal <...>, channel <...>, voice <...>
Framework: <name> - why

<COPY, labelled by element: headline / subhead / body / CTA / preheader etc., with character counts for limited fields>

Variants for testing:
A (control): <...>
B: <...>  Hypothesis: <one variable and why>
C: <...>  Hypothesis: <...>

Claims to verify: [VERIFY] <claim> - <evidence needed>
Compliance notes: <disclaimers, approvals>
```

---

**Phase:** Create. **Next:** when this step is done, offer Test: `/anthropica:test-performance`, `/anthropica:test-validation`. Phase map: `/anthropica:workflow`.
