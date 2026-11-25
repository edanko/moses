#!/bin/bash

# Build script for the moses CLI.
set -e

VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$BUILD_TIME"

echo "Building moses..."
echo "  version: $VERSION"
echo "  commit:  $COMMIT"

echo "Building for $(go env GOOS)/$(go env GOARCH)..."
go build -ldflags "$LDFLAGS" -o moses ./cmd/moses
echo "Built ./moses"

mkdir -p build
PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64"
for platform in $PLATFORMS; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}

    if [ "$GOOS" = "windows" ]; then
        EXTENSION=".exe"
    else
        EXTENSION=""
    fi

    OUTPUT="moses-${GOOS}-${GOARCH}${EXTENSION}"
    echo "Building $platform..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "build/$OUTPUT" ./cmd/moses
done

echo "All builds completed; binaries are in ./build/"