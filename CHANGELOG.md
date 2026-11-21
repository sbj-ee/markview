# Changelog

All notable changes to MarkView are recorded here. Earlier releases are
described on the [GitHub Releases page](https://github.com/sbj-ee/markview/releases).

## [Unreleased]

## [1.2.0] - 2026-11-21

### Fixed

- Paragraphs no longer run off the right side of the preview. A wide code
  block, table or image used to widen the whole page, so all text wrapped
  at that width. Code blocks, math and Mermaid blocks now wrap; tables shrink
  their columns and images scale down to fit; the preview has a right margin
  so text isn't drawn under the scrollbar.

### Changed

- Updates now work like zterminal's and ztodo's:
  - macOS replaces the running MarkView.app in place (mount the verified DMG
    read-only, `ditto` the app alongside, swap; the old bundle is restored if
    the swap fails) instead of opening the DMG for a manual drag.
  - After installing, MarkView offers **Restart Now** and reopens the current
    file, warning first if there are unsaved changes.
  - The signature must be for this release: `SHA256SUMS` must be signed with
    the comment `markview vX.Y.Z SHA256SUMS`, so an older release's signature
    can't be replayed under a newer tag.
  - Downloads come only from the release's exact GitHub URLs and follow
    redirects only to GitHub's release-asset hosts.
  - The installer gets a private copy of the package, re-checked against the
    signed SHA-256 just before install.
  - On Linux, a MarkView not installed from the .deb (`/usr/bin/markview`)
    no longer runs `apt`; it keeps the verified package and says where it is.

## [1.1.0] - 2026-10-29

### Added

- **Help > Check for Updates…** now installs updates, following zwriter's
  flow: it downloads the release's `SHA256SUMS` and `SHA256SUMS.minisig`,
  verifies the minisign signature against the public key built into
  MarkView, downloads the package and checks its SHA-256, then installs the
  .deb with `pkexec apt install` (Linux) or opens the .dmg (macOS).
- **Help > Check for Updates at Startup** toggle (on by default; at most one
  quiet check a day).
- "Skip This Version" now works (the button was never shown before).
- Releases publish a minisign-signed `SHA256SUMS` next to the packages.

### Changed

- Previously "Download" only opened the .dmg link in a browser, with no
  signature check.
