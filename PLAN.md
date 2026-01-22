# Packaging Plan for MarkView

This document outlines the steps to package MarkView as a `.deb` package for Debian/Ubuntu and a `.dmg` disk image for macOS.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Linux .deb Package](#linux-deb-package)
3. [macOS .dmg Package](#macos-dmg-package)
4. [Automation with Makefile](#automation-with-makefile)
5. [CI/CD Integration](#cicd-integration)

---

## Prerequisites

### Common Requirements

- Go 1.21+
- Git
- Make

### For .deb Packaging

- `dpkg-deb` (included in Debian/Ubuntu)
- `fakeroot` (optional, for building without root)
- `lintian` (optional, for package validation)

```bash
sudo apt install dpkg-dev fakeroot lintian
```

### For .dmg Packaging

- macOS with Xcode Command Line Tools
- `create-dmg` (optional, for prettier DMGs)

```bash
# Install Xcode tools
xcode-select --install

# Install create-dmg (optional)
brew install create-dmg
```

### For Cross-Platform Building

- Docker
- fyne-cross

```bash
go install github.com/fyne-io/fyne-cross@latest
```

---

## Linux .deb Package

### Directory Structure

Create the following structure for the .deb package:

```
markview-1.0.0/
├── DEBIAN/
│   ├── control          # Package metadata
│   ├── postinst         # Post-installation script (optional)
│   └── prerm            # Pre-removal script (optional)
├── usr/
│   ├── bin/
│   │   └── markview     # Binary executable
│   ├── share/
│   │   ├── applications/
│   │   │   └── markview.desktop    # Desktop entry
│   │   ├── icons/
│   │   │   └── hicolor/
│   │   │       ├── 64x64/apps/
│   │   │       │   └── markview.png
│   │   │       ├── 128x128/apps/
│   │   │       │   └── markview.png
│   │   │       └── 256x256/apps/
│   │   │           └── markview.png
│   │   ├── doc/
│   │   │   └── markview/
│   │   │       ├── README.md
│   │   │       └── copyright
│   │   └── metainfo/
│   │       └── com.sbj-ee.markview.metainfo.xml  # AppStream metadata
│   └── lib/
│       └── markview/    # Additional resources (if needed)
```

### Step 1: Create DEBIAN/control

```
Package: markview
Version: 1.0.0
Section: editors
Priority: optional
Architecture: amd64
Depends: libc6 (>= 2.17), libgl1, libx11-6, libxcursor1, libxrandr2, libxinerama1, libxi6, libxxf86vm1
Maintainer: Your Name <your.email@example.com>
Homepage: https://github.com/sbj-ee/markview
Description: A powerful markdown viewer and editor
 MarkView is a fast, cross-platform markdown viewer and editor
 with syntax highlighting, live reload, multiple themes, and
 advanced features like full-text search, backlinks, and templates.
```

### Step 2: Create Desktop Entry

Create `usr/share/applications/markview.desktop`:

```ini
[Desktop Entry]
Name=MarkView
Comment=Markdown Viewer and Editor
Exec=markview %F
Icon=markview
Terminal=false
Type=Application
Categories=Office;TextEditor;Utility;
MimeType=text/markdown;text/x-markdown;
Keywords=markdown;editor;viewer;notes;
StartupNotify=true
```

### Step 3: Create AppStream Metadata

Create `usr/share/metainfo/com.sbj-ee.markview.metainfo.xml`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<component type="desktop-application">
  <id>com.sbj-ee.markview</id>
  <name>MarkView</name>
  <summary>A powerful markdown viewer and editor</summary>
  <metadata_license>MIT</metadata_license>
  <project_license>MIT</project_license>
  <description>
    <p>
      MarkView is a fast, cross-platform markdown viewer and editor with
      syntax highlighting, live reload, and multiple themes.
    </p>
    <p>Features include:</p>
    <ul>
      <li>Full Markdown and GFM support</li>
      <li>Syntax highlighting for 250+ languages</li>
      <li>8 color themes</li>
      <li>Full-text search across files</li>
      <li>Document templates</li>
      <li>Export to HTML, PDF, DOCX</li>
    </ul>
  </description>
  <launchable type="desktop-id">markview.desktop</launchable>
  <url type="homepage">https://github.com/sbj-ee/markview</url>
  <url type="bugtracker">https://github.com/sbj-ee/markview/issues</url>
  <screenshots>
    <screenshot type="default">
      <caption>MarkView main window</caption>
      <image>https://raw.githubusercontent.com/sbj-ee/markview/main/screenshots/main.png</image>
    </screenshot>
  </screenshots>
  <releases>
    <release version="1.0.0" date="2026-01-22">
      <description>
        <p>Initial release with full feature set.</p>
      </description>
    </release>
  </releases>
  <content_rating type="oars-1.1"/>
</component>
```

### Step 4: Build Script

Create `scripts/build-deb.sh`:

```bash
#!/bin/bash
set -e

VERSION="1.0.0"
ARCH="amd64"
PKG_NAME="markview"
PKG_DIR="${PKG_NAME}_${VERSION}_${ARCH}"

echo "Building MarkView .deb package v${VERSION}..."

# Clean previous builds
rm -rf "dist/${PKG_DIR}" "dist/${PKG_NAME}_${VERSION}_${ARCH}.deb"
mkdir -p "dist/${PKG_DIR}"

# Build the binary
echo "Compiling binary..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o "dist/${PKG_DIR}/usr/bin/markview" ./cmd/markview

# Create directory structure
mkdir -p "dist/${PKG_DIR}/DEBIAN"
mkdir -p "dist/${PKG_DIR}/usr/bin"
mkdir -p "dist/${PKG_DIR}/usr/share/applications"
mkdir -p "dist/${PKG_DIR}/usr/share/icons/hicolor/64x64/apps"
mkdir -p "dist/${PKG_DIR}/usr/share/icons/hicolor/128x128/apps"
mkdir -p "dist/${PKG_DIR}/usr/share/icons/hicolor/256x256/apps"
mkdir -p "dist/${PKG_DIR}/usr/share/doc/markview"
mkdir -p "dist/${PKG_DIR}/usr/share/metainfo"

# Copy control file
cat > "dist/${PKG_DIR}/DEBIAN/control" << EOF
Package: markview
Version: ${VERSION}
Section: editors
Priority: optional
Architecture: ${ARCH}
Depends: libc6 (>= 2.17), libgl1, libx11-6, libxcursor1, libxrandr2, libxinerama1, libxi6, libxxf86vm1
Maintainer: Your Name <your.email@example.com>
Homepage: https://github.com/sbj-ee/markview
Description: A powerful markdown viewer and editor
 MarkView is a fast, cross-platform markdown viewer and editor
 with syntax highlighting, live reload, multiple themes, and
 advanced features like full-text search, backlinks, and templates.
EOF

# Copy desktop entry
cat > "dist/${PKG_DIR}/usr/share/applications/markview.desktop" << EOF
[Desktop Entry]
Name=MarkView
Comment=Markdown Viewer and Editor
Exec=markview %F
Icon=markview
Terminal=false
Type=Application
Categories=Office;TextEditor;Utility;
MimeType=text/markdown;text/x-markdown;
Keywords=markdown;editor;viewer;notes;
StartupNotify=true
EOF

# Copy icons
cp assets/logo-64.png "dist/${PKG_DIR}/usr/share/icons/hicolor/64x64/apps/markview.png"
cp assets/logo-128.png "dist/${PKG_DIR}/usr/share/icons/hicolor/128x128/apps/markview.png"
cp assets/logo-256.png "dist/${PKG_DIR}/usr/share/icons/hicolor/256x256/apps/markview.png"

# Copy documentation
cp README.md "dist/${PKG_DIR}/usr/share/doc/markview/"
cat > "dist/${PKG_DIR}/usr/share/doc/markview/copyright" << EOF
Format: https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/
Upstream-Name: markview
Source: https://github.com/sbj-ee/markview

Files: *
Copyright: 2026 Your Name
License: MIT
EOF

# Set permissions
chmod 755 "dist/${PKG_DIR}/usr/bin/markview"
chmod 644 "dist/${PKG_DIR}/DEBIAN/control"

# Build the package
echo "Building .deb package..."
dpkg-deb --build "dist/${PKG_DIR}"

# Validate with lintian (optional)
if command -v lintian &> /dev/null; then
    echo "Validating package..."
    lintian "dist/${PKG_NAME}_${VERSION}_${ARCH}.deb" || true
fi

echo "Package built: dist/${PKG_NAME}_${VERSION}_${ARCH}.deb"
```

### Step 5: Build and Install

```bash
# Build the package
chmod +x scripts/build-deb.sh
./scripts/build-deb.sh

# Install
sudo dpkg -i dist/markview_1.0.0_amd64.deb

# Or add to a local apt repository
```

---

## macOS .dmg Package

### Using Fyne Package Tool

Fyne provides a built-in packaging tool:

```bash
# Install fyne CLI
go install fyne.io/fyne/v2/cmd/fyne@latest

# Package as .app bundle
fyne package -os darwin -icon assets/logo-256.png -name MarkView -appID com.sbj-ee.markview
```

### Step 1: Create .app Bundle Structure

```
MarkView.app/
└── Contents/
    ├── Info.plist
    ├── MacOS/
    │   └── markview          # Binary
    └── Resources/
        ├── markview.icns     # App icon
        └── Assets/           # Additional resources
```

### Step 2: Create Info.plist

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>en</string>
    <key>CFBundleExecutable</key>
    <string>markview</string>
    <key>CFBundleIconFile</key>
    <string>markview.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.sbj-ee.markview</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>MarkView</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.14</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>CFBundleDocumentTypes</key>
    <array>
        <dict>
            <key>CFBundleTypeName</key>
            <string>Markdown Document</string>
            <key>CFBundleTypeRole</key>
            <string>Editor</string>
            <key>LSHandlerRank</key>
            <string>Default</string>
            <key>LSItemContentTypes</key>
            <array>
                <string>net.daringfireball.markdown</string>
                <string>public.plain-text</string>
            </array>
            <key>CFBundleTypeExtensions</key>
            <array>
                <string>md</string>
                <string>markdown</string>
            </array>
        </dict>
    </array>
    <key>UTImportedTypeDeclarations</key>
    <array>
        <dict>
            <key>UTTypeIdentifier</key>
            <string>net.daringfireball.markdown</string>
            <key>UTTypeDescription</key>
            <string>Markdown Document</string>
            <key>UTTypeConformsTo</key>
            <array>
                <string>public.plain-text</string>
            </array>
            <key>UTTypeTagSpecification</key>
            <dict>
                <key>public.filename-extension</key>
                <array>
                    <string>md</string>
                    <string>markdown</string>
                </array>
            </dict>
        </dict>
    </array>
</dict>
</plist>
```

### Step 3: Create Icon (.icns)

Convert PNG to icns format:

```bash
# Create iconset directory
mkdir markview.iconset

# Create required sizes
sips -z 16 16     assets/logo-256.png --out markview.iconset/icon_16x16.png
sips -z 32 32     assets/logo-256.png --out markview.iconset/icon_16x16@2x.png
sips -z 32 32     assets/logo-256.png --out markview.iconset/icon_32x32.png
sips -z 64 64     assets/logo-256.png --out markview.iconset/icon_32x32@2x.png
sips -z 128 128   assets/logo-256.png --out markview.iconset/icon_128x128.png
sips -z 256 256   assets/logo-256.png --out markview.iconset/icon_128x128@2x.png
sips -z 256 256   assets/logo-256.png --out markview.iconset/icon_256x256.png
sips -z 512 512   assets/logo-256.png --out markview.iconset/icon_256x256@2x.png
sips -z 512 512   assets/logo-256.png --out markview.iconset/icon_512x512.png
sips -z 1024 1024 assets/logo-256.png --out markview.iconset/icon_512x512@2x.png

# Convert to icns
iconutil -c icns markview.iconset -o markview.icns

# Clean up
rm -rf markview.iconset
```

### Step 4: Build Script for DMG

Create `scripts/build-dmg.sh`:

```bash
#!/bin/bash
set -e

VERSION="1.0.0"
APP_NAME="MarkView"
DMG_NAME="MarkView-${VERSION}"

echo "Building MarkView .dmg package v${VERSION}..."

# Clean previous builds
rm -rf "dist/${APP_NAME}.app" "dist/${DMG_NAME}.dmg"
mkdir -p dist

# Build for macOS (both architectures)
echo "Compiling binary for macOS..."

# For Apple Silicon (arm64)
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o "dist/markview-arm64" ./cmd/markview

# For Intel (amd64)
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o "dist/markview-amd64" ./cmd/markview

# Create universal binary
lipo -create -output "dist/markview" "dist/markview-arm64" "dist/markview-amd64"
rm "dist/markview-arm64" "dist/markview-amd64"

# Create .app bundle structure
echo "Creating .app bundle..."
mkdir -p "dist/${APP_NAME}.app/Contents/MacOS"
mkdir -p "dist/${APP_NAME}.app/Contents/Resources"

# Copy binary
cp "dist/markview" "dist/${APP_NAME}.app/Contents/MacOS/"
chmod +x "dist/${APP_NAME}.app/Contents/MacOS/markview"

# Copy/create icon
if [ -f "assets/markview.icns" ]; then
    cp "assets/markview.icns" "dist/${APP_NAME}.app/Contents/Resources/"
else
    echo "Warning: markview.icns not found, using placeholder"
fi

# Create Info.plist
cat > "dist/${APP_NAME}.app/Contents/Info.plist" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>en</string>
    <key>CFBundleExecutable</key>
    <string>markview</string>
    <key>CFBundleIconFile</key>
    <string>markview.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.sbj-ee.markview</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>MarkView</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.14</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>CFBundleDocumentTypes</key>
    <array>
        <dict>
            <key>CFBundleTypeName</key>
            <string>Markdown Document</string>
            <key>CFBundleTypeRole</key>
            <string>Editor</string>
            <key>CFBundleTypeExtensions</key>
            <array>
                <string>md</string>
                <string>markdown</string>
            </array>
        </dict>
    </array>
</dict>
</plist>
EOF

# Create DMG
echo "Creating .dmg..."

if command -v create-dmg &> /dev/null; then
    # Use create-dmg for a prettier DMG
    create-dmg \
        --volname "${APP_NAME}" \
        --volicon "assets/markview.icns" \
        --window-pos 200 120 \
        --window-size 600 400 \
        --icon-size 100 \
        --icon "${APP_NAME}.app" 150 190 \
        --hide-extension "${APP_NAME}.app" \
        --app-drop-link 450 190 \
        --no-internet-enable \
        "dist/${DMG_NAME}.dmg" \
        "dist/${APP_NAME}.app"
else
    # Fallback to hdiutil
    hdiutil create -volname "${APP_NAME}" \
        -srcfolder "dist/${APP_NAME}.app" \
        -ov -format UDZO \
        "dist/${DMG_NAME}.dmg"
fi

echo "Package built: dist/${DMG_NAME}.dmg"
```

### Step 5: Code Signing (Optional but Recommended)

For distribution outside the App Store:

```bash
# Sign the app
codesign --deep --force --verify --verbose \
    --sign "Developer ID Application: Your Name (TEAM_ID)" \
    "dist/MarkView.app"

# Notarize the app (requires Apple Developer account)
xcrun notarytool submit "dist/MarkView-1.0.0.dmg" \
    --apple-id "your@email.com" \
    --team-id "TEAM_ID" \
    --password "@keychain:AC_PASSWORD" \
    --wait

# Staple the notarization ticket
xcrun stapler staple "dist/MarkView-1.0.0.dmg"
```

---

## Automation with Makefile

Add these targets to the Makefile:

```makefile
VERSION := 1.0.0
ARCH := amd64

.PHONY: package-deb package-dmg package-all

# Build .deb package for Linux
package-deb:
	@echo "Building .deb package..."
	@chmod +x scripts/build-deb.sh
	@./scripts/build-deb.sh

# Build .dmg package for macOS (must run on macOS)
package-dmg:
	@echo "Building .dmg package..."
	@chmod +x scripts/build-dmg.sh
	@./scripts/build-dmg.sh

# Build all packages
package-all: package-deb package-dmg

# Cross-compile using fyne-cross (requires Docker)
package-cross-linux:
	fyne-cross linux -arch=amd64,arm64 -app-id com.sbj-ee.markview

package-cross-darwin:
	fyne-cross darwin -arch=amd64,arm64 -app-id com.sbj-ee.markview

# Clean dist directory
clean-dist:
	rm -rf dist/
```

---

## CI/CD Integration

### GitHub Actions Workflow

Create `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build-linux:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y gcc libgl1-mesa-dev xorg-dev

      - name: Build .deb package
        run: |
          chmod +x scripts/build-deb.sh
          ./scripts/build-deb.sh

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: markview-linux-deb
          path: dist/*.deb

  build-macos:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install create-dmg
        run: brew install create-dmg

      - name: Build .dmg package
        run: |
          chmod +x scripts/build-dmg.sh
          ./scripts/build-dmg.sh

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: markview-macos-dmg
          path: dist/*.dmg

  release:
    needs: [build-linux, build-macos]
    runs-on: ubuntu-latest
    steps:
      - name: Download artifacts
        uses: actions/download-artifact@v4

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            markview-linux-deb/*.deb
            markview-macos-dmg/*.dmg
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Summary

| Platform | Package Format | Build Command | Output |
|----------|---------------|---------------|--------|
| Linux (Debian/Ubuntu) | .deb | `make package-deb` | `dist/markview_1.0.0_amd64.deb` |
| macOS | .dmg | `make package-dmg` | `dist/MarkView-1.0.0.dmg` |

### Quick Start

```bash
# Linux
./scripts/build-deb.sh
sudo dpkg -i dist/markview_1.0.0_amd64.deb

# macOS
./scripts/build-dmg.sh
open dist/MarkView-1.0.0.dmg
```
