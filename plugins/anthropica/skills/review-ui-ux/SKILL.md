---
name: review-ui-ux
description: Review screens, flows, mockups or built UI for usability, WCAG 2.2 AA accessibility, UI states, microcopy, responsive layout and consistency, with severity-rated findings. Use when critiquing a design, checking a page before launch, or auditing accessibility (axe, Lighthouse, keyboard, screen reader). Not for building components (use `build-design-system`) or code review (use `review-pull-request`).
---

# UI/UX Review

Find what stops real users, including disabled users, from completing their task, and say exactly how to fix it.

## Procedure
1. **Frame it.** Confirm the artefact (mockup, prototype, URL, branch or code), target users, the 1-3 primary tasks, devices, locales, the design system and any brand or legal constraints. Ask if unknown; do not review without a primary task.
2. **Walk the primary task** as a first-time user, end to end. Record steps, decisions, dead ends, unclear labels and anything you had to remember from a previous screen.
3. **Heuristic pass** with Nielsen's 10: system status, match with the real world, user control and freedom, consistency and standards, error prevention, recognition over recall, flexibility and efficiency, minimalist design, error recovery, help. Read `references/heuristics-and-content.md` for concrete checks per heuristic and microcopy rules.
4. **States pass.** For every view and component: empty/first use, loading, partial data, error, success, offline, permission denied, 0/1/many/huge counts, long text and long translations, disabled. Missing states are findings.
5. **Accessibility pass** (built UI): automated scan, then keyboard-only pass, then 200%/400% zoom and reflow at 320 CSS px, then a screen reader smoke test, then contrast and motion. Read `references/wcag-22-aa.md` for the full criterion checklist and `references/assistive-tech-testing.md` for keyboard and screen reader scripts. For mockups, check contrast, target sizes, focus design, labels, heading structure and colour-only meaning.
6. **Responsive pass** at 320, 375/390, 768, 1024, 1280 and 1440+ px widths, portrait and landscape: no horizontal scroll (except data tables, maps), no clipped or overlapping content, touch targets usable, primary action reachable, tables and forms still workable.
7. **Consistency pass.** Components, tokens, spacing, type scale, icon meaning and terminology match the design system and the rest of the product. One-off styles and look-alike components are findings.
8. **Write findings** with severity, evidence and a concrete fix. Merge duplicates; group systemic issues (one root cause, many places) into one finding with a list of locations.

## Rules

WCAG 2.2 AA thresholds most often missed:
| Check | Requirement (SC) |
|---|---|
| Text contrast | 4.5:1; large text (24px, or 18.66px bold) 3:1 (1.4.3) |
| Non-text contrast | Icons, input borders, focus indicators, chart lines 3:1 against adjacent colours (1.4.11) |
| Target size | At least 24x24 CSS px or enough spacing; exceptions inline links, equivalent control, essential (2.5.8). Recommend 44x44 for touch |
| Focus visible | Every focusable element shows a visible focus indicator (2.4.7) |
| Focus not obscured | Focused element not fully hidden by sticky headers, banners, chat widgets (2.4.11) |
| Keyboard | All functionality operable by keyboard, no traps, logical order (2.1.1, 2.1.2, 2.4.3) |
| Name, role, value | Custom controls expose name, role, state to assistive tech (4.1.2); visible label text is part of the accessible name (2.5.3) |
| Labels and errors | Inputs have visible labels/instructions (3.3.2); errors identified in text (3.3.1) with a suggestion (3.3.3); status messages announced (4.1.3) |
| Reflow and zoom | Usable at 320 CSS px width without 2D scrolling (1.4.10); text resizable to 200% (1.4.4); text-spacing overrides don't break layout (1.4.12) |
| Dragging | Any drag action has a single-pointer alternative, e.g. buttons or menu (2.5.7) |
| Authentication | No cognitive test to log in without an alternative; allow paste and password managers (3.3.8) |
| Redundant entry / consistent help | Don't ask for the same info twice in a flow (3.3.7); help in the same place across pages (3.2.6) |
| Motion | Auto-moving content over 5 s can be paused (2.2.2); nothing flashes more than 3 times a second (2.3.1); honour `prefers-reduced-motion` for parallax and large transitions (best practice, 2.3.3 AAA) |
| Colour | Colour is never the only way to convey meaning (1.4.1) |

Severity:
| Severity | Meaning |
|---|---|
| Critical | Blocks the primary task for some users (keyboard trap, unlabeled required control, data loss, WCAG A/AA failure on a core flow) |
| Major | Causes errors, significant friction or confusion; WCAG AA failure off the core flow |
| Minor | Slows users or looks inconsistent; workaround obvious |
| Polish | Visual refinement, wording improvement |

Gotchas:
- Automated tools find only part (roughly a third to a half) of WCAG issues. A clean axe run is not "accessible"; always do the manual passes and say which were done.
- Placeholder text is not a label; it disappears and usually fails contrast.
- `outline: none` without a replacement focus style is an automatic finding.
- `div`/`span` with click handlers: not focusable, no role, no Enter/Space handling. Recommend native `button`/`a` first.
- Toasts that vanish in under ~5 s and carry errors or actions fail users who read slowly or use screen readers.
- Disabled submit buttons hide why the form cannot be sent; prefer enabled buttons with validation messages.
- Check dark mode and high-contrast/forced-colors mode separately; contrast pairs differ.
- Infinite scroll without a reachable footer or "load more" breaks keyboard and screen reader users.

Tools (built UI):
- `npx @axe-core/cli <url> --tags wcag2a,wcag2aa,wcag21a,wcag21aa,wcag22aa`
- `npx lighthouse <url> --only-categories=accessibility --output=json --output-path=./a11y.json`
- `npx pa11y <url>` for page-level CI checks; `@axe-core/playwright` inside E2E tests for authenticated flows.
- Browser accessibility tree and contrast checker in dev tools; forced-colors emulation; `prefers-reduced-motion` emulation.

## Safety gates
- Review against staging or local builds. Do not submit forms, place orders, send messages or change settings on production without explicit approval.
- Run automated scanners only on pages you are allowed to test; ask before scanning third-party or production sites at scale.
- Redact personal data in screenshots and recordings attached to findings.

## Output
```
# UI/UX review: <artefact>   Scope: <screens/flows>   Users and task: <...>
Passes done: heuristics | states | axe | Lighthouse | keyboard | zoom/reflow | screen reader (<AT + browser>) | responsive (<widths>)
Summary: <2-3 sentences: can users complete the task? top 3 issues>

| # | Severity | Where (screen, element, file:line) | Issue | Evidence (SC / heuristic, measured value) | Recommendation |
|---|---|---|---|---|---|
| 1 | Critical | Checkout > Pay button | Not reachable by Tab | 2.1.1 Keyboard; div with onClick | Use <button type="submit"> |

Systemic issues: <root cause -> locations>
Not reviewed / limits: <...>
What works well (keep): <...>
```

---

**Phase:** Review. **Next:** when this step is done, offer Secure: `/anthropica:secure-code-review`, `/anthropica:secure-repo-governance`, `/anthropica:secure-supply-chain`. Phase map: `/anthropica:workflow`.
