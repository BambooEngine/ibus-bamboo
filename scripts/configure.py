#!/usr/bin/env python3
"""Generate versioned files from their *.template counterparts.

Usage: configure.py [VERSION]

If VERSION is given, it is written to data/bamboo-version first.
Every @VERSION@ in the templates is replaced by the content of data/bamboo-version.
"""
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
VERSION_FILE = ROOT / "data" / "bamboo-version"
TEMPLATE_DIRS = ["data", "build"]
PLACEHOLDER = "@VERSION@"


def main(argv):
    if len(argv) > 2:
        sys.exit(__doc__)
    if len(argv) == 2:
        VERSION_FILE.write_text(argv[1].strip() + "\n", encoding="utf-8")

    version = VERSION_FILE.read_text(encoding="utf-8").strip()
    for directory in TEMPLATE_DIRS:
        for template in sorted((ROOT / directory).rglob("*.template")):
            output = template.with_suffix("")
            content = template.read_text(encoding="utf-8")
            output.write_text(content.replace(PLACEHOLDER, version), encoding="utf-8")
            print(f"Generated {output.relative_to(ROOT)}")


if __name__ == "__main__":
    main(sys.argv)
