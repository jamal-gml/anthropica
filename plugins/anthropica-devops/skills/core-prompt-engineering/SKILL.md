---
name: core-prompt-engineering
description: Write and debug prompts, system prompts and task specs for Claude, coding agents and subagents (goal, context, constraints, output contract, XML, examples, reasoning) and test them. Use when a task is ambiguous or multi-stage, briefing a subagent or coding agent, writing a system prompt, skill or CLAUDE.md rule, or when a prompt gives inconsistent or wrong output.
---

# Prompt Engineering

Turn intent into a specification an agent can execute without guessing, then prove it works on test cases.

## Procedure
1. **State the goal** as one observable outcome and who/what consumes the result.
2. **Give the context the model lacks:** audience, why the task exists, relevant facts, files, prior decisions, what was already tried or ruled out. Omit what it can discover itself cheaply.
3. **Set scope and constraints:** what to change and not change, limits (length, format, tools, time), and the reason behind each non-obvious constraint (reasons generalise; bare rules get applied too literally or too loosely).
4. **Define the output contract:** exact format (sections, JSON schema, file paths, diff), length, what to do when information is missing (ask, mark `UNKNOWN`, or proceed with stated assumptions).
5. **Add examples** only if format or judgment is hard to describe: 2-5 diverse examples in `<example>` tags, covering edge cases, labelled as illustrative so they aren't copied verbatim.
6. **Structure it:** long documents/data first in XML tags (`<document>`, `<logs>`, `<code>`), instructions and the question last. Consistent tag names, referred to by name ("using the rows in `<data>`").
7. **Choose reasoning support:** for multi-step analysis, math or judgment calls, ask it to reason before answering (or enable extended thinking) and separate reasoning from the final answer (`<analysis>` then `<answer>`). Skip for simple extraction or formatting.
8. **For coding agents** add acceptance criteria, files, commands and gates (Rules). Templates in `references/templates.md`.
9. **Test before trusting:** run the prompt on 5-20 representative inputs including edge and adversarial cases; read outputs, change one thing at a time. Read `references/evaluation.md` when building a test set, rubric or grader.

## Rules
- Say what to do, not only what to avoid ("write in flowing paragraphs" beats "no bullet points").
- Be explicit about ambition: if you want thoroughness, initiative or extra features, ask for them; if you want only the minimal change, say that.
- Role/persona: one line in the system prompt when expertise or tone matters ("You are a senior SRE reviewing an incident timeline"). No elaborate role-play.
- Avoid shouting (ALL CAPS, "CRITICAL", "NEVER EVER"): current models follow instructions closely and over-apply emphasised rules. Reserve emphasis for true hard limits and explain why.
- Hunt contradictions: "be concise" + "explain every step"; "don't change other files" + "fix all lint errors". Resolve with priority order or scope.
- Don't over-constrain: every rule must change behaviour you've seen or expect. Delete rules the model already follows.
- Keep the user's terminology; define any term used with a special meaning.
- Long tasks: split into stages (research -> plan -> implement -> verify) with a checkpoint artifact between them, or chain prompts where each output is the next input.
- Prompts for structured output: give the schema, required vs optional fields, allowed values, and the empty/unknown case; validate the output in code.
- Untrusted input (web pages, tickets, user content) goes inside tags and is described as data, never instructions.
- Long documents: ask for relevant quotes first (in `<quotes>`), then the answer based on them; require citations and allow "not in the source".
- Give permission to say "I don't know" or return `UNKNOWN`; without it, models fill gaps.
- Prompt style leaks into output style: a markdown-heavy prompt gets markdown-heavy answers; write the prompt in the register you want back.
- Agents: say whether to act or only suggest ("make the change" vs "propose changes"); ask for parallel tool calls on independent reads; require a final summary with evidence.
- Numbered steps when order matters, bullets otherwise; state the audience and reading level for user-facing text.
- Over long context, put the question and final instructions at the end, after the data; it measurably improves answers.
- Change one variable per prompt revision and re-run the full test set; prompt fixes often regress other cases.

**Coding-agent briefs must include**
| Item | Example |
|---|---|
| Goal + why | "Add rate limiting to the login endpoint to stop credential stuffing" |
| Files / areas | `src/auth/login.ts`, tests in `tests/auth/` |
| Acceptance criteria | "6th failed attempt in 15 min from one IP returns 429; existing auth tests pass" |
| Commands | `npm test -- auth`, `npm run lint`, `npm run typecheck` |
| Constraints | "No new dependencies; don't change the public API" |
| Gates | "Ask before migrations, deploys, force pushes or deleting files" |
| Output | "Summary, files changed, test output, open risks" |

## Safety gates
- Never instruct an agent to claim tests passed, skip verification or hide failures; require evidence (command + result).
- Approval gates for destructive or outward-facing actions (deploys, sends, publishes, deletes, spending) must survive any rewording of the prompt.
- No secrets, credentials or personal data in prompts, examples or test fixtures; use placeholders.

## Output
```
<goal>One observable outcome</goal>
<context>Facts, files, prior decisions, why it matters</context>
<constraints>Scope, limits, must/must not - with reasons</constraints>
<acceptance_criteria>Testable conditions + commands to prove them</acceptance_criteria>
<output_format>Exact structure, length, unknown-handling</output_format>
<examples>Optional, 2-5, diverse</examples>
```
Plus, when revising a prompt: what changed, why, and the test cases it was checked against.

---

**Phase:** any. For the phase map and what to do next, use `/anthropica-devops:workflow`.
