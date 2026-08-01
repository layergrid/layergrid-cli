#!/usr/bin/env python3
from __future__ import annotations

import os
import subprocess
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont

ROOT = Path(__file__).resolve().parents[1]
MEDIA = ROOT / "docs" / "media"


def main() -> None:
    MEDIA.mkdir(parents=True, exist_ok=True)
    render_terminal(
        title="LayerGrid scan",
        text=(ROOT / "testdata" / "golden" / "human-output.txt").read_text(),
        out=MEDIA / "first-scan.png",
        columns=120,
    )
    render_terminal(
        title="LayerGrid",
        text=banner_text(),
        out=MEDIA / "banner.png",
        columns=100,
    )
    print(MEDIA / "first-scan.png")
    print(MEDIA / "banner.png")


def banner_text() -> str:
    env = os.environ.copy()
    env["NO_COLOR"] = "1"
    env["COLUMNS"] = "100"
    env["GOCACHE"] = str(ROOT / ".gocache")
    result = subprocess.run(
        ["go", "run", "./cmd/layergrid"],
        cwd=ROOT,
        env=env,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.DEVNULL,
        check=True,
    )
    return result.stdout


def render_terminal(title: str, text: str, out: Path, columns: int) -> None:
    font, title_font = load_fonts()
    lines = text.rstrip().splitlines()
    char_width = max(12, int(font.getlength("M")))
    line_height = 32
    chrome_top = 86
    padding_x = 34
    width = max(1180, padding_x * 2 + columns * char_width)
    height = max(720, chrome_top + 46 + len(lines) * line_height)

    image = Image.new("RGB", (width, height), "#080d18")
    draw = ImageDraw.Draw(image)
    draw.rounded_rectangle((22, 22, width - 22, height - 22), radius=18, fill="#101828", outline="#344054", width=2)
    draw.ellipse((54, 52, 74, 72), fill="#ff5f57")
    draw.ellipse((84, 52, 104, 72), fill="#febc2e")
    draw.ellipse((114, 52, 134, 72), fill="#28c840")
    draw.text((156, 48), title, fill="#d0d5dd", font=title_font)

    x, y = padding_x, chrome_top
    for line in lines:
        draw.text((x, y), line, fill=line_color(line), font=font)
        y += line_height
    image.save(out)


def load_fonts():
    candidates = [
        "/System/Library/Fonts/Menlo.ttc",
        "/System/Library/Fonts/SFNSMono.ttf",
        "/System/Library/Fonts/Monaco.ttf",
        "/Library/Fonts/Menlo.ttf",
    ]
    for candidate in candidates:
        try:
            return ImageFont.truetype(candidate, 22), ImageFont.truetype(candidate, 18)
        except OSError:
            pass
    font = ImageFont.load_default()
    return font, font


def line_color(line: str) -> str:
    if line.startswith("╭") or line.startswith("╰") or "▸" in line:
        return "#ff7a45"
    if "CRITICAL" in line:
        return "#ffb4a8"
    if "HIGH" in line:
        return "#ff7a45"
    if "MEDIUM" in line:
        return "#facc15"
    if "Grade          │       A" in line or "100 / 100" in line:
        return "#86efac"
    if line.lstrip().startswith("Next"):
        return "#93c5fd"
    return "#e4e7ec"


if __name__ == "__main__":
    main()
