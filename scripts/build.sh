#!/bin/bash

# Build script for aic

set -e

VERSION=${VERSION:-"dev"}
COMMIT=${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}
BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')

BINARY_NAME="gscribe"
OUTPUT_DIR="bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building ${BINARY_NAME}...${NC}"
echo "Version: ${VERSION}"
echo "Commit: ${COMMIT}"
echo "Build Date: ${BUILD_DATE}"
echo ""

# Create output directory
mkdir -p "${OUTPUT_DIR}"

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildDate=${BUILD_DATE}"

# Build for current platform
echo -e "${YELLOW}Building for $(go env GOOS)/$(go env GOARCH)...${NC}"
go build -ldflags="${LDFLAGS}" -o "${OUTPUT_DIR}/${BINARY_NAME}" ./cmd/gscribe

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful!${NC}"
    echo "Binary: ${OUTPUT_DIR}/${BINARY_NAME}"
    
    # Make binary executable
    chmod +x "${OUTPUT_DIR}/${BINARY_NAME}"
    
    # Show binary info
    echo ""
    echo "Binary size: $(du -h ${OUTPUT_DIR}/${BINARY_NAME} | cut -f1)"
    echo ""
    echo "To install, run: cp ${OUTPUT_DIR}/${BINARY_NAME} /usr/local/bin/"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi
