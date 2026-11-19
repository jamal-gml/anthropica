---
name: token-efficiency
description: Reduce unnecessary Claude context and output while preserving correctness. Use when planning a Claude session, reducing token use, managing context, splitting work, or improving prompt efficiency.
disable-model-invocation: true
---

# Token Efficiency

Optimise for the smallest sufficient context, not the shortest possible answer. Preserve correctness, traceability and required validation.

## Workflow
1. Clarify the requested outcome and acceptance criteria before exploring.
2. Inspect repository instructions, status and only the files relevant to the task.
3. Search by symbol, route, error or feature; read targeted ranges instead of dumping whole directories.
4. Separate discovery from implementation. Summarise verified findings before a large change.
5. Delegate only independent work with bounded inputs and explicit output contracts.
6. Use subagents for isolated research/review when parallelism reduces total context; do not delegate trivial tasks.
7. Keep progress updates short and decision-oriented. Avoid repeating source content already available.
8. Use `/compact` or a fresh session with a concise handover when the active context is no longer useful; do not compact away unresolved requirements.
9. Prefer concise command output, targeted tests and structured summaries.
10. Never skip tests, security checks or necessary file inspection solely to save tokens.

## Output
State the change, key files, validation performed, remaining risks and any next action. Distinguish measured token usage from estimated context savings.
