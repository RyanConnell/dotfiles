#!/bin/bash
set -e

plugin_directory="$HOME/.config/zellij/plugins"
plugins=(
    "https://github.com/dj95/zjstatus/releases/download/v0.23.0/zjstatus.wasm"
)

mkdir -p "$plugin_directory"
for url in "${plugins[@]}"; do
    name=$(echo $url | sed 's/.*\///g')
    curl -sSL -o "$plugin_directory/$name" "$url"
done
