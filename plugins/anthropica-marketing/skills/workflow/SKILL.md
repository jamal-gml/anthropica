---
name: workflow
description: Phase map of the Anthropica Marketing skills (Plan > Create) and advice on which skill to use next. Use when starting a task, when unsure which skill applies, or after finishing a step to pick the next one.
---

# Anthropica Marketing workflow

Phases run in this order. Skip phases that do not apply; loop back when a later phase finds a problem.

**Plan**: Understand the code, requirements, risks and approach before changing anything.
- `/anthropica-marketing:plan-campaign`: Plan marketing campaigns and launches - objective and KPI tree, segments, message house, channel mix by funnel stage, budget and pacing, timeline, UTMs, measurement, experiments and post-campaign review.

**Create**: Write content and copy.
- `/anthropica-marketing:create-marketing-copy`: Write and edit conversion-focused marketing copy (landing pages, emails, ads, social, announcements, push/SMS) that is on-brand, accurate, accessible and testable.
- `/anthropica-marketing:create-seo-content`: Plan, write and optimise search-led content - intent, SERP analysis, keyword clusters, briefs, E-E-A-T, on-page and schema, refreshing or pruning pages.

**Any phase**: Cross-cutting skills for documentation, prompts and token use.
- `/anthropica-marketing:core-documentation`: Write or update docs verified against the repo (README, setup and how-to guides, tutorials, reference, explanations, runbooks) using Diataxis and docs-as-code checks.
- `/anthropica-marketing:core-prompt-engineering`: Write and debug prompts, system prompts and task specs for Claude, coding agents and subagents (goal, context, constraints, output contract, XML, examples, reasoning) and test them.
- `/anthropica-marketing:core-token-efficiency`: Keep Claude Code sessions lean without losing correctness - what loads always vs on use, targeted search and reads, bounded output, subagent fan-out, /compact, /clear, handover notes.

## Choosing the next skill
1. Place the work in a phase from the request and the repository state: no agreed approach yet is Plan; uncommitted or unreviewed changes are Build, Test or Review; a merged but unreleased change is Ship; a production problem is Operate.
2. Recommend one skill and say why in one line. Name the one after it.
3. Use it once the user agrees, or straight away if they already asked you to proceed.
4. If no skill here fits, say so rather than stretching one.
