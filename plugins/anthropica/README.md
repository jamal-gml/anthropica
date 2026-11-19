# Anthropica Plugin

A selective, cross-project Claude plugin. Only skill descriptions stay in context; full instructions load when a task matches or when you invoke a skill by its slash command.

Install:

```text
/plugin marketplace add jamal-gml/anthropica
/plugin install anthropica@anthropica
```

## Skills

| Skill | Use it for |
|---|---|
| `/anthropica:token-efficiency` | planning a Claude session, reducing token use, managing context, splitting work, or improving prompt efficiency |
| `/anthropica:project-bootstrap` | starting a project, onboarding Claude to a repository, or creating CLAUDE.md/AGENTS.md |
| `/anthropica:architecture-design` | planning a new application, major feature, migration or architectural change |
| `/anthropica:prompt-engineering` | a task is ambiguous, complex, multi-stage, or needs clear acceptance criteria and bounded scope |
| `/anthropica:backend-development` | Implement backend services, business logic, API routes, validation and error handling |
| `/anthropica:codebase-analysis` | Map relevant architecture and trace a requested behaviour in an unfamiliar repository. Use before broad or cross-module changes |
| `/anthropica:database-engineering` | Review or change SQL schemas, migrations, indexes, transactions and data-access code |
| `/anthropica:devops-ci-cd` | Create or troubleshoot CI pipelines, GitHub Actions, branch protections, release automation and infrastructure workflows |
| `/anthropica:frontend-development` | Implement or modify frontend UI, components, state, accessibility and responsive behaviour using the existing project conventions |
| `/anthropica:performance-investigation` | Investigate measured performance, latency, memory, database or frontend bottlenecks without speculative optimisation |
| `/anthropica:release-readiness` | Prepare a release readiness review covering changes, tests, migrations, security, rollback and deployment approvals |
| `/anthropica:security-review` | Review application changes for authentication, authorisation, input validation, secrets, injection and common security weaknesses |
| `/anthropica:technical-documentation` | Write or update concise technical documentation, setup guides, architecture notes and operational runbooks |
| `/anthropica:testing-validation` | Design and run focused unit, integration, end-to-end and regression checks for a requested change |

## Specialist agents

`architect`, `code-reviewer`, `security-reviewer` and `test-engineer` provide bounded specialist workflows in Claude Code. The first three are limited to read-only tools.

## Safety defaults

- No MCP servers, external services, shell hooks or automatic file mutations are enabled by default.
- Skills load on demand; add `disable-model-invocation: true` to a skill to make it manual-only.
- Review and approval remain required for destructive or externally visible operations.
- The plugin does not claim token savings; measure comparable tasks.
