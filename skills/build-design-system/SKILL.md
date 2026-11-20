---
name: build-design-system
description: Create, extend or audit a design system - tokens (tiers, W3C format, theming, dark mode), accessible component APIs, stories per state, visual regression, versioning, deprecation and design-to-code handoff. Use when defining tokens, adding a component or variant, theming, or finding design/code drift. Not for reviewing finished screens (use `review-ui-ux`) or feature UI (use `build-frontend`).
---

# Design System

Build tokens and components that are consistent, accessible by default and safe to change.

## Procedure
1. **Inventory first.** List existing tokens, components, variants, consumers and the design library. Grep for drift: hard-coded colours (`#[0-9a-fA-F]{3,8}\b`, `rgba?\(`, `hsla?\(`), raw `px`/`rem` spacing and font sizes, `z-index:` literals, duplicate look-alike components. Reuse or extend before adding.
2. **Tokens.** Three tiers; components consume semantic (or component) tokens only. Store as W3C Design Tokens (DTCG) JSON as the single source, transform to CSS custom properties and platform outputs with a build tool. Theme by remapping the semantic tier per theme and mode; primitives stay fixed. Read `references/tokens.md` when defining, naming, theming or exporting tokens.
3. **Component API.** Design props before code: variants, sizes, states, slots/composition, controlled vs uncontrolled, ref forwarding, prop pass-through. Read `references/components.md` for API rules, accessibility patterns per widget and the state matrix.
4. **Accessibility built in.** Native element first; WAI-ARIA Authoring Practices pattern for composite widgets (keyboard map, roles, states); visible focus from a token; 24x24 CSS px minimum targets (44x44 recommended for touch); labels wired by `id`; reduced-motion handled inside the component; contrast verified for every semantic text/background pair in every theme.
5. **Document and test.** One story per variant x size x state, plus edge content (long text, RTL, missing icon); interaction tests for keyboard behaviour; automated a11y checks on every story; visual regression on every story in each theme. Read `references/docs-testing-release.md` for story matrix, visual-test setup, versioning and deprecation.
6. **Handoff parity.** Design-tool variable collections map to token tiers and modes map to themes; component and variant names match code exactly. Token changes flow design export -> PR to the tokens source -> CI build -> package release; never hand-edit generated outputs.
7. **Release.** Changelog entry per change, semver from the Rules table, migration notes and codemods for breaking changes, deprecation before removal.

## Rules

Token tiers:
| Tier | Example | Who uses it |
|---|---|---|
| Primitive | `color.blue.600`, `space.4`, `font.size.300` | Only semantic tokens |
| Semantic | `color.text.default`, `color.bg.danger.subtle`, `color.border.focus`, `space.inset.md` | Components and app code |
| Component (only when needed) | `button.primary.bg.hover` | That one component |

- Name by purpose, not value: `color.text.danger`, never `color.red-text` or `blue-button`.
- Order: category.property.variant.state (`color.bg.accent.hover`). Same vocabulary across all tokens.
- Dark mode is a remapping, not an inversion: raise surface lightness for elevation, reduce saturation of large fills, re-check every contrast pair.
- CSS custom properties cannot be used inside media queries; emit breakpoints as build-time constants (Sass/JS) too.
- Every `color.text.*` on its intended `color.bg.*` meets 4.5:1 (3:1 for large text, borders, icons, focus) in every theme; test this in CI.

Component API:
- `variant` for mutually exclusive styles (`primary | secondary | ghost | danger`), `size` for scale (`sm | md | lg`); booleans only for independent toggles. No `isPrimary`/`isDanger` boolean sets.
- More than about 3 optional content props or layouts that vary -> composition (compound components or slots) instead of more props.
- Support controlled and uncontrolled: `value`/`defaultValue`/`onChange`, `open`/`defaultOpen`/`onOpenChange`. Never switch modes during the component's lifetime.
- Forward refs, spread remaining props onto the root element, merge `className`/`style`, keep native attribute names (`disabled`, `type`, `name`).
- Same prop names for the same concept across every component.

Semver:
| Change | Bump |
|---|---|
| Remove or rename a component, prop, variant, token or CSS variable; change a default; change DOM structure or keyboard behaviour consumers rely on | MAJOR |
| New component, prop, variant or token; deprecation warning added | MINOR |
| Bug fix, token value tweak inside the same intent, docs | PATCH (announce visible value changes anyway) |

Gotchas:
- A semantic token that aliases another semantic token across themes creates cycles and confusing overrides; alias semantic -> primitive.
- Component tokens for everything multiplies maintenance; add them only when a component needs independent theming.
- `aria-*` sprinkled on a `div` is not a button; prefer native elements and APG patterns.
- Visual regression flakes on fonts, animations and anti-aliasing; render in one pinned browser/OS image with animations off.
- Generated token files and design exports drift if both sides are edited; one direction of truth only.

## Safety gates
- Ask before publishing a package version, removing or renaming tokens, props or components, or changing token values used across multiple products.
- Ask before publishing a design-library update or pushing token changes back to the design tool.
- Ask before bulk-accepting visual regression baselines; review each changed snapshot.
- Never commit design-tool access tokens or API keys; read them from the environment or a secret store.

## Output
```
Change: <tokens/components added or changed>   Version bump: <MAJOR|MINOR|PATCH + reason>
Tokens: | Name | Tier | Value / alias (light | dark) | Replaces |
Components: <name> - props/variants/states, composition, controlled/uncontrolled, keyboard map
Accessibility: contrast pairs checked (<n> pass / <n> fail), keyboard pattern, focus, target size, AT check
Design <-> code mapping: | Design name | Code name | Status (match / renamed / missing) |
Stories and tests: <stories added>, a11y check result, visual diffs reviewed (<n>)
Drift report: hard-coded values <count + top files>, duplicate components, missing in code/design
Migration notes / deprecations: <old -> new, removal version, codemod>
```
