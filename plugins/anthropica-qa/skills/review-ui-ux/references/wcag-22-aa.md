# WCAG 2.2 Level A and AA checklist

Grouped by what you check, with the success criterion (SC) and how to test. "Auto" = automated tools usually catch it; "Manual" = needs a human. Level in brackets. 4.1.1 Parsing is obsolete in WCAG 2.2 and no longer tested.

## Text alternatives and media
| SC | Check | How |
|---|---|---|
| 1.1.1 Non-text content (A) | Informative images have alt conveying purpose; decorative images `alt=""` or CSS; icon buttons have accessible names; complex charts have text equivalent | Auto finds missing alt; Manual judges quality |
| 1.2.1-1.2.5 Media (A/AA) | Captions for prerecorded and live video, transcripts for audio-only, audio description for video where visuals carry meaning | Manual |

## Structure and semantics
| SC | Check | How |
|---|---|---|
| 1.3.1 Info and relationships (A) | Headings are real headings in logical order; lists are lists; tables have header cells; form fields grouped (fieldset/legend for radio groups); landmarks (header, nav, main, footer) | Auto partial; inspect accessibility tree |
| 1.3.2 Meaningful sequence (A) | DOM order matches visual reading order | Manual: read with CSS off or screen reader |
| 1.3.3 Sensory characteristics (A) | Instructions don't rely only on shape, position or sound ("click the round button on the right") | Manual |
| 1.3.4 Orientation (AA) | Works in portrait and landscape unless essential | Manual |
| 1.3.5 Identify input purpose (AA) | Personal-data fields have correct `autocomplete` (name, email, tel, street-address, cc-number...) | Manual / partial auto |
| 2.4.1 Bypass blocks (A) | Skip link or landmarks to jump past repeated navigation | Manual |
| 2.4.2 Page titled (A) | Unique, descriptive `<title>`; SPA route changes update it | Auto + Manual |
| 2.4.4 Link purpose (A) | Link text (plus context) says where it goes; no bare "click here" lists | Manual |
| 2.4.6 Headings and labels (AA) | Headings and labels describe topic or purpose | Manual |
| 3.1.1 / 3.1.2 Language (A/AA) | `lang` on `<html>`; `lang` on passages in another language | Auto + Manual |

## Colour, contrast and visual presentation
| SC | Check | How |
|---|---|---|
| 1.4.1 Use of colour (A) | Errors, links in text, chart series, statuses also use text, icon, underline or pattern | Manual |
| 1.4.3 Contrast minimum (AA) | Text 4.5:1; large text (24px regular or 18.66px bold) 3:1; includes placeholder, disabled-looking but active text, text over images | Auto for plain backgrounds; Manual over images/gradients |
| 1.4.11 Non-text contrast (AA) | UI component boundaries (input borders, toggles, checkboxes), focus indicators, meaningful icons and graph elements 3:1 | Manual with contrast picker |
| 1.4.4 Resize text (AA) | Browser zoom 200%: no loss of content or function | Manual |
| 1.4.5 Images of text (AA) | Real text instead of text baked into images (logos exempt) | Manual |
| 1.4.10 Reflow (AA) | At 320 CSS px width (1280px at 400% zoom) content reflows to one column without horizontal scrolling, except tables, maps, diagrams | Manual |
| 1.4.12 Text spacing (AA) | Line height 1.5, paragraph spacing 2x, letter 0.12em, word 0.16em applied by user: nothing clipped or overlapping | Manual with a text-spacing bookmarklet |
| 1.4.13 Content on hover or focus (AA) | Tooltips/popovers are dismissible (Esc), hoverable (pointer can move onto them) and persistent until dismissed | Manual |

## Keyboard and focus
| SC | Check | How |
|---|---|---|
| 2.1.1 Keyboard (A) | Every action works with keyboard alone, including custom widgets, drag-and-drop, menus, date pickers | Manual |
| 2.1.2 No keyboard trap (A) | Focus can always leave a component (modals trap focus intentionally but close on Esc and return focus) | Manual |
| 2.1.4 Character key shortcuts (A) | Single-key shortcuts can be turned off or remapped, or only work on focus | Manual |
| 2.4.3 Focus order (A) | Tab order follows meaning; no positive `tabindex`; focus moves into opened dialogs and back to the trigger on close | Manual |
| 2.4.7 Focus visible (AA) | Visible focus indicator on every focusable element | Manual |
| 2.4.11 Focus not obscured, minimum (AA, new in 2.2) | Focused element is not entirely hidden by sticky headers/footers, cookie banners, chat launchers | Manual: tab through with banners open |
| 3.2.1 On focus / 3.2.2 On input (A) | Focus or changing a value does not trigger navigation, submission or major context change without warning | Manual |

## Pointer and touch
| SC | Check | How |
|---|---|---|
| 2.5.1 Pointer gestures (A) | Multi-point or path gestures (pinch, swipe) have a single-pointer alternative | Manual |
| 2.5.2 Pointer cancellation (A) | Actions fire on up-event, or can be aborted by moving away | Manual |
| 2.5.3 Label in name (A) | Accessible name contains the visible label text (voice control users say what they see) | Auto partial |
| 2.5.4 Motion actuation (A) | Shake/tilt actions have a UI alternative and can be disabled | Manual |
| 2.5.7 Dragging movements (AA, new) | Every drag (sortable lists, sliders, map pan, kanban) has a click/tap alternative | Manual |
| 2.5.8 Target size minimum (AA, new) | Targets at least 24x24 CSS px, or a 24px circle centred on each target does not overlap another target; exceptions: inline text links, equivalent larger control elsewhere, user-agent default, essential | Manual / partial auto |

## Time, motion and seizures
| SC | Check | How |
|---|---|---|
| 2.2.1 Timing adjustable (A) | Session timeouts warn and allow extension (at least 20 s to respond), unless essential | Manual |
| 2.2.2 Pause, stop, hide (A) | Carousels, auto-updating feeds and animations over 5 s can be paused | Manual |
| 2.3.1 Three flashes (A) | Nothing flashes more than 3 times per second | Manual |
| Best practice (2.3.3 is AAA) | `prefers-reduced-motion` reduces parallax, zoom and large slide transitions | Emulate in dev tools |

## Forms, errors and authentication
| SC | Check | How |
|---|---|---|
| 3.3.1 Error identification (A) | Error described in text, field identified; not colour only | Manual |
| 3.3.2 Labels or instructions (A) | Visible label for every input; format hints (date, password rules) shown before entry | Manual |
| 3.3.3 Error suggestion (AA) | Message says how to fix ("Enter a date like 31/12/2026") | Manual |
| 3.3.4 Error prevention (AA) | Legal, financial or data-changing submissions are reversible, checked, or confirmed | Manual |
| 3.3.7 Redundant entry (A, new) | Information already entered in the same process is auto-filled or selectable | Manual |
| 3.3.8 Accessible authentication, minimum (AA, new) | No memorising or transcribing puzzles to log in unless an alternative exists; paste and password managers allowed; object-recognition CAPTCHA is an allowed exception | Manual |
| 3.2.6 Consistent help (A, new) | Help mechanisms (contact, chat, FAQ link) in the same relative place on every page | Manual |
| 3.2.3 / 3.2.4 Consistent navigation and identification (AA) | Repeated navigation in same order; same function has same name and icon | Manual |

## Assistive technology compatibility
| SC | Check | How |
|---|---|---|
| 4.1.2 Name, role, value (A) | Custom widgets expose role, name, state (expanded, selected, checked, pressed) and update it; follow WAI-ARIA Authoring Practices patterns | Auto partial + screen reader |
| 4.1.3 Status messages (AA) | Success, error counts, "added to cart", search result counts announced via live region without moving focus | Screen reader |

## Reporting a failure
Quote the SC number and name, the measured value (e.g. "contrast 3.2:1, needs 4.5:1"), the element (selector or `file:line`), the user impact, and the fix.
