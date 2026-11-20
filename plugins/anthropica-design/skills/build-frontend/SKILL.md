---
name: build-frontend
description: Build or change frontend UI to production standard (components, state ownership, data fetching, forms, accessibility, Core Web Vitals, XSS-safe rendering, token storage, Testing Library tests, i18n). Use when building or changing components, pages, routes, forms, client-side state or data fetching. Not for design tokens (use `build-design-system`) or UX audits (use `review-ui-ux`).
---

# Frontend Development

Deliver UI that matches the codebase's patterns, works for keyboard and screen-reader users, stays fast, and is safe with untrusted data.

## Procedure
1. **Identify the stack** from manifests and existing code: framework and version, rendering mode (SPA, SSR, SSG, server components), router, styling (CSS modules, utility classes, CSS-in-JS), component library, state and data-fetching libraries, form library, i18n library, test tools.
2. **Find the closest existing page or component** and follow its file layout, naming, props style, styling and test placement. Reuse existing components and design tokens before creating new ones.
3. **Pin down the contract**: API endpoint and response shape (from types, OpenAPI or the backend), loading/empty/error/success states, permissions, and which component owns each piece of state. Never invent endpoints; stub only behind a clearly marked mock if asked.
4. **Choose state placement** with the Rules table. Keep server data in the data-fetching cache, not duplicated into global state.
5. **Build the component**: semantic HTML first, then styling, then interaction. Split when a component has more than one reason to change; keep presentational pieces free of fetching.
6. **Implement every state**: loading (skeleton or spinner without layout shift), empty (with next action), error (message plus retry), success, partially loaded, disabled, and long/overflowing content. Read `references/accessibility.md` when building interactive widgets, forms, modals, menus or custom controls.
7. **Forms**: label every field, validate on submit and on blur (not every keystroke), show errors next to fields and in a summary, move focus to the first error, disable double submit, keep user input on failure, and mirror server validation (the server is authoritative).
8. **Performance check**: avoid new heavy dependencies, lazy-load below-the-fold and route-level code, size images, avoid layout shift. Read `references/performance-security.md` when adding dependencies, images, fonts, third-party scripts, or rendering user-supplied content.
9. **Test** with Testing Library-style tests: query by role and accessible name (`getByRole('button', { name: /save/i })`), then label, then text; `getByTestId` last. Use `userEvent` over `fireEvent`; `findBy*` for async; mock the network at the HTTP layer (MSW or equivalent), not the hook. Cover states, keyboard interaction and error paths. E2E only for critical journeys.
10. **Verify**: typecheck, lint (including a11y lint rules), unit tests, and view the page in a browser at mobile and desktop widths with keyboard only. Inspect the diff; do not reformat unrelated files.

## Rules
- State placement:
  | State | Put it in |
  |---|---|
  | Server data (lists, entities) | data-fetching cache (query library, framework loader, server component) |
  | Shareable or bookmarkable UI (filters, tab, page, search) | URL query params |
  | Form input | form library or local component state |
  | Local UI (open, hover, draft) | component state, lifted only as far as needed |
  | Truly app-wide client state (theme, session, feature flags) | context or a small store |
  Derive values instead of storing them; never sync state with effects when it can be computed.
- Data fetching: one cache key per resource and params; set sensible stale times; invalidate or update the cache after mutations; cancel or ignore stale responses (race conditions); handle 401 by re-auth flow, not silent loops; no fetch waterfalls (parallelise or fetch at route level).
- Optimistic updates only with rollback on error and when the action rarely fails.
- Accessibility minimums: native elements before ARIA (`<button>` not clickable `<div>`); every input has a programmatic label; visible focus indicator; logical tab order with no positive `tabindex`; dialogs trap and restore focus and close on Escape; text contrast 4.5:1 (3:1 for large text and UI boundaries); targets at least 24x24 CSS px; never convey meaning by colour alone; announce async results with a live region; respect `prefers-reduced-motion`; images have meaningful `alt` or `alt=""` if decorative.
- Security: never inject unsanitised HTML (`dangerouslySetInnerHTML`, `v-html`, `innerHTML`, `[innerHTML]`); if unavoidable, sanitise with a vetted sanitiser (DOMPurify) at render time. Validate URLs before using them in `href`/`src` (block `javascript:`). Do not store long-lived tokens in `localStorage`; prefer `HttpOnly`, `Secure`, `SameSite` cookies set by the server. Never put secrets in client code or public env vars (`NEXT_PUBLIC_`, `VITE_`, etc. ship to the browser). Client-side checks are UX, not authorisation.
- Responsive: mobile-first, layouts with flex/grid, no fixed heights on text containers, test at 320 px width and 200% zoom.
- i18n: no hard-coded user-facing strings if the project uses an i18n library; use ICU plurals and placeholders, never string concatenation; format dates, numbers and currency with `Intl`/library; allow 30-50% text expansion; support RTL with logical CSS properties (`margin-inline-start`) where RTL locales exist.
- Effects: only for syncing with external systems; always clean up subscriptions, timers and listeners; include all dependencies.
- Keys for lists are stable IDs, never array indexes for reorderable or filterable lists.
- Read `references/framework-notes.md` when working in React/Next.js, Vue/Nuxt, Angular or Svelte and the framework idiom is unclear.

## Safety gates
- Ask before adding a dependency (state its gzipped size and why the existing stack cannot do it), changing global styles/tokens, changing routing or auth flows, or removing a component used elsewhere (search usages first).
- Do not change API contracts from the frontend side; coordinate with `build-backend`.
- Never commit real user data, API keys or tokens in fixtures, mocks or stories.

## Output
```
## Change
<components/pages and behaviour>

## States covered
loading | empty | error | success | disabled | overflow  (note any not applicable)

## Accessibility
<keyboard path, labels, focus management, contrast checks done>

## Performance
<bundle impact, lazy loading, images>

## Tests
<test files, commands run, results; manual browser check at widths>

## Follow-ups / risks
```

---

**Phase:** Build. **Next:** when this step is done, offer Review: `/anthropica-design:review-ui-ux`. Phase map: `/anthropica-design:workflow`.
