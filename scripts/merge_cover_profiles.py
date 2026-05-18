#!/usr/bin/env python3
"""Merge Go coverage profiles (mode: atomic) from multiple files into one."""
from __future__ import annotations

import sys
from typing import Dict, Tuple


def parse_line(line: str) -> Tuple[str, int, int]:
    """Return (block_key, num_stmts, count) for a coverage data line."""
    parts = line.rsplit(" ", 2)
    if len(parts) != 3:
        raise ValueError(f"unexpected coverage line: {line!r}")
    key, stmts_s, count_s = parts
    return key, int(stmts_s), int(count_s)


def main() -> None:
    if len(sys.argv) < 3:
        print("usage: merge_cover_profiles.py <out> <in1> [in2 ...]", file=sys.stderr)
        raise SystemExit(2)

    out_path = sys.argv[1]
    in_paths = sys.argv[2:]

    mode: str | None = None
    merged: Dict[str, Tuple[int, int]] = {}

    for path in in_paths:
        with open(path, encoding="utf-8") as handle:
            for raw in handle:
                line = raw.strip()
                if not line:
                    continue
                if line.startswith("mode:"):
                    m = line.split(":", 1)[1].strip()
                    if mode is None:
                        mode = m
                    elif mode != m:
                        raise SystemExit(f"mismatched mode {mode!r} vs {m!r}")
                    continue

                key, stmts, count = parse_line(line)
                if key in merged:
                    old_stmts, old_count = merged[key]
                    if old_stmts != stmts:
                        raise SystemExit(f"statement mismatch for {key}: {old_stmts} vs {stmts}")
                    merged[key] = (stmts, old_count + count)
                else:
                    merged[key] = (stmts, count)

    if mode is None:
        raise SystemExit("no mode line found")

    with open(out_path, "w", encoding="utf-8") as out:
        out.write(f"mode: {mode}\n")
        for key, (stmts, count) in merged.items():
            out.write(f"{key} {stmts} {count}\n")


if __name__ == "__main__":
    main()
