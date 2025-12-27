#!/usr/bin/env bash
set -Eeuo pipefail

info() { 
    echo -e "\033[36m==>\033[0m $*"; 
}

success() { 
    echo -e "\033[32m==>\033[0m $*"; 
}

warn() { 
    echo -e "\033[33mWARN:\033[0m $*" >&2; 
}

fail() { 
    echo -e "\033[31mERROR:\033[0m $*" >&2; exit 1; 
}

trap 'fail "Build failed on line $LINENO"' ERR

BINARY_NAME="${BINARY_NAME:-vsharp}"
SOURCE_DIR="${SOURCE_DIR:-./cmd/vsharp}"
DIST_DIR="${DIST_DIR:-./dist}"

VERSION="${VERSION:-0.1.0}"
BUILD_TYPE="${BUILD_TYPE:-alpha}"
FULL_VERSION="${VERSION}-${BUILD_TYPE}"

command -v go >/dev/null 2>&1 || fail "Go is not installed"
command -v shasum >/dev/null 2>&1 || fail "shasum not installed"

info "Go version: $(go version)"

GIT_COMMIT_HASH="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
BUILD_DATE="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

[[ -f go.mod ]] || fail "go.mod not found (not a Go module)"
[[ -d "$SOURCE_DIR" ]] || fail "Source directory not found: $SOURCE_DIR"

LDFLAGS=(
  "-X vsharp/internal/version.Version=$FULL_VERSION"
  "-X vsharp/internal/version.GitCommit=$GIT_COMMIT_HASH"
  "-X vsharp/internal/version.BuildDate=$BUILD_DATE"
  "-s"
  "-w"
)

TARGETS=(
  "linux amd64"
  "linux arm64"
  "linux arm"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
  "windows arm64"
)

info "Preparing dist directory"
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

info "Version: $FULL_VERSION"
info "Git commit: $GIT_COMMIT_HASH"
info "Build date: $BUILD_DATE"

for TARGET in "${TARGETS[@]}"; do
    read -r GOOS GOARCH <<< "$TARGET"

    OUTPUT="$DIST_DIR/${BINARY_NAME}-${GOOS}-${GOARCH}"
    [[ "$GOOS" == "windows" ]] && OUTPUT+=".exe"

    info "Building $OUTPUT"

    env GOOS="$GOOS" GOARCH="$GOARCH" \
        go build \
          -trimpath \
          -ldflags "${LDFLAGS[*]}" \
          -o "$OUTPUT" \
          "$SOURCE_DIR"

    shasum -a 256 "$OUTPUT" > "${OUTPUT}.sha256"
done

echo "$FULL_VERSION"      > "$DIST_DIR/VERSION.txt"
echo "$GIT_COMMIT_HASH"   > "$DIST_DIR/GIT_COMMIT.txt"
echo "$BUILD_DATE"        > "$DIST_DIR/BUILD_DATE.txt"

success "Build finished successfully"
success "Artifacts available in: $DIST_DIR"