#!/bin/bash
set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && cd .. && pwd)"
cd "$PROJECT_ROOT"

export CGO_ENABLED=0 GOARCH=amd64
APP_NAME="fileserver"
MAIN_PATH="./cmd/fileserver/main.go"
BUILD_SCRIPT="$PROJECT_ROOT/build/build.sh"

echo "Building for all platforms..."
echo ""

"$BUILD_SCRIPT" windows $APP_NAME $MAIN_PATH
echo ""

"$BUILD_SCRIPT" darwin $APP_NAME $MAIN_PATH --force-macos
echo ""

"$BUILD_SCRIPT" linux $APP_NAME $MAIN_PATH
echo ""

echo "✓ All builds completed successfully!"
echo "Binaries: bin/fileserver_{win,mac,lnx}"