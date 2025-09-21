#!/bin/bash
set -e

# --- Zellij Installation ---
if ! command -v zellij >/dev/null 2>&1; then
	target=$(uname -sm)
	latest_version=$(curl -s "https://api.github.com/repos/zellij-org/zellij/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")')
	zellij_root="https://github.com/zellij-org/zellij/releases/download/${latest_version}"
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

	echo "Installing zellij from ${download_url}"
	curl -L "${download_url}" | sudo tar -xz -C $HOME/bin zellij
fi
