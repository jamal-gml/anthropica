# Stories, testing, versioning and deprecation

## Story matrix (Storybook or equivalent)
For each component:
- One story per variant, one per size, one per state (hover/focus/active can be forced with pseudo-state tooling or interaction steps).
- Edge content: very long label, no icon, icon only, wrapped text, RTL (`dir="rtl"`), translated strings 40% longer, empty data, many items.
- Each theme and mode (light, dark, high contrast, each brand) via a global toolbar parameter, and a visual test per theme for key stories.
- A "playground" story with controls (args) bound to props; keep other stories fixed so visual tests are stable.
- Docs page (MDX or autodocs) with usage guidance, do/don't, props table from types, accessibility notes.

Story naming mirrors the design library: `Components/Button/Primary`, `Components/Button/States/Disabled`.

## Tests per component
| Layer | What | Tool examples |
|---|---|---|
| Unit / interaction | Props behaviour, controlled/uncontrolled, events, keyboard map from the spec | Testing Library + test runner; story `play` functions |
| Accessibility | axe rules on every story in every theme; fail CI on violations | Storybook a11y addon / test runner with axe, `jest-axe`, `@axe-core/playwright` |
| Visual regression | Pixel or perceptual diff of every story per theme and key viewport | Playwright `toHaveScreenshot`, a hosted visual review service, Loki, BackstopJS |
| Type / API | Public types compile; no accidental export changes | `tsc --noEmit`, API extractor reports |
| Bundle | Size budget per component, tree-shaking works | size-limit or bundler analysis |

Keyboard interaction test example (pseudo):
```ts
await user.tab();                 // focus first tab
expect(tabs[0]).toHaveFocus();
await user.keyboard('{ArrowRight}');
expect(tabs[1]).toHaveAttribute('aria-selected', 'true');
```

## Visual regression setup that does not flake
- Render in one pinned container image (same OS, browser version and fonts) locally and in CI. Do not compare screenshots from different OSes.
- Disable animations and transitions (`animations: 'disabled'` or a global CSS override), freeze time and random data, wait for fonts (`document.fonts.ready`) and images.
- Fixed viewport and device scale factor; hide carets and blinking cursors.
- Small, explicit diff threshold (e.g. `maxDiffPixelRatio` around 0.01) rather than loosening on failure.
- Store baselines per theme and viewport; review every diff in the PR. Updating baselines is a deliberate commit with an explanation, never a blanket "update all".
- Run on every PR that touches tokens or components; token changes may legitimately change many snapshots: review a sample from every component and theme, and call it out in the changelog.

## Versioning
- Semantic versioning for the package(s); see the Rules table in `SKILL.md` for what counts as breaking.
- Changesets or conventional commits generate the changelog; each entry says what changed, why, and migration steps.
- Separate token package and component package only if they release on different cadences; otherwise one package avoids version skew.
- Pre-releases (`-next`, `-beta`) for large changes; consumers test against them before the major lands.
- Status per component: experimental (may change in minors, opt-in import path), stable, deprecated.

## Deprecation policy
1. Mark: `@deprecated Use <Replacement> instead. Removal in vX.0.` in JSDoc/types (editors strike it through); in docs, a deprecated badge and migration snippet.
2. Warn: one development-only console warning per session, never in production builds.
3. Provide: a codemod (jscodeshift/ts-morph or the ecosystem's equivalent) for mechanical renames; tokens keep an alias from old name to new name during the window.
4. Measure: find remaining usages across consumers (code search for import names, token names, CSS variable names) and share the list.
5. Remove: only in the next MAJOR after at least one minor release (or the documented window) with the warning; list the removal in the changelog's breaking section.

## Drift and adoption audit
- Code: hard-coded values (see `tokens.md` commands), local re-implementations of system components (search for similar names: `*Button*`, `*Modal*`), overrides of internal class names, pinned old versions across consumers.
- Design library: detached instances, local styles not bound to variables, components missing in code and vice versa, name mismatches.
- Report: | Item | Design | Code | Status (match / renamed / missing / divergent) | Action |
- Adoption metric: share of UI using system components and tokens per product; track over releases rather than chasing 100%.

## Contribution flow
Proposal (problem, evidence of reuse, alternatives) -> design and API review -> build with stories and tests -> accessibility check -> docs -> release as experimental -> promote to stable after use in two places.
