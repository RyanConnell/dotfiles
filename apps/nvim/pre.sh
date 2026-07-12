#!/bin/bash
set -e

# --- Neovim Installation ---
if ! command -v nvim >/dev/null 2>&1; then
    echo "neovim is not installed. Attempting to install..."
    if command -v apt-get >/dev/null 2>&1; then
        if command -v add-apt-repository; then
            sudo add-apt-repository ppa:neovim-ppa/stable -y
        fi
        sudo apt-get update
        sudo apt-get install -y neovim
    elif command -v dnf >/dev/null 2>&1; then
        sudo dnf install -y neovim
    elif command -v pacman >/dev/null 2>&1; then
        sudo pacman -Syu --noconfirm neovim
    elif command -v brew >/dev/null 2>&1; then
        brew install neovim
    else
        echo "Could not find a supported package manager. Please install neovim manually." >&2
        exit 1
    fi
fi
