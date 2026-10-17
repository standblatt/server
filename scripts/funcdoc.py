#!/usr/bin/env python3
"""Fails when a Go or TypeScript/Svelte function has no doc comment directly above it.

Go: every `func` in internal/ and cmd/ (tests excluded) needs a `//` comment line above it.
TS: every function declaration, class method and arrow function assigned to a const in web/src (components/ui,
api.gen.ts, locales and .d.ts excluded) needs a `/** … */` or `//` comment above it.
Usage: scripts/funcdoc.py [paths]; without paths the whole repo is checked.
"""
import glob
import os
import re
import sys

GO_FUNC = re.compile(r"^func ")
TS_FUNC = re.compile(r"^\s*(export\s+)?(async\s+)?function\s+\w+")
TS_ARROW = re.compile(r"^\s*(export\s+)?const\s+\w+\s*=\s*(async\s*)?(\([^)]*\)|\w+)\s*(:\s*[^=]+)?=>")
TS_METHOD = re.compile(r"^\s+(?:(?:private|public|protected|static|async|get|set)\s+)*(?!if\b|for\b|while\b|switch\b|catch\b|return\b|function\b)[A-Za-z_]\w*\s*\([^)]*\)\s*(?::\s*[^{=]+)?\{\s*$")
TS_SKIP = ("/components/ui/", "api.gen.ts", "/locales/", ".d.ts")


def files(paths):
    out = []
    for a in paths:
        if os.path.isdir(a):
            for ext in ("go", "ts", "svelte"):
                out += glob.glob(os.path.join(a, f"**/*.{ext}"), recursive=True)
        else:
            out.append(a)
    return sorted(set(out))


def missing(path):
    lines = open(path, encoding="utf-8").read().split("\n")
    go = path.endswith(".go")
    if go and path.endswith("_test.go"):
        return []
    if not go and (any(s in path for s in TS_SKIP) or "/node_modules/" in path or "/.svelte-kit/" in path):
        return []
    out = []
    for i, line in enumerate(lines):
        if go:
            hit = GO_FUNC.match(line)
        else:
            hit = TS_FUNC.match(line) or TS_ARROW.match(line) or TS_METHOD.match(line)
        if not hit:
            continue
        prev = lines[i - 1].strip() if i > 0 else ""
        if not (prev.startswith("//") or prev.endswith("*/")):
            out.append(f"{path}:{i + 1}: {line.strip()}")
    return out


def main():
    paths = sys.argv[1:] or ["internal", "cmd", "web/src"]
    found = [m for p in files(paths) for m in missing(p)]
    print("\n".join(found))
    print(f"functions without a doc comment: {len(found)}")
    sys.exit(1 if found else 0)


if __name__ == "__main__":
    main()
