# Component API and built-in accessibility

## Before adding a component
1. Can an existing component, variant or composition cover it? Check the inventory and usage in products.
2. Is it needed in at least two products or features? One-offs stay in the app until they repeat.
3. Write the spec: purpose, when not to use it, anatomy, props, states, keyboard map, ARIA, tokens used, responsive behaviour, content rules.
4. Agree names with design before code; the same component, variant and prop names in both.

## API checklist
| Concern | Rule |
|---|---|
| Visual style | One `variant` enum; one `tone`/`intent` enum only if style and meaning are independent (e.g. `variant="outline" tone="danger"`) |
| Scale | `size="sm|md|lg"`; default `md`; sizes change padding, font and icon size together |
| Booleans | Only independent toggles: `disabled`, `loading`, `fullWidth`, `required`. Never mutually exclusive booleans |
| Content | `children` for the main content; named slots or subcomponents for icons, descriptions, actions |
| Composition | Compound components (`<Tabs><Tabs.List><Tabs.Trigger/></Tabs.List><Tabs.Panel/></Tabs>`) when structure varies; share state through context |
| Polymorphism | `as`/`asChild` only when needed (e.g. Button rendered as a link); keep semantics correct: a link-looking button must still be the right element |
| State | Controlled (`value` + `onChange`) and uncontrolled (`defaultValue`) both supported; warn in development if a component switches mode |
| Events | `onX` naming; callbacks receive the value first, the event second if needed; `onOpenChange(open)` style for toggles |
| DOM access | Forward refs to the main interactive element; spread rest props there; merge `className` and `style` |
| IDs | Generate stable IDs for label/description wiring; accept an override |
| Defaults | Sensible, safe defaults (`type="button"` on buttons to avoid accidental form submit) |
| Styling escape hatch | `className`/`style` allowed on the root; no internal-element overrides except through documented slots or tokens |
| Internationalisation | No hard-coded strings; accept labels as props; support RTL (logical CSS properties `margin-inline-start`) |

Anti-patterns: a `type` prop that changes the whole component (split it), props that only make sense together without grouping, styles controlled by many one-off props (`color`, `padding`, `fontSize`), duplicating native props with new names (`isDisabled` next to `disabled`).

## State matrix (every interactive component)
default, hover, focus-visible, active/pressed, selected/checked, indeterminate (where relevant), disabled, read-only, loading, invalid/error, plus empty and overflow content. Each state has tokens, a story and a visual test. Disabled: decide and document whether it stays focusable (`aria-disabled` keeps it discoverable) or not (`disabled`).

## Accessibility per pattern (WAI-ARIA Authoring Practices)
| Component | Element / roles | Keyboard | Notes |
|---|---|---|---|
| Button | `<button>` | Enter, Space | Toggle: `aria-pressed`; icon-only needs `aria-label` |
| Link | `<a href>` | Enter | Navigation only; never a `<button>` that navigates |
| Text field | `<input>` + `<label for>` | native | Hint and error via `aria-describedby`; `aria-invalid` on error; correct `type` and `autocomplete` |
| Checkbox / switch | `<input type=checkbox>`; switch `role="switch"` | Space | Visible label clickable |
| Radio group | `<fieldset>` + `<legend>` + radios | Arrows move and select; Tab enters/leaves group | One Tab stop |
| Select / combobox | Native `<select>` first; custom: `role="combobox"` + `listbox` | Arrows, Enter, Esc, type-ahead | `aria-expanded`, `aria-activedescendant` or roving focus |
| Menu (actions) | `role="menu"`/`menuitem` | Arrows, Home/End, Esc, type-ahead | Only for app-style action menus; site navigation uses links in `<nav>` |
| Tabs | `tablist`/`tab`/`tabpanel` | Arrows between tabs; Tab to panel | `aria-selected`, `aria-controls` |
| Dialog | `<dialog>` or `role="dialog"` + `aria-modal="true"` | Esc closes; focus trapped | Labelled by title; focus first sensible element; return focus on close; background inert |
| Disclosure / accordion | `<button aria-expanded>` controlling a region | Enter, Space | Heading wraps the button in accordions |
| Tooltip | `role="tooltip"` via `aria-describedby` | Shows on focus and hover; Esc dismisses | Hoverable and persistent (WCAG 1.4.13); never the only place for essential info |
| Toast / alert | `role="status"` (polite) or `role="alert"` (urgent) | - | Live region present in DOM before content changes; time long enough to read; actions also reachable elsewhere |
| Data table | `<table>`, `<th scope>` | native; grid pattern only if cells are interactive | Caption; sortable headers use `aria-sort` |
| Slider | `<input type=range>` or `role="slider"` | Arrows, Page Up/Down, Home/End | `aria-valuetext` for human-readable values; also offer a text input (dragging alternative) |

Cross-cutting:
- Focus ring from `color.border.focus` token, at least 2px, 3:1 against adjacent colours, shown with `:focus-visible`; never removed without replacement.
- Hit area at least 24x24 CSS px (44x44 recommended on touch), achievable with padding or pseudo-element without changing visual size.
- Respect `prefers-reduced-motion` inside the component; motion comes from tokens.
- Components never set `tabindex` greater than 0; roving `tabindex` for composite widgets.
- Labels: required props for accessible names where there is no visible text (TypeScript can enforce `aria-label` or `aria-labelledby` for icon-only variants).
- Test with axe on every story plus keyboard interaction tests; smoke test complex widgets with a screen reader before first release.

## Component documentation page
Purpose and when not to use (point to the alternative); anatomy diagram; variants and sizes; states; content guidelines (label length, casing, truncation); accessibility (keyboard map, announced name/role/state, what consumers must provide); do and don't examples; responsive behaviour; props table generated from types; tokens used; changelog and status (experimental, stable, deprecated).
