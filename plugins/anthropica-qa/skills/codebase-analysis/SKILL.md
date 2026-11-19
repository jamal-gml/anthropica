---
name: codebase-analysis
description: Map relevant architecture and trace a requested behaviour in an unfamiliar repository. Use before broad or cross-module changes.
---

# Codebase Analysis

## Procedure
1. Confirm repository root, branch and working-tree state.
2. Read project instructions and manifests.
3. Search for the named route, symbol, feature or error.
4. Trace only direct dependencies needed to understand the behaviour.
5. Inspect related tests and configuration.
6. Report verified findings before proposing changes.

## Constraints
- Analysis is read-only unless the user separately requests implementation.
- Do not read the entire repository.
- Do not infer architecture from folder names alone.
- Separate observed facts from assumptions.

## Output
Objective; relevant files; execution/data flow; dependencies; risks; recommended smallest change.
