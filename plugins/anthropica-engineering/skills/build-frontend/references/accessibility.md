# Accessibility implementation notes (WCAG 2.2 AA)

For audits and reporting, use `review-ui-ux`. This file is for building.

## Semantics first
| Need | Use | Avoid |
|---|---|---|
| Action | `<button type="button">` | `<div onClick>`, `<a>` without `href` |
| Navigation | `<a href>` | button that changes location |
| Page regions | `<header> <nav> <main> <aside> <footer>`, one `<main>` | div soup |
| Headings | one `<h1>` per page, no skipped levels | headings chosen for size |
| Lists | `<ul>/<ol>/<li>` | divs with bullets |
| Data tables | `<table>`, `<th scope>`, `<caption>` | grid of divs for tabular data |
| Toggle | `<button aria-pressed>` or `<input type="checkbox" role="switch">` | div with class "on" |
| Form grouping | `<fieldset><legend>` | visual boxes only |

First rule of ARIA: if a native element does it, use the native element. ARIA changes what is announced, not behaviour; you must add keyboard support yourself.

## Labels and names
- Inputs: `<label for>` or wrapping label. Placeholder is not a label.
- Icon-only buttons: `aria-label="Close dialog"` or visually hidden text.
- Extra hints and errors: `aria-describedby` pointing at the hint/error element.
- Required: `required` attribute (plus visual indicator explained once).
- Invalid: `aria-invalid="true"` on the field while in error.
- Links: text describes the destination; avoid "click here"; same text means same destination.

## Keyboard
- Everything operable with Tab, Shift+Tab, Enter, Space, Escape, arrow keys where the pattern expects it.
- Never remove focus outlines without a visible replacement (`:focus-visible`), at least 2 px, with 3:1 contrast.
- No positive `tabindex`. Use `tabindex="-1"` only to make an element programmatically focusable (headings after navigation, error summary).
- Focus must not be hidden behind sticky headers or overlays (WCAG 2.4.11).
- Provide a "Skip to main content" link on pages with repeated navigation.

## Common widget patterns
Follow the WAI-ARIA Authoring Practices pattern for each; prefer a tested library component when the project has one.
| Widget | Essentials |
|---|---|
| Modal dialog | prefer native `<dialog>` with `showModal()`; else `role="dialog" aria-modal="true" aria-labelledby`; move focus in on open, trap Tab, Escape closes, restore focus to the trigger on close, inert background |
| Menu button / dropdown | `aria-haspopup`, `aria-expanded`; arrows move between items; Escape closes and returns focus |
| Tabs | `role="tablist"/"tab"/"tabpanel"`, `aria-selected`, `aria-controls`; arrow keys switch tabs; one tab in tab order |
| Combobox / autocomplete | `role="combobox"`, `aria-expanded`, `aria-controls`, `aria-activedescendant`; announce result count |
| Accordion | button inside heading, `aria-expanded`, `aria-controls` |
| Tooltip | appears on hover and focus, dismissable with Escape, does not hold essential info |
| Toast / status | `role="status"` (polite) or `role="alert"` (assertive, errors only); do not auto-dismiss critical messages quickly |
| Carousel | pause control, no auto-advance without it, slides labelled |

## Dynamic content
- Route changes in SPAs: move focus to the new page's `<h1>` (or a wrapper) and update `document.title`.
- Async results ("12 results", "Saved"): announce via a live region that exists in the DOM before the update.
- Loading: `aria-busy="true"` on the updating region; avoid announcing every spinner.
- Do not move focus unexpectedly on input or hover.

## Forms
- Validate on submit; on error, render an error summary at the top with links to fields, focus it, and mark fields `aria-invalid` with `aria-describedby` to the message.
- Error text says what to do ("Enter a date after 1 Jan 2024"), not just "Invalid".
- Use correct `type` and `autocomplete` (`email`, `tel`, `current-password`, `one-time-code`) for autofill and mobile keyboards.
- Do not block paste in password or code fields; avoid cognitive tests for login (WCAG 3.3.8).
- Do not ask for the same information twice in one flow (3.3.7).

## Visual
- Contrast: body text 4.5:1, large text (24 px or 18.66 px bold) 3:1, UI component boundaries and focus indicators 3:1.
- Do not use colour alone for status; add icon or text.
- Content reflows at 320 px wide without horizontal scrolling (except tables, maps, code).
- Text resizes to 200% without loss; text spacing overrides do not break layout.
- Target size at least 24x24 CSS px or enough spacing (2.5.8); dragging actions have a single-pointer alternative (2.5.7).
- Motion: honour `prefers-reduced-motion`; nothing flashes more than 3 times per second.

## Media
- Informative images: concise `alt` describing purpose; decorative: `alt=""`; complex charts: summary text or data table.
- Video: captions; audio: transcript; no autoplay with sound.
- SVG icons: `aria-hidden="true"` when decorative; `role="img"` + `<title>` when meaningful.

## Testing
- Lint: `eslint-plugin-jsx-a11y` / `vuejs-accessibility` / Angular template a11y rules.
- Automated: axe-core via `jest-axe`/`vitest-axe` or Playwright `@axe-core/playwright` (catches ~30-40% of issues).
- Manual: keyboard-only pass, screen reader smoke test (NVDA or JAWS on Windows, VoiceOver on macOS/iOS, TalkBack on Android), 200% zoom, 320 px width.
- Testing Library role queries double as accessibility checks: if `getByRole` cannot find it, assistive tech likely cannot either.
