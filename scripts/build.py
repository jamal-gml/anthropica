#!/usr/bin/env python3
"""Build every role pack from the single source in skills/ and agents/.

Generates (do not edit these by hand):
  plugins/<pack>/                      one Claude Code plugin per pack in packs.json
  .claude-plugin/marketplace.json      the marketplace listing every pack
  dist/claude-ai-skills/<skill>.zip    one portable ZIP per skill for Claude.ai uploads
  README.md                            the pack table between the packs:start/end markers

Options:
  --packs   also write dist/packs/<pack>-claude-ai.zip (a ZIP of that pack's skill ZIPs),
            used by the release workflow; not committed.

Standard library only. Output is deterministic so CI can check it is up to date.
"""
from __future__ import annotations

import io
import json
import re
import shutil
import sys
import zipfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SKILLS = ROOT / "skills"
AGENTS = ROOT / "agents"
PLUGINS = ROOT / "plugins"
DIST = ROOT / "dist"
NAME_RE = re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*$")
ZIP_DATE = (1980, 1, 1, 0, 0, 0)
WORKFLOW = "workflow"  # generated per pack; reserved as a skill name

errors: list[str] = []


def read(path: Path) -> str:
    return path.read_text(encoding="utf-8").replace("\r\n", "\n")


def write(path: Path, text: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(text, encoding="utf-8", newline="\n")


def frontmatter(path: Path) -> dict[str, str]:
    m = re.match(r"^---\n(.*?)\n---\n", read(path), re.S)
    if not m:
        errors.append(f"{path.relative_to(ROOT)}: missing YAML frontmatter")
        return {}
    meta = {}
    for line in m.group(1).splitlines():
        if ":" in line and not line.startswith((" ", "\t")):
            key, value = line.split(":", 1)
            meta[key.strip()] = value.strip()
    return meta


def load_items(folder: Path, kind: str) -> dict[str, dict[str, str]]:
    items = {}
    paths = sorted(folder.glob("*/SKILL.md")) if kind == "skill" else sorted(folder.glob("*.md"))
    for path in paths:
        expected = path.parent.name if kind == "skill" else path.stem
        meta = frontmatter(path)
        rel = path.relative_to(ROOT)
        if meta.get("name") != expected:
            errors.append(f"{rel}: name '{meta.get('name')}' must equal '{expected}'")
        if not NAME_RE.match(expected) or len(expected) > 64:
            errors.append(f"{rel}: name must be lowercase-hyphenated, 64 chars max")
        desc = meta.get("description", "")
        if not desc:
            errors.append(f"{rel}: description is required")
        elif len(desc) > 1024:
            errors.append(f"{rel}: description is {len(desc)} chars (max 1024)")
        elif not re.search(r"\buse (when|before|after|for)\b", desc, re.I):
            errors.append(f"{rel}: description must say when to use it ('Use when ...')")
        items[expected] = {"description": desc, "path": str(path)}
    return items


def resolve(pack: dict, key: str, available: dict, core: list[str]) -> list[str]:
    wanted = pack.get(key, [])
    if wanted == "*":
        return sorted(available)
    names = list(dict.fromkeys(wanted + (core if key == "skills" else [])))
    for name in names:
        if name not in available:
            errors.append(f"packs.json: pack '{pack['name']}' lists unknown {key[:-1]} '{name}'")
    return names


def zip_bytes(files: dict[str, bytes]) -> bytes:
    buf = io.BytesIO()
    with zipfile.ZipFile(buf, "w", zipfile.ZIP_STORED) as zf:
        for name in sorted(files):
            info = zipfile.ZipInfo(name, date_time=ZIP_DATE)
            info.create_system = 3  # Unix; the default depends on the host OS and breaks reproducibility
            info.external_attr = 0o644 << 16
            zf.writestr(info, files[name])
    return buf.getvalue()


def phase_of(skill: str) -> str:
    return skill.split("-", 1)[0]


def by_phase(skills: list[str], phases: list[dict]) -> list[str]:
    rank = {p["name"]: i for i, p in enumerate(phases)}
    return sorted(skills, key=lambda s: (rank[phase_of(s)], s))


def summary(description: str) -> str:
    """The 'what it does' half of a description, without the 'Use when ...' trigger."""
    return re.split(r"\s+Use (?:when|before|after|for)\b", description, maxsplit=1)[0].strip()


def footer(skill: str, pack_skills: list[str], phases: list[dict], prefix: str) -> str:
    """Pack-aware hint appended to every skill: its phase and the skills of the next phase."""
    cmd = lambda s: f"`{prefix}{s}`"
    titles = {p["name"]: p["title"] for p in phases}
    phase = phase_of(skill)
    if phase == "core":
        line = f"**Phase:** any. For the phase map and what to do next, use {cmd(WORKFLOW)}."
    else:
        present = [p["name"] for p in phases if p["name"] != "core" and any(phase_of(s) == p["name"] for s in pack_skills)]
        i = present.index(phase)
        nxt = present[(i + 1) % len(present)]
        nexts = ", ".join(cmd(s) for s in by_phase([s for s in pack_skills if phase_of(s) == nxt], phases))
        when = "when this step is done, offer" if i + 1 < len(present) else "last phase; feed findings back into"
        line = f"**Phase:** {titles[phase]}. **Next:** {when} {titles[nxt]}: {nexts}. Phase map: {cmd(WORKFLOW)}."
    return f"\n---\n\n{line}\n"


def workflow_skill(title: str, pack_skills: list[str], all_skills: dict, phases: list[dict], prefix: str) -> str:
    """Generated per pack: the pack's skills grouped by phase, and how to choose the next one."""
    ordered = by_phase(pack_skills, phases)
    titles = {p["name"]: p["title"] for p in phases}
    present = [p for p in phases if any(phase_of(s) == p["name"] for s in ordered)]
    flow = " > ".join(p["title"] for p in present if p["name"] != "core")
    lines = [
        "---",
        f"name: {WORKFLOW}",
        f"description: Phase map of the {title} skills ({flow}) and advice on which skill to use next. "
        "Use when starting a task, when unsure which skill applies, or after finishing a step to pick the next one.",
        "---",
        "",
        f"# {title} workflow",
        "",
        "Phases run in this order. Skip phases that do not apply; loop back when a later phase finds a problem.",
    ]
    for p in present:
        lines += ["", f"**{p['title']}**: {p['when']}"]
        lines += [f"- `{prefix}{s}`: {summary(all_skills[s]['description'])}" for s in ordered if phase_of(s) == p["name"]]
    lines += [
        "",
        "## Choosing the next skill",
        "1. Place the work in a phase from the request and the repository state: no agreed approach yet is Plan; "
        "uncommitted or unreviewed changes are Build, Test or Review; a merged but unreleased change is Ship; "
        "a production problem is Operate.",
        "2. Recommend one skill and say why in one line. Name the one after it.",
        "3. Use it once the user agrees, or straight away if they already asked you to proceed.",
        "4. If no skill here fits, say so rather than stretching one.",
        "",
    ]
    return "\n".join(lines)


def skill_files(name: str, extra: str, portable: bool) -> dict[str, bytes]:
    """Every file in a skill folder; SKILL.md gets the pack footer appended."""
    files = {}
    for path in sorted((SKILLS / name).rglob("*")):
        if not path.is_file():
            continue
        rel = path.relative_to(SKILLS / name).as_posix()
        if rel == "SKILL.md":
            text = read(path).rstrip("\n") + "\n" + extra
            if portable:
                # disable-model-invocation is a Claude Code extension, not portable Agent Skills metadata.
                text = re.sub(r"^disable-model-invocation:.*\n", "", text, flags=re.M)
            files[rel] = text.encode("utf-8")
        else:
            files[rel] = path.read_bytes()
    return files


def portable_pack(pack_skills: list[str], all_skills: dict, phases: list[dict], title: str) -> dict[str, bytes]:
    """Claude.ai skill ZIPs for one pack, keyed by ZIP name, footers scoped to that pack."""
    zips = {}
    for s in pack_skills:
        files = skill_files(s, footer(s, pack_skills, phases, ""), portable=True)
        zips[s] = zip_bytes({f"{s}/{rel}": data for rel, data in files.items()})
    wf = workflow_skill(title, pack_skills, all_skills, phases, "").encode("utf-8")
    zips[WORKFLOW] = zip_bytes({f"{WORKFLOW}/SKILL.md": wf})
    return zips


def pack_readme(pack: dict, skills: list[str], agents: list[str], all_skills: dict, all_agents: dict, cfg: dict) -> str:
    name = pack["name"]
    mk = cfg["marketplace"]["name"]
    lines = [
        f"# {pack['title']} pack (`{name}`)",
        "",
        "> Generated by `scripts/build.py` from `packs.json`. Do not edit files in this folder; edit `skills/`, `agents/` or `packs.json` and rebuild.",
        "",
        f"For: {pack['role']}. {pack['description']}",
        "",
        "## Install",
        "",
        "```text",
        f"/plugin marketplace add {cfg['repository'].removeprefix('https://github.com/')}",
        f"/plugin install {name}@{mk}",
        "```",
        "",
        "Claude loads a skill automatically when a task matches its description. Invoke one directly with its slash command.",
        f"Skills are named `<phase>-<topic>`. Start with `/{name}:{WORKFLOW}` to see the phases and get a suggested next skill;",
        "each skill also ends with a pointer to the next phase.",
        "",
        "## Skills",
        "",
        "| Phase | Command | Use it for |",
        "|---|---|---|",
        f"| Any | `/{name}:{WORKFLOW}` | Phase map of this pack and advice on which skill to use next. |",
    ]
    titles = {p["name"]: p["title"] for p in cfg["phases"]}
    for s in by_phase(skills, cfg["phases"]):
        lines.append(f"| {titles[phase_of(s)]} | `/{name}:{s}` | {all_skills[s]['description'].replace('|', '/')} |")
    if agents:
        lines += ["", "## Agents", "", "| Agent | Role |", "|---|---|"]
        for a in agents:
            lines.append(f"| `{a}` | {all_agents[a]['description'].replace('|', '/')} |")
    return "\n".join(lines) + "\n"


def readme_table(packs: list[dict], resolved: dict, cfg: dict) -> str:
    mk = cfg["marketplace"]["name"]
    rows = ["| Pack | For | Skills | Agents | Install |", "|---|---|---|---|---|"]
    for p in packs:
        skills, agents = resolved[p["name"]]
        rows.append(
            f"| **{p['title']}** | {p['role']} | {len(skills)} | {len(agents)} | `/plugin install {p['name']}@{mk}` |"
        )
    return "\n".join(rows)


def main() -> int:
    cfg = json.loads(read(ROOT / "packs.json"))
    version = cfg["version"]
    if not re.match(r"^\d+\.\d+\.\d+$", version):
        errors.append(f"packs.json: version '{version}' is not semver")
    skills = load_items(SKILLS, "skill")
    agents = load_items(AGENTS, "agent")
    core = cfg.get("core", [])
    phases = cfg["phases"]
    phase_names = [p["name"] for p in phases]
    if WORKFLOW in skills:
        errors.append(f"skills/{WORKFLOW}: the name is reserved for the generated phase map")
    for s in skills:
        if phase_of(s) not in phase_names:
            errors.append(f"skills/{s}: name must start with a phase ({', '.join(phase_names)})")
    for s in core:
        if phase_of(s) != "core":
            errors.append(f"packs.json: core skill '{s}' must be named core-<topic>")

    resolved = {}
    for pack in cfg["packs"]:
        if not NAME_RE.match(pack["name"]):
            errors.append(f"packs.json: pack name '{pack['name']}' must be lowercase-hyphenated")
        resolved[pack["name"]] = (resolve(pack, "skills", skills, core), resolve(pack, "agents", agents, []))

    role_skills = {s for p in cfg["packs"] if p.get("skills") != "*" for s in resolved[p["name"]][0]}
    for orphan in sorted(set(skills) - role_skills):
        print(f"note: skill '{orphan}' is only in the full pack")

    if errors:
        print("Build failed:\n  " + "\n  ".join(errors), file=sys.stderr)
        return 1

    shutil.rmtree(PLUGINS, ignore_errors=True)
    shutil.rmtree(DIST / "claude-ai-skills", ignore_errors=True)

    entries = []
    for pack in cfg["packs"]:
        name = pack["name"]
        pack_skills, pack_agents = resolved[name]
        base = PLUGINS / name
        manifest = {
            "name": name,
            "displayName": "Anthropica" if name == cfg["marketplace"]["name"] else f"Anthropica {pack['title']}",
            "description": pack["description"],
            "version": version,
            "author": cfg["author"],
            "homepage": cfg["repository"],
            "repository": cfg["repository"],
            "license": cfg["license"],
            "keywords": [pack["title"].lower(), "skills"],
        }
        write(base / ".claude-plugin" / "plugin.json", json.dumps(manifest, indent=2) + "\n")
        prefix = f"/{name}:"
        for s in pack_skills:
            for rel, data in skill_files(s, footer(s, pack_skills, phases, prefix), portable=False).items():
                path = base / "skills" / s / rel
                path.parent.mkdir(parents=True, exist_ok=True)
                path.write_bytes(data)
        title = "Anthropica" if name == cfg["marketplace"]["name"] else f"Anthropica {pack['title']}"
        write(base / "skills" / WORKFLOW / "SKILL.md", workflow_skill(title, pack_skills, skills, phases, prefix))
        for a in pack_agents:
            write(base / "agents" / f"{a}.md", read(AGENTS / f"{a}.md"))
        write(base / "README.md", pack_readme(pack, pack_skills, pack_agents, skills, agents, cfg))
        entries.append({
            "name": name,
            "source": f"./plugins/{name}",
            "description": pack["description"],
            "version": version,
            "category": pack["category"],
        })
        print(f"plugin  {name}: {len(pack_skills)} skills, {len(pack_agents)} agents")

    marketplace = {
        "$schema": "https://json.schemastore.org/claude-code-marketplace.json",
        "name": cfg["marketplace"]["name"],
        "owner": cfg["marketplace"]["owner"],
        "metadata": {"description": cfg["marketplace"]["description"], "version": version},
        "plugins": entries,
    }
    write(ROOT / ".claude-plugin" / "marketplace.json", json.dumps(marketplace, indent=2) + "\n")

    # Single skill ZIPs carry footers for the full library; pack ZIPs carry footers scoped to their pack.
    out = DIST / "claude-ai-skills"
    out.mkdir(parents=True, exist_ok=True)
    skill_zips = portable_pack(sorted(skills), skills, phases, "Anthropica")
    for s, data in skill_zips.items():
        (out / f"{s}.zip").write_bytes(data)
    print(f"dist    {len(skill_zips)} skill ZIPs")

    if "--packs" in sys.argv:
        packs_out = DIST / "packs"
        shutil.rmtree(packs_out, ignore_errors=True)
        packs_out.mkdir(parents=True)
        for pack in cfg["packs"]:
            zips = portable_pack(resolved[pack["name"]][0], skills, phases, f"Anthropica {pack['title']}")
            files = {f"{s}.zip": data for s, data in zips.items()}
            (packs_out / f"{pack['name']}-claude-ai.zip").write_bytes(zip_bytes(files))
        print(f"dist    {len(cfg['packs'])} pack ZIPs in dist/packs")

    readme = ROOT / "README.md"
    text = read(readme)
    new = re.sub(
        r"(<!-- packs:start -->\n).*?(<!-- packs:end -->)",
        lambda m: m.group(1) + readme_table(cfg["packs"], resolved, cfg) + "\n" + m.group(2),
        text,
        flags=re.S,
    )
    if new == text and "<!-- packs:start -->" not in text:
        print("note: README.md has no packs:start/end markers; table not updated")
    write(readme, new)
    return 0


if __name__ == "__main__":
    sys.exit(main())
