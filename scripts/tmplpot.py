#!/usr/bin/env python3
"""Extracts the {{T "..."}} and {{Tf "..." ...}} strings of html/template files into a POT file for msgcat."""
import re
import sys

PATTERN = re.compile(r'\bTf?\s+"((?:[^"\\]|\\.)*)"')


def main(out, files):
    """Writes the strings of files, with their file as location, to out."""
    entries = {}
    for path in files:
        for m in PATTERN.finditer(open(path, encoding="utf-8").read()):
            entries.setdefault(m.group(1), set()).add(path)
    with open(out, "w", encoding="utf-8") as f:
        f.write('msgid ""\nmsgstr ""\n"Content-Type: text/plain; charset=UTF-8\\n"\n\n')
        for msgid in sorted(entries):
            f.write("#: " + " ".join(sorted(entries[msgid])) + "\n")
            f.write('msgid "' + msgid + '"\nmsgstr ""\n\n')


if __name__ == "__main__":
    main(sys.argv[1], sys.argv[2:])
