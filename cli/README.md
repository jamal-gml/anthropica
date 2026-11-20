# Anthropica CLI

A small, provider-neutral command-line tool for preparing projects for AI coding assistants and working with the Anthropica skill library. It generates `AGENTS.md` (read by most assistants) and, optionally, assistant-specific files such as `CLAUDE.md`. It never overwrites existing files and makes no AI provider calls.

## Install the CLI

### npm (any OS, easiest)

Requires Node.js 18 or later.

```bat
npm install -g anthropica@latest
anthropica version
```

Or run it without installing: `npx anthropica@latest doctor`. npm installs a small launcher plus the prebuilt binary for your platform; there are no install scripts. Update with the same command; uninstall with `npm uninstall -g anthropica`.

### Windows (CMD)

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

### macOS and Linux

```bash
os=$(uname -s | tr '[:upper:]' '[:lower:]'); arch=$(uname -m); [ "$arch" = x86_64 ] && arch=amd64; [ "$arch" = aarch64 ] && arch=arm64
curl -fLO "https://github.com/jamal-gml/anthropica/releases/latest/download/anthropica_${os}_${arch}.tar.gz"
curl -fLO https://github.com/jamal-gml/anthropica/releases/latest/download/checksums.txt
grep "anthropica_${os}_${arch}.tar.gz" checksums.txt | shasum -a 256 -c -
tar -xzf "anthropica_${os}_${arch}.tar.gz" --strip-components 1 "anthropica_${os}_${arch}/anthropica"
sudo install -m 0755 anthropica /usr/local/bin/anthropica && rm anthropica "anthropica_${os}_${arch}.tar.gz" checksums.txt
anthropica version
```

### With Go (any OS)

```bat
go install github.com/jamal-gml/anthropica/cli@latest
```

The binary goes to `%USERPROFILE%\go\bin` (Windows) or `~/go/bin`; make sure that folder is on your `PATH`. The command is named `cli` when installed this way; rename it to `anthropica`, or use a download above.

### Build from source

Requires Git and Go 1.23 or later.

```bat
git clone https://github.com/jamal-gml/anthropica.git
cd anthropica\cli
go build -o anthropica.exe .
.\anthropica.exe version
```

On macOS and Linux use `cd anthropica/cli`, `go build -o anthropica .` and `./anthropica version`.

## Run the CLI

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
