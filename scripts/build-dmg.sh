#!/bin/bash
#
# Build script for MarkView .dmg package (macOS)
#
# Usage:
#   VERSION=1.0.6 ./build-dmg.sh           # Universal binary (default)
#   VERSION=1.0.6 ARCH=arm64 ./build-dmg.sh  # Apple Silicon only
#   VERSION=1.0.6 ARCH=amd64 ./build-dmg.sh  # Intel only
#   VERSION=1.0.6 ARCH=all ./build-dmg.sh    # Build all three DMGs
#
set -e

VERSION="${VERSION:-1.0.0}"
ARCH="${ARCH:-universal}"  # universal, arm64, amd64, or all
APP_NAME="MarkView"
BUNDLE_ID="com.sbj-ee.markview"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Set DMG name based on architecture
case "$ARCH" in
    arm64)
        DMG_NAME="MarkView-${VERSION}-arm64"
        ARCH_DISPLAY="Apple Silicon (arm64)"
        ;;
    amd64|x86_64)
        ARCH="amd64"
        DMG_NAME="MarkView-${VERSION}-x86_64"
        ARCH_DISPLAY="Intel (x86_64)"
        ;;
    all)
        # Build all architectures by calling self recursively
        echo "Building all architectures..."
        ARCH=arm64 VERSION="$VERSION" "$0"
        ARCH=amd64 VERSION="$VERSION" "$0"
        ARCH=universal VERSION="$VERSION" "$0"
        echo ""
        echo "============================================"
        echo "All builds complete!"
        ls -lh dist/dmg/*.dmg
        echo "============================================"
        exit 0
        ;;
    *)
        DMG_NAME="MarkView-${VERSION}"
        ARCH_DISPLAY="Universal (arm64 + x86_64)"
        ARCH="universal"
        ;;
esac

echo "============================================"
echo "Building MarkView .dmg package"
echo "Version: ${VERSION}"
echo "Architecture: ${ARCH_DISPLAY}"
echo "============================================"

# Check if running on macOS
if [[ "$(uname)" != "Darwin" ]]; then
    echo "Error: This script must be run on macOS"
    exit 1
fi

# Clean previous builds for this architecture
echo "[1/6] Cleaning previous builds..."
rm -rf "dist/dmg/${APP_NAME}.app" "dist/dmg/${DMG_NAME}.dmg" "dist/dmg/markview-arm64" "dist/dmg/markview-amd64" "dist/dmg/markview"
mkdir -p dist/dmg

# Build based on architecture selection
if [ "$ARCH" = "universal" ]; then
    # Build for Apple Silicon (arm64)
    echo "[2/6] Compiling binary for Apple Silicon (arm64)..."
    CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o "dist/dmg/markview-arm64" ./cmd/markview

    # Build for Intel (amd64)
    echo "[2/6] Compiling binary for Intel (amd64)..."
    CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o "dist/dmg/markview-amd64" ./cmd/markview

    # Create universal binary
    echo "[2/6] Creating universal binary..."
    lipo -create -output "dist/dmg/markview" "dist/dmg/markview-arm64" "dist/dmg/markview-amd64"
    rm "dist/dmg/markview-arm64" "dist/dmg/markview-amd64"
    echo "   Binary architectures: $(lipo -archs dist/dmg/markview)"
elif [ "$ARCH" = "arm64" ]; then
    echo "[2/6] Compiling binary for Apple Silicon (arm64)..."
    CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o "dist/dmg/markview" ./cmd/markview
    echo "   Binary architecture: arm64"
elif [ "$ARCH" = "amd64" ]; then
    echo "[2/6] Compiling binary for Intel (x86_64)..."
    CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o "dist/dmg/markview" ./cmd/markview
    echo "   Binary architecture: x86_64"
fi

# Create .app bundle structure
echo "[3/6] Creating .app bundle..."
mkdir -p "dist/dmg/${APP_NAME}.app/Contents/MacOS"
mkdir -p "dist/dmg/${APP_NAME}.app/Contents/Resources"

# Copy binary
cp "dist/dmg/markview" "dist/dmg/${APP_NAME}.app/Contents/MacOS/"
chmod +x "dist/dmg/${APP_NAME}.app/Contents/MacOS/markview"

# Create/copy icon
echo "[4/6] Creating app icon..."
if [ -f "assets/markview.icns" ]; then
    cp "assets/markview.icns" "dist/dmg/${APP_NAME}.app/Contents/Resources/"
elif [ -f "assets/logo-256.png" ]; then
    # Create iconset from PNG
    ICONSET_DIR="dist/markview.iconset"
    mkdir -p "$ICONSET_DIR"

    # Generate required sizes
    sips -z 16 16     "assets/logo-256.png" --out "$ICONSET_DIR/icon_16x16.png" 2>/dev/null
    sips -z 32 32     "assets/logo-256.png" --out "$ICONSET_DIR/icon_16x16@2x.png" 2>/dev/null
    sips -z 32 32     "assets/logo-256.png" --out "$ICONSET_DIR/icon_32x32.png" 2>/dev/null
    sips -z 64 64     "assets/logo-256.png" --out "$ICONSET_DIR/icon_32x32@2x.png" 2>/dev/null
    sips -z 128 128   "assets/logo-256.png" --out "$ICONSET_DIR/icon_128x128.png" 2>/dev/null
    sips -z 256 256   "assets/logo-256.png" --out "$ICONSET_DIR/icon_128x128@2x.png" 2>/dev/null
    sips -z 256 256   "assets/logo-256.png" --out "$ICONSET_DIR/icon_256x256.png" 2>/dev/null
    cp "assets/logo-256.png" "$ICONSET_DIR/icon_256x256@2x.png"

    # Try to create 512px versions if possible
    if [ -f "assets/logo-512.png" ]; then
        cp "assets/logo-512.png" "$ICONSET_DIR/icon_512x512.png"
        cp "assets/logo-512.png" "$ICONSET_DIR/icon_512x512@2x.png"
    else
        cp "assets/logo-256.png" "$ICONSET_DIR/icon_512x512.png"
        cp "assets/logo-256.png" "$ICONSET_DIR/icon_512x512@2x.png"
    fi

    # Convert to icns
    iconutil -c icns "$ICONSET_DIR" -o "dist/dmg/${APP_NAME}.app/Contents/Resources/markview.icns"
    rm -rf "$ICONSET_DIR"

    # Also save to assets for future use
    cp "dist/dmg/${APP_NAME}.app/Contents/Resources/markview.icns" "assets/markview.icns" 2>/dev/null || true
else
    echo "   Warning: No icon found, app will use default icon"
fi

# Create Info.plist
echo "[5/6] Creating Info.plist..."
cat > "dist/dmg/${APP_NAME}.app/Contents/Info.plist" << EOF
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
    <string>${BUNDLE_ID}</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>${APP_NAME}</string>
    <key>CFBundleDisplayName</key>
    <string>${APP_NAME}</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>${VERSION}</string>
    <key>CFBundleVersion</key>
    <string>${VERSION}</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.14</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSSupportsAutomaticGraphicsSwitching</key>
    <true/>
    <key>LSApplicationCategoryType</key>
    <string>public.app-category.productivity</string>
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
                <string>mdown</string>
                <string>mkd</string>
            </array>
            <key>CFBundleTypeIconFile</key>
            <string>markview.icns</string>
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
                    <string>mdown</string>
                    <string>mkd</string>
                </array>
                <key>public.mime-type</key>
                <array>
                    <string>text/markdown</string>
                    <string>text/x-markdown</string>
                </array>
            </dict>
        </dict>
    </array>
    <key>NSHumanReadableCopyright</key>
    <string>Copyright © $(date +%Y) sbj-ee. MIT License.</string>
</dict>
</plist>
EOF

# Create PkgInfo
echo -n "APPL????" > "dist/dmg/${APP_NAME}.app/Contents/PkgInfo"

# Create DMG
echo "[6/6] Creating .dmg..."
rm -f "dist/dmg/${DMG_NAME}.dmg"

if command -v create-dmg &> /dev/null; then
    # Use create-dmg for a prettier DMG with background and icon positioning
    create-dmg \
        --volname "${APP_NAME}" \
        --window-pos 200 120 \
        --window-size 600 400 \
        --icon-size 100 \
        --icon "${APP_NAME}.app" 150 190 \
        --hide-extension "${APP_NAME}.app" \
        --app-drop-link 450 190 \
        --no-internet-enable \
        "dist/dmg/${DMG_NAME}.dmg" \
        "dist/dmg/${APP_NAME}.app" || {
            # Fallback if create-dmg fails
            echo "   create-dmg failed, falling back to hdiutil..."
            hdiutil create -volname "${APP_NAME}" \
                -srcfolder "dist/dmg/${APP_NAME}.app" \
                -ov -format UDZO \
                "dist/dmg/${DMG_NAME}.dmg"
        }
else
    # Fallback to hdiutil
    hdiutil create -volname "${APP_NAME}" \
        -srcfolder "dist/dmg/${APP_NAME}.app" \
        -ov -format UDZO \
        "dist/dmg/${DMG_NAME}.dmg"
fi

# Clean up intermediate files
rm -f "dist/dmg/markview"

echo ""
echo "============================================"
echo "Build complete!"
echo ""
echo "App bundle: dist/dmg/${APP_NAME}.app"
echo "DMG file:   dist/dmg/${DMG_NAME}.dmg"
echo ""
echo "To install:"
echo "  1. Open dist/${DMG_NAME}.dmg"
echo "  2. Drag ${APP_NAME} to Applications"
echo ""
echo "To test the app bundle directly:"
echo "  open dist/dmg/${APP_NAME}.app"
echo "============================================"

# Optional: Code signing reminder
echo ""
echo "NOTE: For distribution, you should sign the app:"
echo "  codesign --deep --force --verify --verbose \\"
echo "    --sign \"Developer ID Application: Your Name (TEAM_ID)\" \\"
echo "    \"dist/dmg/${APP_NAME}.app\""
echo ""
echo "And notarize for Gatekeeper:"
echo "  xcrun notarytool submit \"dist/dmg/${DMG_NAME}.dmg\" \\"
echo "    --apple-id \"your@email.com\" \\"
echo "    --team-id \"TEAM_ID\" \\"
echo "    --password \"@keychain:AC_PASSWORD\" \\"
echo "    --wait"
