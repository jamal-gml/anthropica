---
name: ui-ux-review
description: Review screens, flows, mockups or built UI for usability, accessibility (WCAG 2.2 AA), visual consistency and content clarity, with prioritised, actionable findings. Use when critiquing a design, checking a page before launch, or auditing accessibility.
---

# UI/UX Review

## Inputs
Confirm what is being reviewed (mockup, prototype, live URL or code), the target users, their main task, devices, and any design system or brand rules.

## Review lenses
1. **Task flow:** can the primary user complete the main task without guessing? Count steps, dead ends and decisions.
2. **Usability heuristics:** visibility of system status, match with real-world language, user control (undo, cancel), consistency, error prevention, recognition over recall, clear error recovery.
3. **Accessibility (WCAG 2.2 AA):**
   - Text contrast at least 4.5:1 (3:1 for large text and UI components).
   - Everything reachable and operable by keyboard, with a visible focus indicator and logical tab order.
   - Semantic structure: one H1, ordered headings, landmarks, labelled form fields, alt text that conveys purpose.
   - Touch targets at least 24×24 CSS px. No information conveyed by colour alone. Respect reduced-motion settings.
   - Error messages identify the field and explain how to fix it.
4. **Visual design:** hierarchy (one clear primary action per view), spacing rhythm, alignment, typography scale, and consistent use of design-system components and tokens.
5. **Content:** plain language, specific button labels ("Save changes", not "OK"), helpful empty, loading and error states.
6. **Responsive:** check narrow mobile (320–375 px), tablet and desktop. No horizontal scroll or clipped content.

## Evidence
Cite the exact screen, element or `file:line`. Where tools are available, use automated checks (axe, Lighthouse accessibility) and state that automated tools catch only part of WCAG issues.

## Output
A prioritised list:

| # | Severity | Where | Issue | Why it matters | Recommendation |
|---|---|---|---|---|---|

Severity: **Critical** (blocks the task or excludes users), **Major** (causes errors or significant friction), **Minor** (polish). Finish with what works well, so it is kept.
