---
name: core-token-efficiency
description: Keep Claude Code sessions lean without losing correctness - what loads always vs on use, targeted search and reads, bounded output, subagent fan-out, /compact, /clear, handover notes. Use when planning a long session, context is filling or auto-compacting, output is noisy, deciding whether to delegate, or trimming CLAUDE.md, skills or MCP setup.
---

# Token Efficiency

Use the smallest context that is still sufficient; never trade correctness, validation or safety for tokens.

## Procedure
1. **Pin the goal first.** Restate the outcome and acceptance criteria in one or two lines before exploring; unclear goals cause the most wasted reads.
2. **Locate before reading.** Glob for paths, Grep with `files_with_matches` or `count` first, then `content` with small `-C` context on the 1-3 hits that matter. Search by symbol, route, error text, config key.
3. **Read ranges, not files.** Use offset/limit around the relevant lines for large files; read whole files only when small or when the whole file is the subject.
4. **Bound command output.** Quiet flags, `--stat`/`--name-only`, `| tail -n 50`, `-q`, filters on test runners (run the failing test, not the suite, while iterating; run the full suite once before declaring done).
5. **Fan out with subagents** only for independent, bounded work (sweeping many files, parallel reviews, research) where you need the conclusion, not the raw output. Give a precise brief and an output contract (see `core-prompt-engineering`).
6. **Checkpoint.** After discovery, write a short findings summary (facts, file:line, decisions) before implementing. It survives compaction and makes `/compact` safe.
7. **Manage the window.** `/context` to see what is using space; `/compact <what to keep>` at natural breakpoints; `/clear` (with a handover note) when switching to an unrelated task.
8. **Reply concisely.** Lead with the result; reference `path:line` instead of pasting code back; no recap of what the user just saw.

Read `references/context-mechanics.md` when deciding what to put in CLAUDE.md, skills or MCP config, tuning a long session, or writing a handover note.

## Rules
| Situation | Do | Avoid |
|---|---|---|
| Find where X is handled | Grep symbol/error -> read 30-80 lines around hit | Reading every file in the folder |
| Large log or JSON | Grep/filter for the error or key; head/tail | Dumping the whole file |
| Git history | `git log --oneline -n 20 -- path`, `git diff --stat` then targeted diff | Unbounded `git log -p` |
| Many files to survey | One subagent returns a summary table | Reading them all in the main thread |
| 2-3 known files | Read them yourself | Spawning a subagent for a trivial lookup |
| Context > ~70% full mid-task | Write findings summary, then `/compact` with focus | Waiting for auto-compact at a bad moment |
| New unrelated task | `/clear` after saving a handover note if needed | Carrying stale context forward |

- Always-loaded text (CLAUDE.md, imported files, skill/agent descriptions, MCP tool definitions) costs tokens on every turn of every session: keep it short and move detail into on-demand files.
- Tool output stays in context until compaction: one careless 5,000-line read taxes every later turn.
- Subagents save main-context tokens but usually cost more total tokens; use them when isolation or parallelism is the point.
- Don't re-read a file you just edited to check it; the edit tool already reports failure.
- Make small changes with targeted edits; rewriting a whole file resends all of it.
- Batch independent tool calls in one turn; sequence only dependent ones.
- Plan before large multi-file changes (plan mode or a written plan); a wrong implementation is the most expensive token sink.
- After ~3 failed attempts at the same fix, stop: summarise what was tried and why it failed, then rethink (or `/clear` with that summary) instead of piling on context.
- Output tokens are slower and pricier than input: no restating the question, no narrating each step, no full-file echoes.
- Skill descriptions load every turn: make them specific triggers, push detail into the body and `references/`.
- Disconnect MCP servers a project doesn't use; each adds tool schemas to every turn. Keep `@` imports in CLAUDE.md to small files.
- Ask one clarifying question early rather than exploring three interpretations of an ambiguous request.
- Run long commands in the background and read the result once, instead of polling output repeatedly.
- Send follow-ups to an existing subagent rather than spawning a new one that must rediscover the same context.
- Images, PDFs and screenshots are expensive: read only the pages or regions needed.
- Prefer exact commands and paths in instructions over prose explanations.

## Safety gates
- Never skip tests, security checks, required reads or reviewing a diff to save tokens. If a budget forces a trade-off, say what was not checked.
- Before `/compact` or `/clear`, make sure unresolved requirements, user constraints and open risks are written down; they are the first thing lost.
- Never paste secrets into handover notes or summaries.

## Output
Status reply:
```
<Result in one line>
Evidence: <path:line>, `<command>` -> <key output line>
Not checked / risks: <...>
Then: <next step or question>
```
- Distinguish measured token numbers (from `/context` or `/cost`) from estimates.
- Handover note (before `/clear` or ending a session): goal, done (with files), decisions and why, open issues, next steps, verify commands - template in `references/context-mechanics.md`.

---

**Phase:** any. For the phase map and what to do next, use `/anthropica-security:workflow`.
