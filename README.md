```
 █████╗ ███╗   ██╗████████╗██╗  ██╗██████╗  ██████╗ ██████╗ ██╗ ██████╗ █████╗ 
██╔══██╗████╗  ██║╚══██╔══╝██║  ██║██╔══██╗██╔═══██╗██╔══██╗██║██╔════╝██╔══██╗
███████║██╔██╗ ██║   ██║   ███████║██████╔╝██║   ██║██████╔╝██║██║     ███████║
██╔══██║██║╚██╗██║   ██║   ██╔══██║██╔══██╗██║   ██║██╔═══╝ ██║██║     ██╔══██║
██║  ██║██║ ╚████║   ██║   ██║  ██║██║  ██║╚██████╔╝██║     ██║╚██████╗██║  ██║
╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝ ╚═════╝╚═╝  ╚═╝
```

# Anthropica

Role-based Claude skill packs for Claude Code, Claude Desktop, Claude Cowork and Claude.ai. Install only the pack for your job, or the full library, and Claude follows industry-standard practice for that work: commits and pull requests, GitHub Actions, infrastructure, releases, design reviews, marketing copy, SEO, QA and security.

## Choose your pack

<!-- packs:start -->
| Pack | For | Skills | Agents | Install |
|---|---|---|---|---|
| **Full library** | Everyone who wants every skill and agent | 32 | 4 | `/plugin install anthropica@anthropica` |
| **DevOps** | DevOps and platform engineers | 16 | 3 | `/plugin install anthropica-devops@anthropica` |
| **Infrastructure** | Infrastructure, cloud and site reliability engineers | 16 | 3 | `/plugin install anthropica-infrastructure@anthropica` |
| **Software engineering** | Backend, frontend and full-stack developers | 14 | 3 | `/plugin install anthropica-engineering@anthropica` |
| **Design** | Product, UX and UI designers | 6 | 0 | `/plugin install anthropica-design@anthropica` |
| **Marketing** | Marketing, content and growth teams | 6 | 0 | `/plugin install anthropica-marketing@anthropica` |
| **SEO** | SEO specialists and web content teams | 7 | 0 | `/plugin install anthropica-seo@anthropica` |
| **Quality assurance** | QA engineers and testers | 10 | 2 | `/plugin install anthropica-qa@anthropica` |
| **Security** | Application, cloud and product security teams | 11 | 2 | `/plugin install anthropica-security@anthropica` |
<!-- packs:end -->

Every pack also includes the core skills: `token-efficiency`, `prompt-engineering` and `technical-documentation`. Each pack's README in `plugins/<pack>/` lists its skills and agents.

## Install in Claude Code

Run inside Claude Code:

```text
/plugin marketplace add jamal-gml/anthropica
/plugin install anthropica-devops@anthropica
```

Replace `anthropica-devops` with your pack from the table. Choose **User** scope to use it across all your projects, or **Project** scope to share it with a team through the repository's `.claude/settings.json`. Run `/reload-plugins` if prompted. Install more than one pack if your work spans roles.

Claude loads a skill automatically when a task matches its description; only the one-line descriptions sit in context until then. You can also invoke a skill directly, for example `/anthropica-devops:github-pull-requests`. Use `/skills` to list skills and `/plugin` to disable or uninstall.

### Rolling out to a team

Add the marketplace and the team's pack to the repository's `.claude/settings.json` so everyone who trusts the folder gets the same pack:

```json
{
  "extraKnownMarketplaces": {
    "anthropica": { "source": { "source": "github", "repo": "jamal-gml/anthropica" } }
  },
  "enabledPlugins": { "anthropica-devops@anthropica": true }
}
```

## Claude.ai and Claude Desktop

Download the pack you need from the [latest release](https://github.com/jamal-gml/anthropica/releases/latest) (`<pack>-claude-ai.zip`), unzip it, and upload each skill ZIP through **Customize → Skills → Upload a skill**. Single skills are also in `dist/claude-ai-skills/`. Organisation owners can provision skills for everyone.

## Agents

| Agent | Tools | Use |
|---|---|---|
| `architect` | read-only | Design decisions and phased plans before code is written |
| `code-reviewer` | read-only | Reviewing a diff before commit or merge |
| `security-reviewer` | read-only | Security review of changes touching auth, input, secrets or infrastructure |
| `test-engineer` | full | Writing and running focused tests |

## Optional global instructions

Copy `global/CLAUDE.md` to `~/.claude/CLAUDE.md` for concise baseline behaviour in every Claude Code project. If that file already exists, merge by hand. The plugins never change your configuration.

## Updating

Run `/plugin marketplace update anthropica`, then `/reload-plugins`.

## Security and scope

No MCP servers, hooks, background processes or external connections. Skills are instructions, not a security boundary; they tell Claude to ask before destructive or externally visible actions such as deployments, force pushes or production changes. Never put credentials in skill files.

## Contributing

Each skill has one source file in `skills/<name>/SKILL.md`, and each agent one in `agents/<name>.md`. `packs.json` decides which skills and agents each role pack gets. `plugins/`, `.claude-plugin/marketplace.json`, `dist/claude-ai-skills/` and the pack table above are generated:

```bash
python scripts/build.py
```

CI fails if the generated files are out of date. See [CLAUDE.md](CLAUDE.md) for how to add a skill or a pack and how to cut a release.

## Repository layout

```text
skills/<name>/SKILL.md        source of every skill
agents/<name>.md              source of every agent
packs.json                    role packs, version, marketplace metadata
scripts/build.py              generates everything below
plugins/<pack>/               generated Claude Code plugins
.claude-plugin/marketplace.json
dist/claude-ai-skills/*.zip   generated portable skill ZIPs
global/CLAUDE.md              optional personal baseline
.github/workflows/            CI check and tag-driven release
```
