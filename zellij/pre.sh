#!/bin/bash
set -e

ZELLIJ_VERSION="v0.43.1"

# --- Zellij Installation ---
installed_version=""
if command -v zellij >/dev/null 2>&1; then
	installed_version="v$(zellij --version | awk '{print $2}')"
fi

if [ "$installed_version" != "$ZELLIJ_VERSION" ]; then
	target=$(uname -sm)
	zellij_root="https://github.com/zellij-org/zellij/releases/download/${ZELLIJ_VERSION}"
	if [ "$target" = "Linux x86_64" ]; then
		download_url="${zellij_root}/zellij-x86_64-unknown-linux-musl.tar.gz"
	elif [ "$target" = "Linux aarch64" ]; then
		download_url="${zellij_root}/zellij-aarch64-unknown-linux-musl.tar.gz"
	elif [ "$target" = "Darwin arm64" ]; then
		download_url="${zellij_root}/zellij-aarch64-apple-darwin.tar.gz"
	else
		echo "Unsupported architecture: $target" >&2
		exit 1
	fi

	echo "Installing zellij ${ZELLIJ_VERSION} from ${download_url}"
	mkdir -p $HOME/bin && curl -L "${download_url}" | sudo tar -xz -C $HOME/bin zellij
fi
