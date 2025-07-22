#!/bin/sh
set -eu

REPO="Cybernetics354/mayohttp"
DEFAULT_INSTALL_DIR="$HOME/.local/bin"
INSTALL_DIR=${DESTINATION:-$DEFAULT_INSTALL_DIR}
VERSION=${VERSION:-}
TMP_DIR=$(mktemp -d)

# --- Pre-checks ---
echo "✅ Checking required tools..."

check_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "❌ Required command not found: $1"
    echo "👉 Please install it and re-run the script."
    exit 1
  fi
}

check_command tar

if command -v curl >/dev/null 2>&1; then
    DLCMD="curl -sLO"
elif command -v wget >/dev/null 2>&1; then
    DLCMD="wget -q"
else
    echo "❌ Neither 'curl' nor 'wget' found."
    echo "👉 Please install one of them and re-run the script."
    exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
    SHACMD="sha256sum"
elif command -v shasum >/dev/null 2>&1; then
    SHACMD="shasum -a 256"
else
    echo "❌ Neither 'sha256sum' nor 'shasum' found."
    echo "👉 Please install one of them to verify checksums."
    exit 1
fi

# --- Parse arguments ---
while [ $# -gt 0 ]; do
  case "$1" in
    --destination)
      INSTALL_DIR="$2"
      shift 2
      ;;
    --version)
      VERSION="$2"
      shift 2
      ;;
    *)
      echo "❌ Unknown option: $1"
      echo "Usage: $0 [--destination <dir>] [--version <version>]"
      exit 1
      ;;
  esac
done

# --- Detect version ---
if [ -z "$VERSION" ]; then
  echo "🔍 Fetching latest version from GitHub..."
  if command -v curl >/dev/null 2>&1; then
    VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | awk -F'"tag_name"[[:space:]]*:[[:space:]]*"' 'NF>1 { split($2, a, "\""); print a[1] }')
  else
    VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | awk -F'"tag_name"[[:space:]]*:[[:space:]]*"' 'NF>1 { split($2, a, "\""); print a[1] }')
  fi
  if [ -z "$VERSION" ]; then
    echo "❌ Failed to fetch latest version"
    exit 1
  fi
fi

echo "📦 Installing mayohttp version: $VERSION"

# Build checksum file name based on version
CHECKSUM_FILE="mayohttp_${VERSION#v}_checksums.txt"

# --- Detect OS ---
echo "👉 Detecting OS..."
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux|darwin) ;;
    *)
        echo "❌ Unsupported OS: $OS"
        exit 1
        ;;
esac

# --- Detect architecture ---
echo "👉 Detecting architecture..."
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    i386|i686) ARCH="386" ;;
    arm64|aarch64) ARCH="arm64" ;;
    armv6l) ARCH="armv6" ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

ASSET="mayohttp_${VERSION#v}_${OS}_${ARCH}.tar.gz"

ORIGIN_PWD=$(pwd)
cd "$TMP_DIR"

cleanup() {
  echo "🧹 Cleaning up..."
  rm -rf "$TMP_DIR"
  cd "$ORIGIN_PWD"
}

echo "📥 Downloading mayohttp archive: $ASSET"
$DLCMD "https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"

echo "📥 Downloading checksum file..."
$DLCMD "https://github.com/${REPO}/releases/download/${VERSION}/${CHECKSUM_FILE}"

echo "🔑 Verifying checksum..."
grep " ${ASSET}$" "$CHECKSUM_FILE" | $SHACMD -c -

echo "📦 Extracting archive..."
tar -xzf "${ASSET}"

if [ ! -f mayohttp ]; then
  echo "❌ Expected file 'mayohttp' not found after extracting"
  cleanup
  exit 1
fi

echo "🚀 Installing to $INSTALL_DIR..."
mkdir -p "$INSTALL_DIR"
chmod +x mayohttp
mv mayohttp "$INSTALL_DIR/mayohttp"

# --- Check if INSTALL_DIR is in PATH ---
#if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
if ! command -v mayohttp >/dev/null 2>&1; then
    echo
    echo " *-----------------------------------------------------------------------------------------"
    echo " |  ⚠️  $INSTALL_DIR is not in your PATH."
    echo " |  👉 To fix, add this line to your shell profile (~/.bashrc, ~/.zshrc, or ~/.profile):"
    echo " |"
    echo " |      export PATH=\"$INSTALL_DIR:\$PATH\""
    echo " |"
    echo " |  Then run: source ~/.profile (or your shell config)"
    echo " *-----------------------------------------------------------------------------------------"
    echo
fi

cleanup

echo "✅ mayohttp installed successfully!"
echo "✨ Run 'mayohttp' to get started."
