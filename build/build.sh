#!/bin/bash
set -e

# Использование: ./build.sh <GOOS> <APP_NAME> <MAIN_PATH> [UPX_FLAGS]
# Примеры:
#   ./build.sh linux fileserver ./cmd/fileserver/main.go
#   ./build.sh darwin fileserver ./cmd/fileserver/main.go --force-macos
#   ./build.sh windows fileserver ./cmd/fileserver/main.go

if [ $# -lt 3 ]; then
    echo "Usage: $0 <GOOS> <APP_NAME> <MAIN_PATH> [UPX_FLAGS]"
    echo "Example: $0 linux fileserver ./cmd/fileserver/main.go"
    exit 1
fi

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && cd .. && pwd)"
cd "$PROJECT_ROOT"

mkdir -p ./tmp ./bin

GOOS=$1
APP_NAME=$2
MAIN_PATH=$3
UPX_FLAGS="${4:-}"

# Определяем суффикс платформы
case "$GOOS" in
    linux)  SUFFIX="_lnx" ;;
    darwin) SUFFIX="_mac" ;;
    windows) SUFFIX="_win" ;;
    *)      echo "Unknown OS: $GOOS"; exit 1 ;;
esac

FILE_NAME="${APP_NAME}${SUFFIX}"
TEMP_FILE="./tmp/${FILE_NAME}"

echo "Building $FILE_NAME for $GOOS..."
GOOS=$GOOS \
go build -ldflags="-s -w" -trimpath -o "$TEMP_FILE" "$MAIN_PATH"

echo "Compressing $FILE_NAME with UPX..."
if [ -n "$UPX_FLAGS" ]; then
    upx --lzma $UPX_FLAGS -o "./bin/${FILE_NAME}" "$TEMP_FILE"
else
    upx --lzma -o "./bin/${FILE_NAME}" "$TEMP_FILE"
fi

rm "$TEMP_FILE"
echo "✓ Done: bin/$FILE_NAME"
