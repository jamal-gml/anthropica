# How Claude Code context works

Check current behaviour with `/context` and the official docs; details evolve between versions.

## What is in context, and when
| Content | When it loads | Cost implication |
|---|---|---|
| System prompt and built-in tool definitions | Every turn | Fixed overhead |
| MCP server tool definitions | Every turn (unless the client defers them behind a tool search) | Each connected server adds its schemas; disconnect servers you don't use in this project |
| User, project and parent-directory `CLAUDE.md` (and local variants) | Session start, every turn | Keep each short; every line is paid for constantly |
| Files imported with `@path` from CLAUDE.md | Session start, every turn | Import small files only; link large ones as "read X when Y" instead |
| `CLAUDE.md` in subdirectories | When Claude reads files in that directory | Good place for module-specific rules |
| Auto memory index | Session start | Keep entries one line each |
| Skill and subagent `name` + `description` | Every turn | Keep descriptions tight but specific; they decide triggering |
| Skill `SKILL.md` body | When the skill is invoked | Keep the core procedure here, ~60-150 lines |
| Skill `references/*.md` | Only when read | Put long checklists, templates, examples here |
| Conversation: your messages, tool calls, tool results, replies | Accumulates until compaction | The largest variable cost; tool output dominates |
| Subagent work | In the subagent's own window | Parent receives only the final report |

## Compaction and clearing
- `/compact [instructions]`: replaces the conversation with a summary. Pass what must survive, for example `/compact keep the failing test names, the API contract decision and remaining TODOs`.
- Auto-compaction triggers near the window limit and can happen mid-task; compact earlier at a clean breakpoint instead.
- After compaction, exact file contents and earlier tool output are gone; re-read targeted ranges rather than trusting memory of them.
- `/clear`: empty conversation, same project instructions. Use between unrelated tasks; save a handover note first if work continues later.
- Resuming an earlier session restores its context; start fresh instead if most of it is stale.

## CLAUDE.md hygiene
- Include: build/test/lint commands, repo layout in a few lines, conventions Claude would not guess, hard rules (safety, attribution, branches), pointers to deeper docs.
- Exclude: things discoverable from code in one grep, long style guides, generic advice, changelogs, duplicated README content.
- Target well under ~200 lines per file; prefer "Read `docs/x.md` when Y" pointers.
- Move situational rules into skills (loaded on use) or subdirectory CLAUDE.md files.

## Subagent economics
Use a subagent when:
- The work reads many files or long outputs and you need only the conclusion.
- Tasks are independent and can run in parallel (for example reviewing 5 modules).
- You want an unanchored second opinion (give it the code, not your conclusion).

Don't use one when:
- You already know the 1-3 files; reading them yourself is cheaper.
- The task needs back-and-forth with the user or shared evolving state.
- The brief would be longer than doing the work.

Brief = goal, scope (files/dirs), what is already ruled out, exact output format and length cap. Ask for `path:line` evidence so results are verifiable.

## Output-heavy commands: lean versions
| Instead of | Use |
|---|---|
| `cat big.log` | `grep -n "ERROR" big.log \| tail -n 40` |
| `ls -R` / `tree` | Glob with a pattern, or `tree -L 2` |
| `git log -p` | `git log --oneline -n 20`, then `git show <sha> -- path` |
| `git diff` on a big change | `git diff --stat`, then `git diff -- path` |
| Full test suite while iterating | Single test file/name; full suite once at the end |
| Verbose build | Default or `-q`, then rerun failing step verbosely |
| `npm install` logs | `--silent`/`--no-progress`, check exit code |
| Large JSON API response | `jq` to the needed fields |

## Handover note template
```
# Handover: <task>  (<date>)
Goal: <outcome + acceptance criteria>
Branch / state: <branch, uncommitted changes?>
Done:
- <change> (<files>)
Decisions:
- <decision> - because <reason>
Findings (verified):
- <fact> (<path:line> or command)
Open issues / risks:
- <issue>
Next steps:
1. <step>
Verify with:
- <exact command> -> expected <result>
Constraints from the user: <must/must not>
```
