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

ensure_unpack_tool() {
  local os="$1"
  if [ "$os" = "windows" ]; then
    if ! command -v unzip >/dev/null 2>&1; then
      log "unzip is required to extract zip archives"
      exit 1
    fi
  else
    if ! command -v tar >/dev/null 2>&1; then
      log "tar is required to extract tar.gz archives"
      exit 1
    fi
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

latest_tag() {
  # Uses GitHub API to get latest release tag, e.g. "v1.0.1"
  local api="https://api.github.com/repos/${REPO}/releases/latest"
  local tag
  tag="$(curl -fsSL "$api" | grep -m 1 '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
  if [ -z "$tag" ]; then
    log "Failed to determine latest release tag from GitHub API"
    exit 1
  fi
  echo "$tag"
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
  ensure_unpack_tool "$os"

  tag="$(latest_tag)"        # e.g. "v1.0.1"
  version="${tag#v}"         # e.g. "1.0.1"

  # Archive naming must match your workflow:
  #   lx_${VERSION}_${GOOS}_${GOARCH}.tar.gz (or .zip on windows)
  archive_base="lx_${version}_${os}_${arch}"
  if [ "$os" = "windows" ]; then
    archive_file="${archive_base}.zip"
  else
    archive_file="${archive_base}.tar.gz"
  fi

  url="https://github.com/${REPO}/releases/download/${tag}/${archive_file}"

  tmpdir="$(mktemp -d)"
  trap 'rm -rf "$tmpdir"' EXIT

  tmp_archive="${tmpdir}/${archive_file}"
  download "$url" "$tmp_archive"

  # Extract archive
  if [ "$os" = "windows" ]; then
    (cd "$tmpdir" && unzip -q "$tmp_archive")
  else
    (cd "$tmpdir" && tar -xzf "$tmp_archive")
  fi

  # Binary name inside archive must match your build step:
  #   lx-${VERSION}-${GOOS}-${GOARCH}[.exe]
  bin_name="lx-${version}-${os}-${arch}"
  if [ "$os" = "windows" ]; then
    bin_name="${bin_name}.exe"
  fi

  src_bin="${tmpdir}/${bin_name}"
  if [ ! -f "$src_bin" ]; then
    log "Binary not found in archive: $bin_name"
    exit 1
  fi

  ensure_install_dir
  chmod +x "$src_bin"
  mv "$src_bin" "$INSTALL_DIR/lx"

  check_path

  log "Installed lx ${version} to: $INSTALL_DIR/lx"
  log "Run: lx --help"
}

main "$@"

