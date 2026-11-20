# Evaluating prompts

A prompt is untested until it has run against inputs you didn't write it for.

## 1. Define success first
- Write success criteria before iterating: specific, measurable, tied to the use case (for example "extracts all invoice line items with exact amounts on 95% of test docs", "never answers out-of-scope questions").
- Separate hard requirements (format valid, no fabricated facts, safety) from quality dimensions (tone, completeness, concision).

## 2. Build a test set
| Bucket | Share | Examples |
|---|---|---|
| Typical | ~50% | Real, representative inputs (anonymised) |
| Edge | ~30% | Empty/missing fields, very long input, ambiguous request, mixed languages, unusual formats |
| Adversarial / negative | ~20% | Out-of-scope requests, prompt-injection inside data, requests that must be refused or escalated |
- Start with 10-20 cases; grow it with every real failure (regression set).
- Store as data (JSONL/CSV): `id, input, expected (or rubric), tags`.
- Keep a held-out set you don't tune against.

## 3. Grade
| Method | Use for | Notes |
|---|---|---|
| Exact/code check | Classification labels, JSON validity, schema, required fields, regex | Cheapest, most reliable; use whenever possible |
| Reference comparison | Extraction, short answers | Normalise whitespace/case; fuzzy match for strings |
| Rubric + model grader | Tone, helpfulness, completeness | Give the grader a specific rubric and ask for reasoning then a score; spot-check it against human judgments |
| Human review | High-stakes, subjective, final sign-off | Blind A/B between prompt versions |
- Model-graded rubrics: 1-5 scale with a definition per level, or pass/fail per criterion; one criterion per question beats one holistic score.
- Measure variance: run each case 3+ times at the production temperature for non-deterministic outputs.

## 4. Iterate
1. Run the baseline prompt on the full set; record scores per case.
2. Read failures; group them by cause (missing context, ambiguous instruction, format drift, hallucination).
3. Change one thing addressing the largest group.
4. Re-run the whole set (not only the failures) to catch regressions.
5. Stop when criteria are met on the held-out set, or when changes stop moving scores - then the fix is elsewhere (better context, tools, a different model, splitting the task).

## 5. Report
```
Prompt version: <id/date>  Model and settings: <model, temperature, thinking>
Test set: <n cases: typical/edge/adversarial>
Results: <pass rate per criterion vs previous version>
Regressions: <cases that got worse>
Changes made: <one line each, with reason>
Known gaps: <...>
```

## Agent and skill checks
- Skills/CLAUDE.md rules: test triggering (does the description make it load for the intended requests and not for unrelated ones?) and behaviour (does the agent follow the procedure on 3-5 realistic tasks?).
- Coding-agent prompts: the acceptance test is the repository's own test/lint/typecheck commands plus a human diff review.
- Check for over-triggering after adding emphatic rules: a rule that fires where it shouldn't is a failure too.
