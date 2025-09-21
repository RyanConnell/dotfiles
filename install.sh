#!/bin/bash
set -e

# --- Ensure that Stow is installed ---
if ! command -v stow >/dev/null 2>&1; then
    echo "stow is not installed. Attempting to install..."
    if command -v apt-get >/dev/null 2>&1; then
        sudo apt-get update && sudo apt-get install -y stow
    elif command -v pacman >/dev/null 2>&1; then
        sudo pacman -Syu --noconfirm stow
    elif command -v dnf >/dev/null 2>&1; then
        sudo dnf install -y stow
    elif command -v brew >/dev/null 2>&1; then
        brew install stow
    else
        echo "Error: Could not find a supported package manager (apt, pacman, dnf, brew)." >&2
        echo "Please install stow manually and re-run this script." >&2
        exit 1
    fi
fi

PACKAGES=$(ls -d */ | sed 's/\///' | xargs)

# --- Run Package Pre-Install Scripts ---
for pkg in $PACKAGES; do
    if [ -f "${pkg}/pre.sh" ]; then
        /bin/bash "${pkg}/pre.sh" | sed "s/^/[${pkg}\/pre.sh]: /"
    fi
done

# --- Install config/dotfiles for each application ---
for pkg in $PACKAGES; do
	stow -R "$pkg" --ignore '(pre|post).sh'
done

# --- Run Package Post-Install Scripts ---
for pkg in $PACKAGES; do
    if [ -f "${pkg}/post.sh" ]; then
        /bin/bash "${pkg}/post.sh" | sed "s/^/[${pkg}\/post.sh]: /"
    fi
done

echo "Dotfiles setup complete!"
