#!/bin/bash
set -e

# --- Tmux Installation ---
if ! command -v tmux >/dev/null 2>&1; then
    echo "tmux is not installed. Attempting to install..."
    if command -v apt-get >/dev/null 2>&1; then
        sudo apt-get update && sudo apt-get install -y tmux
    elif command -v dnf >/dev/null 2>&1; then
        sudo dnf install -y tmux
    elif command -v pacman >/dev/null 2>&1; then
        sudo pacman -Syu --noconfirm tmux
    elif command -v brew >/dev/null 2>&1; then
        brew install tmux
    else
        echo "Could not find a supported package manager. Please install tmux manually." >&2
        exit 1
    fi
fi
