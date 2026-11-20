# Usability heuristics, UI states and microcopy

## Nielsen's 10 heuristics - concrete checks
| Heuristic | Look for |
|---|---|
| 1. Visibility of system status | Feedback within ~100 ms for clicks; spinner or skeleton after ~1 s; progress with estimate for operations over ~10 s; saved/unsaved indicator; current step in multi-step flows; current location in navigation |
| 2. Match with the real world | Users' words, not internal or technical terms ("Payment failed", not "Error 402"); familiar icons; logical order (address fields in local order); units and formats for the locale |
| 3. User control and freedom | Undo for destructive or bulk actions (prefer undo over "Are you sure?"); cancel on every long operation and modal; back button works and keeps form data; easy exit from flows |
| 4. Consistency and standards | Same action, same label and position everywhere; platform conventions (links look like links, primary action placement consistent); design-system components used |
| 5. Error prevention | Constraints over validation (date picker, input masks, disabled impossible options); inline validation after the field is left, not on every keystroke; confirm irreversible actions with specific wording ("Delete 3 invoices") |
| 6. Recognition over recall | Options visible; recent items and defaults; no need to remember codes from an earlier screen; labels stay visible when typing |
| 7. Flexibility and efficiency | Keyboard shortcuts for power users (documented, non-conflicting); bulk actions; saved filters; sensible defaults; autocomplete |
| 8. Aesthetic and minimalist design | One clear primary action per view; secondary actions visually quieter; no competing calls to action; content above the fold supports the main task |
| 9. Help users recover from errors | Plain language, says what happened and how to fix it, keeps the user's input, focuses or links to the first error, offers a retry for transient failures |
| 10. Help and documentation | Contextual help where users get stuck (inline hints, tooltips that are keyboard accessible); searchable help; contact route |

## UI state checklist (per view or component)
| State | Must have |
|---|---|
| Empty / first use | Explains what will appear and the next action ("Create your first project"), not a blank area |
| No results | Echo the query, suggest fixes (clear filters, check spelling) |
| Loading | Skeleton matching final layout for content; spinner only for short waits; no layout shift when content arrives; long tasks show progress and can run in background |
| Partial | Some data failed: show what loaded and flag what didn't, with retry |
| Error | Specific cause, recovery action, input preserved; distinguish user error, server error, offline, permission |
| Success | Confirmation of what happened and what next; announced to screen readers; not only a transient toast for important outcomes |
| Offline / slow network | Clear indicator; queued or blocked actions explained; no silent data loss |
| Permission denied | Explains why and how to get access; does not reveal data that should be hidden |
| Disabled | Explains why (tooltip or helper text) if the user might expect it to work |
| Data extremes | 0, 1, many, thousands; very long names, emails, numbers; missing optional fields; images missing or slow |
| Destructive | Confirmation or undo, consequence spelled out |

## Forms
- One column for most forms; labels above fields; group related fields; mark optional fields rather than required when most are required (or mark both clearly).
- Input types and `autocomplete` set for mobile keyboards and autofill (`type=email`, `inputmode=numeric`).
- Validate on blur and on submit; show all errors on submit with a summary linked to fields; never clear entered data on error.
- Don't split single values (phone, date) into multiple boxes unless essential; accept common formats and normalise.
- Primary button label states the outcome ("Create account"), placed consistently.

## Microcopy rules
- Buttons: verb + object ("Save changes", "Send invite"), not "OK", "Submit", "Yes".
- Confirmations: question names the action and object; buttons repeat the verb ("Delete project" / "Keep project"), not "Yes/No".
- Errors: what happened + how to fix, no blame, no codes alone, no "Oops". "Card was declined. Try another card or contact your bank."
- Empty states: why it is empty + one action.
- Use sentence case, consistent terminology (pick "sign in" or "log in", not both), plain language at about reading age 12-14, front-load key words.
- Avoid ambiguous link text ("click here", "learn more" repeated) unless context makes it unique to screen readers.
- Numbers, dates, currency formatted for the locale; relative times ("2 min ago") with absolute time available.
- Leave room for translation: text can grow 30-40% (more for short labels); avoid text in images and string concatenation that breaks grammar in other languages; check RTL mirroring.

## Visual hierarchy quick checks
- Squint test: the primary action and the page purpose stand out.
- Spacing follows the scale; related items closer than unrelated (proximity).
- Type scale limited (typically 4-6 sizes); line length 45-80 characters for body text; body text at least 16px on mobile.
- Alignment to a grid; consistent icon style and size; no more than one emphasis technique at once.
- Motion has purpose (feedback, orientation), lasts about 100-300 ms for UI transitions, and respects reduced-motion settings.
