#!/usr/bin/env bash
set -e

VERSION="v0.1.0"
REPO="rasros/lx"

uname_os=$(uname -s)
uname_arch=$(uname -m)

case "$uname_os" in
  Linux*)  os="linux" ;;
  Darwin*) os="darwin" ;;
  *) echo "Unsupported OS: $uname_os"; exit 1 ;;
esac

case "$uname_arch" in
  x86_64)  arch="amd64" ;;
  aarch64) arch="arm64" ;;
  arm64)   arch="arm64" ;;
  *) echo "Unsupported architecture: $uname_arch"; exit 1 ;;
esac

file="lx-${os}-${arch}"
url="https://github.com/${REPO}/releases/download/${VERSION}/${file}"

echo "Downloading $url..."
curl -L "$url" -o lx

chmod +x lx
mkdir -p "$HOME/.local/bin"
mv lx "$HOME/.local/bin/lx"

echo "Installed to $HOME/.local/bin/lx"
echo "Make sure it's in your PATH."

