# Changelog

All notable changes to MarkView are recorded here. Earlier releases are
described on the [GitHub Releases page](https://github.com/sbj-ee/markview/releases).

## [Unreleased]

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
