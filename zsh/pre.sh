#!/bin/bash
set -e

# --- Zsh Installation ---
if ! command -v zsh >/dev/null 2>&1; then
    echo "zsh is not installed. Attempting to install..."
    if command -v apt-get >/dev/null 2>&1; then
        sudo apt-get update && sudo apt-get install -y zsh
    elif command -v dnf >/dev/null 2>&1; then
        sudo dnf install -y zsh
    elif command -v pacman >/dev/null 2>&1; then
        sudo pacman -Syu --noconfirm zsh
    elif command -v brew >/dev/null 2>&1; then
        brew install zsh
    else
        echo "Could not find a supported package manager. Please install zsh manually." >&2
        exit 1
    fi
fi

# --- Oh My Zsh Installation ---
if [ ! -d "$HOME/.oh-my-zsh" ]; then
    echo "Installing Oh My Zsh..."
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended
fi
