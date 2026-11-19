# anthropica

Provider-neutral command-line tool for preparing projects for AI coding assistants and working with the [Anthropica](https://github.com/jamal-gml/anthropica) role-based skill packs.

## Install

```bash
npm install -g anthropica@latest
anthropica version
```

Or run without installing:

```bash
npx anthropica@latest doctor
```

npm installs a small launcher plus one prebuilt native binary for your platform (Windows, macOS or Linux on x64 or arm64). There are no install scripts.

## Usage

```bash
anthropica doctor                                  # check Git, GitHub CLI and token setup
anthropica analyse --path .                        # metadata scan of a project
anthropica init --path . --assistant all           # create AGENTS.md / CLAUDE.md (never overwrites)
anthropica skills list --repo <anthropica clone>   # list skills or a role pack (--pack)
anthropica github inspect owner/repo               # read-only repository overview
```

Full documentation: [cli/README.md](https://github.com/jamal-gml/anthropica/blob/main/cli/README.md).

## Troubleshooting

- **"platform package ... is not installed"**: optional dependencies were skipped. Reinstall without `--omit=optional` or `--no-optional`.
- **Unsupported platform**: build from source with Go; see the repository.

Each release is built and published by GitHub Actions with npm provenance; the version matches the [GitHub release](https://github.com/jamal-gml/anthropica/releases) of the same number.

License: MIT
