# Claude Optimisation Toolkit Plugin

A selective, cross-project Claude plugin. Skills are explicitly invoked with slash commands; they do not automatically load their full instructions on every task.

## Skills

| Skill | Use it for |
|---|---|
| `/claude-optimisation:token-efficiency` | planning a Claude session, reducing token use, managing context, splitting work, or improving prompt efficiency |
| `/claude-optimisation:project-bootstrap` | starting a project, onboarding Claude to a repository, or creating CLAUDE.md/AGENTS.md |
| `/claude-optimisation:architecture-design` | planning a new application, major feature, migration or architectural change |
| `/claude-optimisation:prompt-engineering` | a task is ambiguous, complex, multi-stage, or needs clear acceptance criteria and bounded scope |
| `/claude-optimisation:backend-development` | Implement backend services, business logic, API routes, validation and error handling |
| `/claude-optimisation:codebase-analysis` | Map relevant architecture and trace a requested behaviour in an unfamiliar repository. Use before broad or cross-module changes |
| `/claude-optimisation:database-engineering` | Review or change SQL schemas, migrations, indexes, transactions and data-access code |
| `/claude-optimisation:devops-ci-cd` | Create or troubleshoot CI pipelines, GitHub Actions, branch protections, release automation and infrastructure workflows |
| `/claude-optimisation:frontend-development` | Implement or modify frontend UI, components, state, accessibility and responsive behaviour using the existing project conventions |
| `/claude-optimisation:performance-investigation` | Investigate measured performance, latency, memory, database or frontend bottlenecks without speculative optimisation |
| `/claude-optimisation:release-readiness` | Prepare a release readiness review covering changes, tests, migrations, security, rollback and deployment approvals |
| `/claude-optimisation:security-review` | Review application changes for authentication, authorisation, input validation, secrets, injection and common security weaknesses |
| `/claude-optimisation:technical-documentation` | Write or update concise technical documentation, setup guides, architecture notes and operational runbooks |
| `/claude-optimisation:testing-validation` | Design and run focused unit, integration, end-to-end and regression checks for a requested change |

## Specialist agents

`architect`, `code-reviewer`, `security-reviewer`, and `test-engineer` provide bounded specialist workflows in Claude Code and supported agent surfaces.

## Safety defaults

- No MCP servers, external services, shell hooks or automatic file mutations are enabled by default.
- Skills are manually invoked (`disable-model-invocation: true`).
- Review and approval remain required for destructive or externally visible operations.
- The plugin does not claim token savings; measure comparable tasks.
