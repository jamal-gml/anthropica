 █████╗ ███╗   ██╗████████╗██╗  ██╗██████╗  ██████╗ ██████╗ ██╗ ██████╗ █████╗ 
██╔══██╗████╗  ██║╚══██╔══╝██║  ██║██╔══██╗██╔═══██╗██╔══██╗██║██╔════╝██╔══██╗
███████║██╔██╗ ██║   ██║   ███████║██████╔╝██║   ██║██████╔╝██║██║     ███████║
██╔══██║██║╚██╗██║   ██║   ██╔══██║██╔══██╗██║   ██║██╔═══╝ ██║██║     ██╔══██║
██║  ██║██║ ╚████║   ██║   ██║  ██║██║  ██║╚██████╔╝██║     ██║╚██████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝ ╚═════╝╚═╝  ╚═╝
                                                                               
# Anthropica

A reusable Claude plugin and skill library for Claude Code, Claude Desktop, Claude Cowork and Claude.ai. The repository is designed for selective use: install the plugin once, then invoke only the workflow needed.

## What is included

- A Claude Code plugin marketplace (`.claude-plugin/marketplace.json`).
- A versioned plugin with 14 manually invoked Skills and four specialist agents.
- A small optional global `CLAUDE.md` for Claude Code users who want baseline behaviour across projects.
- Individually packaged Skills for upload to Claude.ai.
- No default MCP servers, hooks, background processes or external service connections.

## Install in Claude Code (available across projects)

Run in a terminal:

```bash
/plugin marketplace add jamal-gml/anthropic-claude-optimization
/plugin install claude-optimisation@jamal-claude-toolkit
```

Choose **User** scope when prompted to make it available across your projects. If prompted, run `/reload-plugins`.

To use a workflow, invoke it explicitly, for example:

```text
/claude-optimisation:token-efficiency
/claude-optimisation:project-bootstrap
/claude-optimisation:security-review
/claude-optimisation:codebase-analysis
```

To list skills, use `/skills`. To disable or uninstall the plugin, use `/plugin`.

## Claude Desktop and Claude.ai

- **Claude.ai / Claude Desktop Skills:** upload an individual ZIP from `dist/claude-ai-skills/` through Customize → Skills → Create skill → Upload a skill. Each skill can be enabled or disabled independently.
- **Unified Claude plugin directory:** to make the complete plugin discoverable inside Claude's web/Desktop plugin directory, submit it through `https://claude.ai/directory/manage` or have an organisation owner provision it. A GitHub Claude Code marketplace is not itself the same thing as a Claude directory listing.
- The standalone ZIPs use the portable Agent Skills format; Claude Code plugin-specific invocation metadata is removed from those ZIPs.

## Optional global Claude Code instructions

Copy `global/CLAUDE.md` to `~/.claude/CLAUDE.md` if you want concise baseline instructions in every Claude Code project on that machine. If that file already exists, merge manually; do not overwrite it. The plugin itself does not silently alter user configuration.

## Selective design

All plugin skills set `disable-model-invocation: true`. Claude will not decide to load a workflow automatically; the user invokes it with its slash command. This makes activation deliberate and predictable. For Claude.ai, individual skills can be toggled separately.

## Updating

After pulling a repository update, run `/plugin marketplace update` and then `/reload-plugins` if needed. Bump both the plugin and marketplace version for releases.

## Security and scope

The initial release contains no MCP configuration, hooks or auto-running scripts. Review future additions before enabling them. Skills are instructions, not a security boundary. Never put credentials in the repository or skill files.

## Repository layout

```text
.claude-plugin/marketplace.json
plugins/claude-optimisation/.claude-plugin/plugin.json
plugins/claude-optimisation/skills/*/SKILL.md
plugins/claude-optimisation/agents/*.md
global/CLAUDE.md
scripts/package-skills.py
dist/claude-ai-skills/*.zip
```
