# Anthropica CLI

A small, provider-neutral command-line tool for preparing projects for AI coding assistants and working with the Anthropica skill library. It generates `AGENTS.md` (read by most assistants) and, optionally, assistant-specific files such as `CLAUDE.md`. It never overwrites existing files and makes no AI provider calls.

## Install

**Download** a binary for Windows, macOS or Linux from the [latest release](https://github.com/jamal-gml/anthropica/releases/latest) and put it on your `PATH`.

**Or build from source** (Go 1.23+):

```powershell
git clone https://github.com/jamal-gml/anthropica.git
cd anthropica\cli
go build -o anthropica.exe .
.\anthropica.exe version
```

## Commands

| Command | What it does |
|---|---|
| `anthropica version` | Print the CLI version (also `--version`). |
| `anthropica doctor` | Check for Git, the GitHub CLI and a GitHub token (the token is never printed). |
| `anthropica analyse --path <dir>` | Metadata scan: root files, detected technologies and AI instruction files. Alias: `analyze`. |
| `anthropica init --path <dir> [--assistant universal\|claude\|all]` | Create `.anthropica/config.json` plus `AGENTS.md` (universal), `CLAUDE.md` (claude) or both (all). Existing files are left unchanged; the directory must already exist. |
| `anthropica skills list [--repo <dir>] [--pack <pack>]` | List skills from an Anthropica checkout. `--repo` can be any directory inside the checkout. `--pack anthropica-devops` lists one role pack. |
| `anthropica github inspect OWNER/REPO` | Read-only view of a repository's metadata and root entries. Accepts `OWNER/REPO` or a `github.com` URL. |

## GitHub access

`github inspect` works without a token for public repositories. For private repositories or higher rate limits, set a token in the environment rather than passing `--token` (which can end up in shell history):

```powershell
$env:GITHUB_TOKEN = (gh auth token)
anthropica github inspect your-org/your-repo
```

`GH_TOKEN` is also accepted. A fine-grained token with read-only **Metadata** and **Contents** access is enough.

## Development

```powershell
cd cli
go vet ./...
go test ./...
go build -ldflags "-X github.com/jamal-gml/anthropica/cli/internal/cli.version=1.2.0" -o anthropica.exe .
```

- `main.go` only calls `cli.Execute()`. Each command lives in its own file in `internal/cli/` and is registered in `root.go`.
- Commands write to `cmd.OutOrStdout()` / `cmd.ErrOrStderr()` so tests can run them in-process (`cli_test.go`). The GitHub API base URL is a variable so tests use a local server; tests never call the real API.
- The version is injected at build time with `-ldflags`; release binaries are built by `.github/workflows/release.yml` from the tag.
- Keep the CLI provider-neutral: assistant-specific behaviour belongs behind `--assistant` options, not in core commands.
