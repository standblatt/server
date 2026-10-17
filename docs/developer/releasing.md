# Releasing

1. Every change has an entry under `## [Unreleased]` in `CHANGELOG.md` ([Keep a Changelog](https://keepachangelog.com/en/1.1.0/), entries end with `(by @HB9HIL)`).
2. Pick the version by [SemVer](https://semver.org/): new features raise the minor version, fixes the patch version.
3. Move the entries into a new section `## [X.Y.Z] - YYYY-MM-DD` below an empty `## [Unreleased]`, commit as `chore(release): vX.Y.Z`.
4. Create an annotated tag `vX.Y.Z` and push the commit and the tag.

The tag run of GitHub Actions builds the image with version `vX.Y.Z` (tags `vX.Y.Z` and `latest` on ghcr.io), publishes the Helm chart `X.Y.Z` to `oci://ghcr.io/standblatt/charts` and creates a GitHub release with the changelog section and the print agent binaries. The hosted installation follows `latest` through Keel.

Release tags are signed (`git tag -s`).
