"""Writes the translated help pages ../src/help/<lang>/ from locales/<lang>.po.

po4a writes only the pages that are at least 80 % translated; the app shows the English page for
the others. Run with --update to refresh help.pot and the PO files after the English pages changed
(make help-po).
"""

import shutil
import subprocess
import sys
from pathlib import Path

LANGS = ["de", "fr", "it"]

here = Path(__file__).parent


def main() -> None:
    update = "--update" in sys.argv
    for lang in LANGS:
        shutil.rmtree(here.parent / "src" / "help" / lang, ignore_errors=True)
    subprocess.run(["po4a", "--no-backups", "--wrap-po", "no", "--package-name", "standblatt-help",
                    *([] if update else ["--no-update"]), "po4a.cfg"], cwd=here, check=True)


if __name__ == "__main__":
    main()
