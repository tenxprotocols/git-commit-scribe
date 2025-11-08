#!/bin/bash

# Release script for git-commit-scribe
# Creates a new release with cross-platform binaries

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BINARY_NAME="gscribe"
OUTPUT_DIR="dist"
REPO_URL="github.com/tenxprotocols/git-commit-scribe"

# Platforms to build for
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

# Print usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS] <version>

Create a new release for git-commit-scribe

Arguments:
  version          Version to release (e.g., v1.2.3 or 1.2.3)

Options:
  -h, --help       Show this help message
  -d, --dry-run    Dry run - don't create tags or push
  -s, --skip-build Skip building binaries
  --no-tag         Don't create and push git tag
  --github         Create GitHub release (requires gh CLI)

Examples:
  $0 v1.2.3                    # Create release v1.2.3
  $0 1.2.3                     # Create release v1.2.3 (auto-adds 'v')
  $0 -d v1.2.3                 # Dry run
  $0 --github v1.2.3           # Create release and GitHub release

EOF
    exit 1
}

# Parse command line arguments
DRY_RUN=false
SKIP_BUILD=false
CREATE_TAG=true
CREATE_GITHUB_RELEASE=false
VERSION=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -s|--skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --no-tag)
            CREATE_TAG=false
            shift
            ;;
        --github)
            CREATE_GITHUB_RELEASE=true
            shift
            ;;
        *)
            if [[ -z "$VERSION" ]]; then
                VERSION="$1"
            else
                echo -e "${RED}Error: Unknown argument: $1${NC}"
                usage
            fi
            shift
            ;;
    esac
done

# Check if version was provided
if [[ -z "$VERSION" ]]; then
    echo -e "${RED}Error: Version is required${NC}"
    usage
fi

# Normalize version (add 'v' prefix if missing)
if [[ ! "$VERSION" =~ ^v ]]; then
    VERSION="v${VERSION}"
fi

# Validate version format
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
    echo -e "${RED}Error: Invalid version format: ${VERSION}${NC}"
    echo "Version must be in format: v1.2.3 or v1.2.3-beta1"
    exit 1
fi

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}git-commit-scribe Release${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "Version:       ${GREEN}${VERSION}${NC}"
echo -e "Dry Run:       ${DRY_RUN}"
echo -e "Skip Build:    ${SKIP_BUILD}"
echo -e "Create Tag:    ${CREATE_TAG}"
echo -e "GitHub Release: ${CREATE_GITHUB_RELEASE}"
echo ""

# Check if we're in a git repository
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Check for uncommitted changes
if [[ -n $(git status -s) ]]; then
    echo -e "${YELLOW}Warning: You have uncommitted changes${NC}"
    git status -s
    echo ""
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
fi

# Check if tag already exists
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo -e "${RED}Error: Tag ${VERSION} already exists${NC}"
    exit 1
fi

# Get current commit
COMMIT=$(git rev-parse --short HEAD)
BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')

echo -e "${YELLOW}Preparing release...${NC}"
echo "Commit: ${COMMIT}"
echo "Date: ${BUILD_DATE}"
echo ""

# Build binaries
if [[ "$SKIP_BUILD" == false ]]; then
    echo -e "${GREEN}Building binaries...${NC}"
    
    # Clean and create dist directory
    rm -rf "${OUTPUT_DIR}"
    mkdir -p "${OUTPUT_DIR}"
    
    # Build for each platform
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r -a platform_split <<< "$platform"
        GOOS="${platform_split[0]}"
        GOARCH="${platform_split[1]}"
        
        output_name="${BINARY_NAME}-${VERSION}-${GOOS}-${GOARCH}"
        
        if [[ "$GOOS" == "windows" ]]; then
            output_name="${output_name}.exe"
        fi
        
        echo -e "${YELLOW}  Building ${GOOS}/${GOARCH}...${NC}"
        
        LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildDate=${BUILD_DATE}"
        
        if ! GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="${LDFLAGS}" -o "${OUTPUT_DIR}/${output_name}" ./cmd/gscribe; then
            echo -e "${RED}  ✗ Build failed for ${GOOS}/${GOARCH}${NC}"
            exit 1
        fi
        
        # Create archive
        archive_name="${BINARY_NAME}-${VERSION}-${GOOS}-${GOARCH}"
        
        pushd "${OUTPUT_DIR}" > /dev/null
        
        if [[ "$GOOS" == "windows" ]]; then
            zip -q "${archive_name}.zip" "${output_name}"
            echo -e "${GREEN}  ✓ Created ${archive_name}.zip${NC}"
        else
            tar czf "${archive_name}.tar.gz" "${output_name}"
            echo -e "${GREEN}  ✓ Created ${archive_name}.tar.gz${NC}"
        fi
        
        # Remove unarchived binary
        rm "${output_name}"
        
        popd > /dev/null
    done
    
    # Generate checksums
    echo ""
    echo -e "${YELLOW}Generating checksums...${NC}"
    pushd "${OUTPUT_DIR}" > /dev/null
    
    if command -v shasum &> /dev/null; then
        shasum -a 256 * > SHA256SUMS
    elif command -v sha256sum &> /dev/null; then
        sha256sum * > SHA256SUMS
    else
        echo -e "${YELLOW}Warning: Neither shasum nor sha256sum found, skipping checksums${NC}"
    fi
    
    popd > /dev/null
    
    echo -e "${GREEN}✓ All binaries built successfully${NC}"
    echo ""
    echo "Release artifacts in ${OUTPUT_DIR}/:"
    ls -lh "${OUTPUT_DIR}"
    echo ""
fi

# Create and push git tag
if [[ "$CREATE_TAG" == true ]]; then
    if [[ "$DRY_RUN" == true ]]; then
        echo -e "${YELLOW}[DRY RUN] Would create tag: ${VERSION}${NC}"
        echo -e "${YELLOW}[DRY RUN] Would push tag to origin${NC}"
    else
        echo -e "${GREEN}Creating git tag...${NC}"
        git tag -a "$VERSION" -m "Release $VERSION"
        
        echo -e "${GREEN}Pushing tag to origin...${NC}"
        git push origin "$VERSION"
        
        echo -e "${GREEN}✓ Tag ${VERSION} created and pushed${NC}"
    fi
    echo ""
fi

# Create GitHub release
if [[ "$CREATE_GITHUB_RELEASE" == true ]]; then
    if ! command -v gh &> /dev/null; then
        echo -e "${YELLOW}Warning: GitHub CLI (gh) not found. Skipping GitHub release.${NC}"
        echo "Install from: https://cli.github.com/"
    elif [[ "$DRY_RUN" == true ]]; then
        echo -e "${YELLOW}[DRY RUN] Would create GitHub release for ${VERSION}${NC}"
    else
        echo -e "${GREEN}Creating GitHub release...${NC}"
        
        # Create release notes file
        RELEASE_NOTES=$(mktemp)
        cat > "$RELEASE_NOTES" << EOF
Release $VERSION

## Installation

### Homebrew (macOS/Linux)
\`\`\`bash
brew install tenxprotocols/tap/gscribe
\`\`\`

### Go Install
\`\`\`bash
go install ${REPO_URL}/cmd/gscribe@${VERSION}
\`\`\`

### Binary Download
Download the appropriate binary for your platform from the assets below.

## Checksums
See SHA256SUMS file for checksums of all release artifacts.
EOF
        
        # Create GitHub release with artifacts
        gh release create "$VERSION" \
            --title "Release $VERSION" \
            --notes-file "$RELEASE_NOTES" \
            "${OUTPUT_DIR}"/*
        
        rm "$RELEASE_NOTES"
        
        echo -e "${GREEN}✓ GitHub release created${NC}"
    fi
    echo ""
fi

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Release ${VERSION} completed successfully!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

if [[ "$DRY_RUN" == false ]]; then
    echo "Next steps:"
    echo "  1. Verify the release on GitHub"
    echo "  2. Update announcement/changelog if needed"
    if [[ "$CREATE_GITHUB_RELEASE" == false ]]; then
        echo "  3. Create GitHub release manually with: gh release create ${VERSION}"
    fi
fi
