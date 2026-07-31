#!/usr/bin/env python3
from pathlib import Path
from PIL import Image, ImageDraw, ImageFont

ROOT = Path(__file__).resolve().parents[1]
TEXT = """Scanning agent stack in /Users/nikhil/dev/my-app...

Discovered
  Agents      2     (crewai × 1, langchain × 1)
  Tools       3     (1 function, 1 MCP, 1 shell, 0 code)
  MCP Servers 1     (1 remote, 1 unverified publishers)
  Datasources 1

Lethal Trifecta Detected  ·  1 paths

  path #1  CRITICAL  score −30    [LG-LETHAL-TRIFECTA-01]
    support-agent  ->  Slack MCP  ->  Secrets store
    Fix: scope Slack MCP to chat:read only  (see internal/mcp/slack.py:41)

────────────────────────────────────────
  Trifecta Score      65 / 100
  Grade                C
  Findings             2  (1 critical, 0 high, 1 medium, 0 low)
  Scan time           4.2s
────────────────────────────────────────

Run  layergrid explain LG-LETHAL-TRIFECTA-01  for details.
"""


def main() -> None:
    out = ROOT / "docs" / "media" / "first-scan.png"
    out.parent.mkdir(parents=True, exist_ok=True)
    width, height = 1480, 1080
    image = Image.new("RGB", (width, height), "#0b1020")
    draw = ImageDraw.Draw(image)
    try:
        font = ImageFont.truetype("/System/Library/Fonts/Monaco.ttf", 28)
        title = ImageFont.truetype("/System/Library/Fonts/Monaco.ttf", 22)
    except OSError:
        font = ImageFont.load_default()
        title = font

    draw.rounded_rectangle((28, 28, width - 28, height - 28), radius=18, fill="#101828", outline="#344054", width=2)
    draw.ellipse((58, 58, 78, 78), fill="#ff5f57")
    draw.ellipse((88, 58, 108, 78), fill="#febc2e")
    draw.ellipse((118, 58, 138, 78), fill="#28c840")
    draw.text((160, 55), "LayerGrid scan", fill="#d0d5dd", font=title)

    x, y = 62, 112
    line_height = 40
    for line in TEXT.splitlines():
        fill = "#e4e7ec"
        if "CRITICAL" in line or "Lethal Trifecta" in line:
            fill = "#ffb4a8"
        elif "Trifecta Score" in line or "Grade" in line:
            fill = "#b7f7d8"
        elif line.startswith("Run "):
            fill = "#93c5fd"
        draw.text((x, y), line, fill=fill, font=font)
        y += line_height
    image.save(out)
    print(out)


if __name__ == "__main__":
    main()
