#!/usr/bin/env bash
set -euo pipefail

REPO="rasros/lx"

# Allow override, e.g. LX_INSTALL_DIR=/usr/local/bin ./install.sh
INSTALL_DIR="${LX_INSTALL_DIR:-"$HOME/.local/bin"}"

log() {
  echo "[lx install] $*" >&2
}

detect_os() {
  local uname_os
  uname_os="$(uname -s)"

  case "$uname_os" in
    Linux*)  echo "linux" ;;
    Darwin*) echo "darwin" ;;
    MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
    *)
      log "Unsupported OS: $uname_os"
      exit 1
      ;;
  esac
}

detect_arch() {
  local uname_arch
  uname_arch="$(uname -m)"

  case "$uname_arch" in
    x86_64|amd64) echo "amd64" ;;
    aarch64|arm64) echo "arm64" ;;
    *)
      log "Unsupported architecture: $uname_arch"
      exit 1
      ;;
  esac
}

ensure_curl() {
  if ! command -v curl >/dev/null 2>&1; then
    log "curl is required but was not found"
    exit 1
  fi
}

download() {
  local url="$1"
  local out="$2"

  log "Downloading $url"
  curl -fL "$url" -o "$out"
}

ensure_install_dir() {
  if [ ! -d "$INSTALL_DIR" ]; then
    log "Creating install directory: $INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"
  fi
}

check_path() {
  case ":$PATH:" in
    *":$INSTALL_DIR:"*) return 0 ;;
    *)
      log "Warning: $INSTALL_DIR is not in your PATH."
      log "Add this to your shell config, for example:"
      log "  export PATH=\"$INSTALL_DIR:\$PATH\""
      ;;
  esac
}

main() {
  if [ "${1:-}" = "-h" ] || [ "${1:-}" = "--help" ]; then
    cat <<EOF
Install lx from the latest GitHub release.

Usage:
  curl -fsSL https://raw.githubusercontent.com/rasros/lx/main/install.sh | bash

Options:
  LX_INSTALL_DIR   Override install directory (default: \$HOME/.local/bin)
EOF
    exit 0
  fi

  ensure_curl

  os="$(detect_os)"
  arch="$(detect_arch)"

  file="lx-${os}-${arch}"
  if [ "$os" = "windows" ]; then
    file="${file}.exe"
  fi

  url="https://github.com/${REPO}/releases/latest/download/${file}"

  tmpdir="$(mktemp -d)"
  trap 'rm -rf "$tmpdir"' EXIT

  tmpfile="${tmpdir}/lx"
  download "$url" "$tmpfile"

  chmod +x "$tmpfile"
  ensure_install_dir
  mv "$tmpfile" "$INSTALL_DIR/lx"

  check_path

  log "Installed lx to: $INSTALL_DIR/lx"
  log "Run: lx --help"
}

main "$@"
