#!/usr/bin/env python3
"""Build portable Claude.ai-compatible ZIPs from the plugin Skills."""
from pathlib import Path
from zipfile import ZipFile, ZIP_DEFLATED
import re

ROOT = Path(__file__).resolve().parents[1]
SKILLS = ROOT / "plugins" / "claude-optimisation" / "skills"
OUT = ROOT / "dist" / "claude-ai-skills"
OUT.mkdir(parents=True, exist_ok=True)
for skill_dir in sorted(p for p in SKILLS.iterdir() if p.is_dir()):
    skill_file = skill_dir / "SKILL.md"
    if not skill_file.is_file():
        continue
    content = skill_file.read_text(encoding="utf-8")
    # disable-model-invocation is a Claude Code extension, not portable Agent Skills metadata.
    content = re.sub(r"^disable-model-invocation:\s*true\s*\n", "", content, count=1, flags=re.M)
    archive = OUT / f"{skill_dir.name}.zip"
    with ZipFile(archive, "w", ZIP_DEFLATED) as zf:
        zf.writestr(f"{skill_dir.name}/SKILL.md", content)
    print(f"Created {archive.relative_to(ROOT)}")
