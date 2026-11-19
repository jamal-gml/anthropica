---
name: prompt-engineering
description: Create precise prompts and task specifications for Claude and coding agents. Use when a task is ambiguous, complex, multi-stage, or needs clear acceptance criteria and bounded scope.
---

# Prompt Engineering

Convert the request into an executable task specification.

## Required structure
- Objective: one clear outcome.
- Context: only relevant project facts and constraints.
- Scope: what to change and what not to change.
- Requirements: observable behaviours, not vague adjectives.
- Acceptance criteria: testable completion conditions.
- Validation: exact checks to run and expected evidence.
- Safety: permissions, destructive actions and approval gates.
- Output: files, report or implementation expected.

## Rules
- Ask one focused clarification only when a missing detail blocks safe progress.
- Preserve user terminology and constraints.
- Separate research, planning and implementation when mixing them would create unnecessary context.
- Do not instruct an agent to claim tests passed without running them.
- Avoid excessive role-play, repeated rules and giant context dumps.
- For large tasks, define phases and handover checkpoints.
