# TODO

## Recently Completed (usability & infrastructure)

- [x] Toolbar tooltips and native menu bar for discoverability
- [x] Command palette (Cmd/Ctrl+Shift+P)
- [x] Reclaim content space when file tree / TOC panes are hidden
- [x] Accept a markdown file as a positional CLI argument
- [x] Transient "Saved" confirmation in the status bar
- [x] Welcome / empty state when no document is open
- [x] Highlight active panel toggle buttons
- [x] Unsaved-changes indicator in the status bar
- [x] Preserve scroll position when reloading the same document
- [x] CI workflow (gofmt, vet, test, build) as a required check on `main`
- [x] Branch protection ruleset on `main`
- [x] Dialog consistency — unified sizes and dismiss-button labels
- [x] Dialog consistency — normalized error-message casing
- [x] First-run onboarding
- [x] Smoother theme/font switching (in-place re-render, no disk reload, scroll preserved)
- [x] Mermaid diagram rendering in HTML/Print/PDF export
- [x] Math/LaTeX rendering in HTML/Print/PDF export
- [x] Release automation — tag-triggered GitHub Release (Linux .deb/tarball, macOS .dmg)
- [x] Improved table editing (per-column alignment, pipe escaping, edit existing tables in place)

## Features

- [ ] Code signing and notarization for macOS distribution
- [ ] Linux .deb package testing and distribution
- [ ] Vim keybindings mode
- [ ] Custom CSS themes support
- [ ] Plugin system for extensions
- [ ] Sync with cloud storage (iCloud, Dropbox, etc.)
- [ ] Collaborative editing support

## Improvements

- [ ] Performance optimization for large files
- [ ] Better image handling (drag & drop, paste from clipboard)
- [ ] Footnote support
- [ ] Bibliography/citation support

## Bug Fixes

- _None open._ (Theme-switch rendering glitch fixed by the in-place re-render above.)

## Documentation

- [ ] Video tutorials
- [ ] User guide with screenshots
- [ ] Contributing guidelines

## Technical Debt

- [ ] Increase test coverage
- [ ] Add integration tests
- [ ] macOS Developer ID signing + notarization in the release workflow (needs Apple secrets; currently ad-hoc signed)
