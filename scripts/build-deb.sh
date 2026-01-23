#!/bin/bash
#
# Build script for MarkView .deb package (Debian/Ubuntu)
#
set -e

VERSION="${VERSION:-1.0.0}"
ARCH="${ARCH:-amd64}"
PKG_NAME="markview"
PKG_DIR="${PKG_NAME}_${VERSION}_${ARCH}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "============================================"
echo "Building MarkView .deb package"
echo "Version: ${VERSION}"
echo "Architecture: ${ARCH}"
echo "============================================"

# Clean previous builds
echo "[1/7] Cleaning previous builds..."
rm -rf "dist/deb/${PKG_DIR}" "dist/deb/${PKG_NAME}_${VERSION}_${ARCH}.deb"
mkdir -p "dist/deb/${PKG_DIR}"

# Build the binary
echo "[2/7] Compiling binary..."
CGO_ENABLED=1 GOOS=linux GOARCH=${ARCH} go build -ldflags="-s -w -X main.version=${VERSION}" -o "dist/deb/${PKG_DIR}/usr/bin/markview" ./cmd/markview

# Create directory structure
echo "[3/7] Creating package structure..."
mkdir -p "dist/deb/${PKG_DIR}/DEBIAN"
mkdir -p "dist/deb/${PKG_DIR}/usr/bin"
mkdir -p "dist/deb/${PKG_DIR}/usr/share/applications"
mkdir -p "dist/deb/${PKG_DIR}/usr/share/icons/hicolor/64x64/apps"
mkdir -p "dist/deb/${PKG_DIR}/usr/share/icons/hicolor/128x128/apps"
mkdir -p "dist/deb/${PKG_DIR}/usr/share/icons/hicolor/256x256/apps"
mkdir -p "dist/deb/${PKG_DIR}/usr/share/icons/hicolor/scalable/apps"
mkdir -p "dist/deb/${PKG_DIR}/usr/share/doc/markview"
mkdir -p "dist/deb/${PKG_DIR}/usr/share/metainfo"
mkdir -p "dist/deb/${PKG_DIR}/usr/share/mime/packages"

# Create control file
echo "[4/7] Creating package metadata..."
cat > "dist/deb/${PKG_DIR}/DEBIAN/control" << EOF
Package: markview
Version: ${VERSION}
Section: editors
Priority: optional
Architecture: ${ARCH}
Depends: libc6 (>= 2.17), libgl1 | libgl1-mesa-glx, libx11-6, libxcursor1, libxrandr2, libxinerama1, libxi6, libxxf86vm1
Recommends: aspell, wkhtmltopdf, pandoc
Maintainer: sbj-ee <sbj-ee@users.noreply.github.com>
Homepage: https://github.com/sbj-ee/markview
Description: A powerful markdown viewer and editor
 MarkView is a fast, cross-platform markdown viewer and editor with
 syntax highlighting, live reload, multiple themes, and advanced features.
 .
 Features include:
  - Full Markdown and GitHub Flavored Markdown support
  - Syntax highlighting for 250+ languages
  - 8 color themes (Light, Dark, Nord, Solarized, Monokai, etc.)
  - Live file reload with auto-watch
  - Edit mode with formatting toolbar
  - Split view (side-by-side editor and preview)
  - Quick file switcher and full-text search
  - Document templates and tag support
  - Export to HTML, PDF, DOCX, RTF
EOF

# Create postinst script
cat > "dist/deb/${PKG_DIR}/DEBIAN/postinst" << 'EOF'
#!/bin/bash
set -e

# Update icon cache
if command -v gtk-update-icon-cache &> /dev/null; then
    gtk-update-icon-cache -f -t /usr/share/icons/hicolor || true
fi

# Update desktop database
if command -v update-desktop-database &> /dev/null; then
    update-desktop-database /usr/share/applications || true
fi

# Update MIME database
if command -v update-mime-database &> /dev/null; then
    update-mime-database /usr/share/mime || true
fi

exit 0
EOF
chmod 755 "dist/deb/${PKG_DIR}/DEBIAN/postinst"

# Create postrm script
cat > "dist/deb/${PKG_DIR}/DEBIAN/postrm" << 'EOF'
#!/bin/bash
set -e

# Update icon cache
if command -v gtk-update-icon-cache &> /dev/null; then
    gtk-update-icon-cache -f -t /usr/share/icons/hicolor || true
fi

# Update desktop database
if command -v update-desktop-database &> /dev/null; then
    update-desktop-database /usr/share/applications || true
fi

exit 0
EOF
chmod 755 "dist/deb/${PKG_DIR}/DEBIAN/postrm"

# Create desktop entry
cat > "dist/deb/${PKG_DIR}/usr/share/applications/markview.desktop" << EOF
[Desktop Entry]
Name=MarkView
GenericName=Markdown Editor
Comment=A powerful markdown viewer and editor
Exec=markview %F
Icon=markview
Terminal=false
Type=Application
Categories=Office;TextEditor;Utility;Development;
MimeType=text/markdown;text/x-markdown;
Keywords=markdown;editor;viewer;notes;writing;documentation;
StartupNotify=true
StartupWMClass=markview
EOF

# Create MIME type definition
cat > "dist/deb/${PKG_DIR}/usr/share/mime/packages/markview.xml" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<mime-info xmlns="http://www.freedesktop.org/standards/shared-mime-info">
  <mime-type type="text/markdown">
    <comment>Markdown document</comment>
    <glob pattern="*.md"/>
    <glob pattern="*.markdown"/>
    <glob pattern="*.mdown"/>
    <glob pattern="*.mkd"/>
  </mime-type>
</mime-info>
EOF

# Create AppStream metadata
cat > "dist/deb/${PKG_DIR}/usr/share/metainfo/com.sbj-ee.markview.metainfo.xml" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<component type="desktop-application">
  <id>com.sbj-ee.markview</id>
  <name>MarkView</name>
  <summary>A powerful markdown viewer and editor</summary>
  <metadata_license>MIT</metadata_license>
  <project_license>MIT</project_license>
  <developer_name>sbj-ee</developer_name>
  <description>
    <p>
      MarkView is a fast, cross-platform markdown viewer and editor with
      syntax highlighting, live reload, and multiple themes.
    </p>
    <p>Features include:</p>
    <ul>
      <li>Full Markdown and GitHub Flavored Markdown support</li>
      <li>Syntax highlighting for 250+ languages</li>
      <li>8 color themes including Dark, Nord, Solarized, and Monokai</li>
      <li>Live file reload with auto-watch</li>
      <li>Edit mode with formatting toolbar</li>
      <li>Split view for side-by-side editing and preview</li>
      <li>Quick file switcher and full-text search</li>
      <li>Document templates for common formats</li>
      <li>Tag support for document organization</li>
      <li>Export to HTML, PDF, DOCX, and RTF</li>
    </ul>
  </description>
  <launchable type="desktop-id">markview.desktop</launchable>
  <url type="homepage">https://github.com/sbj-ee/markview</url>
  <url type="bugtracker">https://github.com/sbj-ee/markview/issues</url>
  <provides>
    <binary>markview</binary>
  </provides>
  <releases>
    <release version="${VERSION}" date="$(date +%Y-%m-%d)">
      <description>
        <p>Release version ${VERSION}</p>
      </description>
    </release>
  </releases>
  <content_rating type="oars-1.1"/>
  <categories>
    <category>Office</category>
    <category>TextEditor</category>
  </categories>
  <keywords>
    <keyword>markdown</keyword>
    <keyword>editor</keyword>
    <keyword>viewer</keyword>
    <keyword>notes</keyword>
    <keyword>writing</keyword>
  </keywords>
</component>
EOF

# Copy icons
echo "[5/7] Copying icons..."
if [ -f "assets/logo-64.png" ]; then
    cp "assets/logo-64.png" "dist/deb/${PKG_DIR}/usr/share/icons/hicolor/64x64/apps/markview.png"
fi
if [ -f "assets/logo-128.png" ]; then
    cp "assets/logo-128.png" "dist/deb/${PKG_DIR}/usr/share/icons/hicolor/128x128/apps/markview.png"
fi
if [ -f "assets/logo-256.png" ]; then
    cp "assets/logo-256.png" "dist/deb/${PKG_DIR}/usr/share/icons/hicolor/256x256/apps/markview.png"
fi
if [ -f "assets/logo.svg" ]; then
    cp "assets/logo.svg" "dist/deb/${PKG_DIR}/usr/share/icons/hicolor/scalable/apps/markview.svg"
fi

# Copy documentation
echo "[6/7] Copying documentation..."
cp README.md "dist/deb/${PKG_DIR}/usr/share/doc/markview/"

cat > "dist/deb/${PKG_DIR}/usr/share/doc/markview/copyright" << EOF
Format: https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/
Upstream-Name: markview
Upstream-Contact: sbj-ee <sbj-ee@users.noreply.github.com>
Source: https://github.com/sbj-ee/markview

Files: *
Copyright: $(date +%Y) sbj-ee
License: MIT

License: MIT
 Permission is hereby granted, free of charge, to any person obtaining a
 copy of this software and associated documentation files (the "Software"),
 to deal in the Software without restriction, including without limitation
 the rights to use, copy, modify, merge, publish, distribute, sublicense,
 and/or sell copies of the Software, and to permit persons to whom the
 Software is furnished to do so, subject to the following conditions:
 .
 The above copyright notice and this permission notice shall be included
 in all copies or substantial portions of the Software.
 .
 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
 OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
 SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
EOF

# Set permissions
chmod 755 "dist/deb/${PKG_DIR}/usr/bin/markview"
chmod 644 "dist/deb/${PKG_DIR}/DEBIAN/control"
find "dist/deb/${PKG_DIR}/usr/share" -type f -exec chmod 644 {} \;
find "dist/deb/${PKG_DIR}/usr/share" -type d -exec chmod 755 {} \;

# Build the package
echo "[7/7] Building .deb package..."
if command -v fakeroot &> /dev/null; then
    fakeroot dpkg-deb --build "dist/deb/${PKG_DIR}" "dist/deb/${PKG_NAME}_${VERSION}_${ARCH}.deb"
else
    dpkg-deb --build "dist/deb/${PKG_DIR}" "dist/deb/${PKG_NAME}_${VERSION}_${ARCH}.deb"
fi

# Clean up build directory
rm -rf "dist/deb/${PKG_DIR}"

# Validate with lintian (optional)
if command -v lintian &> /dev/null; then
    echo ""
    echo "Validating package with lintian..."
    lintian --no-tag-display-limit "dist/deb/${PKG_NAME}_${VERSION}_${ARCH}.deb" || true
fi

echo ""
echo "============================================"
echo "Build complete!"
echo "Package: dist/${PKG_NAME}_${VERSION}_${ARCH}.deb"
echo ""
echo "Install with:"
echo "  sudo dpkg -i dist/${PKG_NAME}_${VERSION}_${ARCH}.deb"
echo "============================================"
