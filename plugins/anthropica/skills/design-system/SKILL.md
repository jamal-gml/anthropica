---
name: design-system
description: Create, extend or audit a design system (design tokens, component APIs, variants and states, documentation, design-to-code parity). Use when defining tokens or components, adding a variant, documenting usage, or finding inconsistencies between design files and code.
---

# Design System

## Foundations (tokens)
- Define tokens in tiers: **primitive** (raw values: `blue-600`, `space-4`), **semantic** (purpose: `color-text-primary`, `color-bg-danger`), and **component** only when needed (`button-padding-x`).
- Components consume semantic tokens only, never primitives or hard-coded values.
- Cover colour (with light and dark modes), typography scale, spacing scale, radius, elevation, motion and breakpoints.
- Verify that every text and background semantic pair meets WCAG AA contrast in each mode.
- One source of truth (for example a tokens JSON file transformed with Style Dictionary) feeding both design tools and code.

## Components
1. Check whether an existing component or variant already covers the need before adding one.
2. Define the API: props and variants (size, emphasis, tone), states (default, hover, focus, active, disabled, loading, error), content slots, and accessibility behaviour (role, keyboard interaction, ARIA).
3. Keep variants orthogonal and minimal; prefer composition over a prop for every case.
4. Match naming exactly between the design tool and code (component names, variant names, token names).

## Documentation (per component)
Purpose and when **not** to use it; anatomy; variants and states with examples; content guidelines; accessibility notes (keyboard, screen reader); do and don't examples; a code usage snippet; the token list used.

## Auditing for drift
- Search code for hard-coded colours, spacing and font sizes (`#[0-9a-f]{3,6}`, `px` values outside tokens).
- Compare component inventories between the design library and code. List missing, renamed and divergent items.
- Flag one-off components that duplicate system ones.

## Change management
- Version the system (semver). Breaking token or prop changes get a migration note and a deprecation period.
- Announce changes to consuming teams with a changelog.

## Output
Tokens or components added or changed, naming mapping between design and code, accessibility checks, and a drift report where relevant.
