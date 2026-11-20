# Keyboard and screen reader testing scripts

Manual passes catch what automated scanners cannot. Time-box: about 10 minutes per primary flow for keyboard, 10-15 for a screen reader smoke test.

## Keyboard-only pass
Unplug or ignore the mouse. Start from the browser address bar.

| Key | Expect |
|---|---|
| Tab / Shift+Tab | Moves through interactive elements in visual and logical order; skip link appears first |
| Enter | Activates links and buttons; submits forms |
| Space | Activates buttons, toggles checkboxes; scrolls page when on non-interactive content |
| Arrow keys | Move within composite widgets: radio groups, tabs, menus, listboxes, grids, sliders |
| Esc | Closes dialogs, menus, popovers, tooltips; focus returns to the trigger |
| Home / End | First/last item in lists, menus, sliders |

Check on every step:
1. Focus indicator always visible and clearly distinguishable (3:1 against adjacent colours); not hidden under sticky headers, cookie banners or chat widgets.
2. No focus on hidden or off-screen elements (closed menus, inactive carousel slides, visually hidden drawers).
3. Modal dialogs: focus moves inside on open, Tab cycles within, Esc closes, focus returns to the trigger, background is inert.
4. Custom widgets follow WAI-ARIA Authoring Practices keyboard patterns (one Tab stop per composite widget, arrows inside).
5. After dynamic changes (adding a row, deleting an item, route change in a single-page app) focus lands somewhere sensible, not on `body`.
6. Drag-and-drop, sliders and carousels have keyboard and single-click alternatives.
7. No keyboard trap anywhere, including embedded third-party content (video players, maps, payment iframes).
8. The whole primary task can be completed.

## Screen reader smoke test
Pick the pair your users most likely use; test at least one desktop and, for mobile-first products, one mobile pair.

| Platform | Pair | Start / stop |
|---|---|---|
| Windows | NVDA + Firefox or Chrome | Ctrl+Alt+N / Insert+Q |
| Windows | JAWS + Chrome or Edge | Commercial; use if users require it |
| macOS | VoiceOver + Safari | Cmd+F5 |
| iOS | VoiceOver + Safari | Settings > Accessibility, or triple-click side button shortcut |
| Android | TalkBack + Chrome | Volume up + down held 3 s (if shortcut enabled) |

Essential commands:
| Action | NVDA (browse mode) | VoiceOver macOS (VO = Ctrl+Option) |
|---|---|---|
| Next heading / list headings | H / Insert+F7 (elements list) | VO+Cmd+H / VO+U (rotor) |
| Next landmark | D | Rotor > Landmarks |
| Next form field | F | VO+Cmd+J |
| Next link / button | K / B | VO+Cmd+L / rotor |
| Read from here | Insert+Down | VO+A |
| Interact with form | Enter (focus mode) or Tab | VO+Space |

Mobile: swipe right/left moves to next/previous element, double-tap activates, rotor (two-finger twist on iOS) jumps by headings or landmarks.

Check:
1. Page title announced on load and on route change.
2. Headings list gives an outline of the page; one `h1`; no skipped levels used for styling.
3. Landmarks present: banner, navigation, main, contentinfo; multiple navs have distinct labels.
4. Every control announces name, role and state: "Search, edit text", "Menu, collapsed, button", "Accept terms, checkbox, not checked".
5. Images: informative ones described; decorative ones silent; icon-only buttons have names.
6. Form fields announce label, required state, format hint and error message (via `aria-describedby`), and errors are announced on submit.
7. Dynamic updates (toasts, validation, cart count, search results) announced through live regions without stealing focus; not announced twice.
8. Tables announce headers when moving between cells.
9. Nothing important is only in `title` attributes, placeholders, CSS pseudo-content or visual position.
10. Reading order matches visual order.

## Other quick checks
- **Zoom:** browser zoom 200% and 400% (equals 320 CSS px wide at 1280 px window). No clipped text, no overlapping, no 2D scrolling for text content.
- **Text spacing:** apply line-height 1.5, letter-spacing 0.12em, word-spacing 0.16em, paragraph spacing 2em; nothing truncated.
- **Forced colours / high contrast:** emulate `forced-colors: active`; icons, focus rings, borders and selected states remain visible (avoid relying on background colour or box-shadow only).
- **Reduced motion:** emulate `prefers-reduced-motion: reduce`; large animations and parallax are reduced or removed.
- **Voice control:** visible labels match accessible names so "click Save" works.
- **Orientation:** rotate device; content works in both.

## Recording evidence
For each failure note: AT and browser versions, the exact announcement heard (or missing), the keystrokes used, and the element (selector or `file:line`). Short screen recordings with audio help developers reproduce screen reader issues.
