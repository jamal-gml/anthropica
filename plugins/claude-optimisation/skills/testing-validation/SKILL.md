---
name: testing-validation
description: Design and run focused unit, integration, end-to-end and regression checks for a requested change.
disable-model-invocation: true
---

# Testing Validation

# Testing and Validation

1. Identify the behaviour changed and its acceptance criteria.
2. Find existing tests and conventions.
3. Add the smallest useful regression test.
4. Run the narrowest relevant test first.
5. Expand to integration/build checks when risk warrants it.
6. Capture concise failure excerpts and investigate root causes.
7. Inspect whether tests genuinely exercise the requested behaviour.

Never delete, weaken or bypass tests just to obtain a passing result.
Report exact commands and actual outcomes.
