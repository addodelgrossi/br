#!/usr/bin/env sh
set -eu

repo="${BR_REPO:-addodelgrossi/br}"
version="${BR_VERSION:-latest}"
install_dir="${BR_INSTALL_DIR:-$HOME/.local/bin}"
bin_name="br"

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "br: '$1' is required" >&2
    exit 1
  fi
}

need curl
need tar

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  linux|darwin) ;;
  *)
    echo "br: unsupported OS: $os" >&2
    exit 1
    ;;
esac

arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)
    echo "br: unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

asset="br_${os}_${arch}.tar.gz"
if [ "$version" = "latest" ]; then
  url="https://github.com/${repo}/releases/latest/download/${asset}"
else
  url="https://github.com/${repo}/releases/download/${version}/${asset}"
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT INT TERM

mkdir -p "$install_dir"

echo "Downloading $url"
curl -fsSL "$url" -o "$tmp_dir/$asset"
tar -xzf "$tmp_dir/$asset" -C "$tmp_dir"

if [ ! -f "$tmp_dir/$bin_name" ]; then
  echo "br: archive did not contain '$bin_name'" >&2
  exit 1
fi

install -m 0755 "$tmp_dir/$bin_name" "$install_dir/$bin_name"

echo "Installed br to $install_dir/$bin_name"
case ":$PATH:" in
  *":$install_dir:"*) ;;
  *)
    echo "Add this to your shell profile if needed:"
    echo "  export PATH=\"$install_dir:\$PATH\""
    ;;
esac
