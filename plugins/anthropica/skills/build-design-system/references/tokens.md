# Design tokens: format, naming, theming and handoff

## W3C Design Tokens (DTCG) format
Files use the `.tokens.json` (or `.tokens`) extension. Properties prefixed with `$` are reserved; groups are plain nested objects; aliases use curly-brace paths. Excerpt (aliased primitives such as `color.neutral.*` live in the same source):

```json
{
  "color": {
    "blue": {
      "$type": "color",
      "600": { "$value": "#2563eb" },
      "700": { "$value": "#1d4ed8" }
    },
    "text": {
      "$type": "color",
      "default": { "$value": "{color.neutral.900}", "$description": "Body text on default surfaces" },
      "link": { "$value": "{color.blue.600}" }
    }
  },
  "space": {
    "$type": "dimension",
    "4": { "$value": "16px" }
  },
  "font": {
    "body": {
      "$type": "typography",
      "$value": { "fontFamily": "{font.family.sans}", "fontSize": "{font.size.300}", "fontWeight": 400, "lineHeight": 1.5 }
    }
  }
}
```

- Types: `color`, `dimension`, `fontFamily`, `fontWeight`, `duration`, `cubicBezier`, `number`, and composites `typography`, `shadow`, `border`, `transition`, `gradient`, `strokeStyle`. Set `$type` on the group to inherit.
- `$description` on every semantic token: when to use it.
- Use `$extensions` (namespaced, reverse-domain key) for tool-specific metadata; keep it out of names.
- Check the transformer's supported spec version; the format evolved (e.g. structured colour and dimension values in newer drafts). Pin the tool version.

## Naming
Pattern: `<category>.<property>.<variant>.<state>` with a fixed vocabulary.

| Category | Properties | Variants | States |
|---|---|---|---|
| color | text, bg, border, icon, focus | default, subtle, strong, inverse, accent, success, warning, danger, info, disabled | hover, active, selected |
| space | inset, stack, inline, gap | xs, sm, md, lg, xl | - |
| radius, border-width, shadow/elevation | - | sm, md, lg, full / 0-5 | - |
| font | family, size, weight, line-height, letter-spacing; or composite `typography.heading.lg` | scale steps | - |
| motion | duration, easing | fast, base, slow; standard, enter, exit | - |
| z | index | dropdown, sticky, overlay, modal, toast | - |

Rules:
- Primitives may use scales (`blue.600`, `space.4` = 4 x base unit). Pick one base (4px or 8px) and a modular type scale; avoid ad-hoc steps.
- Semantic names never contain a colour name or a raw value.
- Do not encode the component in semantic names; component tokens live under the component name.
- Prefer a numbered or t-shirt scale consistently, not a mix.
- Never reuse a removed token name with a new meaning.

## Theming and dark mode
Primitives are the palette; each theme maps semantic tokens to primitives.

```css
:root, [data-theme="light"] {
  --color-text-default: var(--color-neutral-900);
  --color-bg-default: var(--color-neutral-0);
  --color-border-focus: var(--color-blue-600);
}
[data-theme="dark"] {
  --color-text-default: var(--color-neutral-50);
  --color-bg-default: var(--color-neutral-950);
  --color-border-focus: var(--color-blue-400);
}
@media (prefers-color-scheme: dark) {
  :root:not([data-theme]) { /* same as dark */ }
}
```
- Respect the OS preference by default, allow an explicit override, persist the choice, and set `color-scheme: light dark` so native controls and scrollbars match.
- Avoid a flash of the wrong theme: set the theme attribute before first paint (inline script or server-rendered attribute).
- Dark mode rules: elevation = lighter surfaces, not darker shadows; avoid pure black and pure white for large areas; desaturate brand fills; check that semantic statuses remain distinguishable; re-run every contrast pair.
- Multi-brand: brand is another mode axis (brand x colour scheme x density). Keep the axes independent and limit combinations you actually ship and test.
- High contrast: test `forced-colors: active`; use `CanvasText`, `Highlight`, `ButtonText` system colours for focus rings and borders where needed.
- Motion: provide `motion.duration.*` tokens and set them to `0ms` or reduced values under `prefers-reduced-motion: reduce`.

## Transformation pipeline
1. Source: DTCG JSON in the repo (one folder per tier; theme files for semantic mappings).
2. Transform with a token build tool (for example Style Dictionary, which supports DTCG) to:
   - CSS custom properties (per theme selector), plus Sass/JS constants for values needed at build time (breakpoints in media queries cannot use CSS variables).
   - TypeScript types or a JS object for type-safe access.
   - iOS/Android outputs if native apps consume the system.
3. Validate in CI: schema/format, unresolved or circular aliases, naming pattern lint, contrast pairs per theme, no primitive referenced directly by components (lint rule or grep).
4. Publish generated outputs as a versioned package; consumers never copy values.

## Design-tool variables to tokens
| Design tool concept | Token concept |
|---|---|
| Variable collection "Primitives" | Primitive tier |
| Collection "Semantic" / "Theme" with modes Light/Dark | Semantic tier, one theme file per mode |
| Variable alias | `{path.to.token}` alias |
| Variable scopes (fill, stroke, text, gap) | `$type` and usage docs |
| Variable description | `$description` |
| Styles (text, effect) | Composite tokens (`typography`, `shadow`) |

Flow: export variables (tool API, plugin or native DTCG export) -> normalise to the repo's DTCG structure with a script -> open a PR with a readable diff -> CI builds and validates -> release. Keep names identical; if the design tool's path separators differ (`/` vs `.`), convert mechanically, never by hand. Decide which side is the source of truth per tier and document it; reverse sync (code -> design) only through the same scripted path.

## Drift audit commands
```bash
G='*.{css,scss,less,tsx,jsx,vue,svelte}'
rg -n -g "$G" -g '!**/tokens/**' '#[0-9a-fA-F]{3,8}\b|rgba?\(|hsla?\('
rg -n -g "$G" '\b(margin|padding|gap|font-size|border-radius)\s*:\s*[0-9.]+(px|rem)'
rg -n -g "$G" 'z-index\s*:\s*[0-9]+'
```
Report counts per file, the closest existing token for each value, and values that need a new token (rare; challenge each one).
