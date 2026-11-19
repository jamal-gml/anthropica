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

## Command-line tool

The `anthropica` CLI prepares any project for AI coding assistants and lists skills and packs. It is provider-neutral: `init` writes `AGENTS.md`, plus `CLAUDE.md` with `--assistant claude` or `all`, and never overwrites existing files.

```powershell
anthropica doctor
anthropica analyse --path .
anthropica init --path . --assistant all
anthropica skills list --pack anthropica-devops
anthropica github inspect your-org/your-repo
```

Full command reference: [cli/README.md](cli/README.md).

### Install the CLI

#### npm (any OS, easiest)

Requires Node.js 18 or later.

```bat
npm install -g anthropica@latest
anthropica version
```

Or run it without installing: `npx anthropica@latest doctor`. npm installs a small launcher plus the prebuilt binary for your platform; there are no install scripts. Update with the same command; uninstall with `npm uninstall -g anthropica`.

#### Windows (CMD)

Downloads the latest release, verifies its checksum and adds it to your user `PATH`. Uses `curl`, `tar` and `certutil`, which ship with Windows 10 and 11.

```bat
mkdir "%LOCALAPPDATA%\Programs\anthropica"
cd /d "%LOCALAPPDATA%\Programs\anthropica"
curl -fLO https://github.com/jamal-gml/anthropica/releases/latest/download/anthropica_windows_amd64.zip
curl -fLO https://github.com/jamal-gml/anthropica/releases/latest/download/checksums.txt
certutil -hashfile anthropica_windows_amd64.zip SHA256
findstr anthropica_windows_amd64.zip checksums.txt
```

The two hashes printed must match. Then unpack, add to `PATH`, and open a **new** CMD window:

```bat
%SystemRoot%\System32\tar.exe -xf anthropica_windows_amd64.zip --strip-components 1
del anthropica_windows_amd64.zip checksums.txt
powershell -NoProfile -Command "$p=[Environment]::GetEnvironmentVariable('Path','User'); if (($p -split ';') -notcontains $env:LOCALAPPDATA+'\Programs\anthropica') { [Environment]::SetEnvironmentVariable('Path', $p.TrimEnd(';')+';'+$env:LOCALAPPDATA+'\Programs\anthropica', 'User') }"
```

```bat
anthropica version
anthropica doctor
```

On ARM devices, use `anthropica_windows_arm64.zip`. To update, repeat the download and unpack steps. To uninstall, delete `%LOCALAPPDATA%\Programs\anthropica` and remove it from your user `PATH`.

#### macOS and Linux

```bash
os=$(uname -s | tr '[:upper:]' '[:lower:]'); arch=$(uname -m); [ "$arch" = x86_64 ] && arch=amd64; [ "$arch" = aarch64 ] && arch=arm64
curl -fLO "https://github.com/jamal-gml/anthropica/releases/latest/download/anthropica_${os}_${arch}.tar.gz"
curl -fLO https://github.com/jamal-gml/anthropica/releases/latest/download/checksums.txt
grep "anthropica_${os}_${arch}.tar.gz" checksums.txt | shasum -a 256 -c -
tar -xzf "anthropica_${os}_${arch}.tar.gz" --strip-components 1 "anthropica_${os}_${arch}/anthropica"
sudo install -m 0755 anthropica /usr/local/bin/anthropica && rm anthropica "anthropica_${os}_${arch}.tar.gz" checksums.txt
anthropica version
```

#### With Go (any OS)

```bat
go install github.com/jamal-gml/anthropica/cli@latest
```

The binary goes to `%USERPROFILE%\go\bin` (Windows) or `~/go/bin`; make sure that folder is on your `PATH`. The command is named `cli` when installed this way; rename it to `anthropica`, or use a download above.

#### Build from source

Requires Git and Go 1.23 or later.

```bat
git clone https://github.com/jamal-gml/anthropica.git
cd anthropica\cli
go build -o anthropica.exe .
.\anthropica.exe version
```

On macOS and Linux use `cd anthropica/cli`, `go build -o anthropica .` and `./anthropica version`.

### Run the CLI

```bat
anthropica --help
anthropica doctor
anthropica analyse --path C:\path\to\project
anthropica init --path C:\path\to\project --assistant all
anthropica skills list --repo C:\path\to\anthropica --pack anthropica-devops
set GITHUB_TOKEN=<your token>
anthropica github inspect your-org/your-repo
```

In PowerShell set the token with `$env:GITHUB_TOKEN = (gh auth token)`. `skills list` needs a local clone of this repository; point `--repo` at it or run it from inside it.

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
cli/                          anthropica command-line tool (Go)
plugins/<pack>/               generated Claude Code plugins
.claude-plugin/marketplace.json
dist/claude-ai-skills/*.zip   generated portable skill ZIPs
global/CLAUDE.md              optional personal baseline
.github/workflows/            CI checks and tag-driven release (packs and CLI binaries)
```
