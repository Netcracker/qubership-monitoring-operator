#!/usr/bin/env python3
"""Print the symptom catalog or one troubleshooting section."""

import sys

with open(sys.argv[1], encoding="utf-8") as catalog:
    lines = catalog.read().splitlines()

if len(sys.argv) == 2:
    for index, line in enumerate(lines):
        if line == "**Symptoms:**":
            end = lines.index("**Root cause:**", index)
            print("\n".join(lines[slice(index - 2, end)]))
else:
    start = lines.index("### " + sys.argv[2])
    end = start + 1
    while end < len(lines) and not lines[end].startswith(("## ", "### ")):
        end += 1
    print("\n".join(lines[start:end]))
