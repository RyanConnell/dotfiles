#!/bin/bash
set -e

PACKAGES=$(ls -d */ | sed 's/\///' | xargs)

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

# --- Create backup directory ---
BACKUP_DIR="backups"

# --- Run Package Pre-Install Scripts ---
for pkg in $PACKAGES; do
    if [ -f "${pkg}/pre.sh" ]; then
        /bin/bash "${pkg}/pre.sh" | sed "s/^/[${pkg}\/pre.sh]: /"
    fi
done

# --- Install config/dotfiles for each application ---
TIMESTAMP=$(date +%Y-%m-%d_%H-%M-%S)
for pkg in $PACKAGES; do
    # Attempt to stow
    if ! stow -R "${pkg}" --ignore '(pre|post).sh'; then
        echo "Stow for '${pkg}' failed; Attempting to backup files"
        # Find all files in the package directory, excluding pre/post scripts
        files=$(find "$pkg" -type f -not -path "*/pre.sh" -not -path "*/post.sh")

        # --- Backup any existing files to avoid loss of data ---
        for file in ${files#$pkg/}; do
            target_file="$HOME/$file"
            if [ -f "$target_file" ] && [ ! -L "$target_file" ]; then
                backup_file="$BACKUP_DIR/$TIMESTAMP/$file"
                echo "[${pkg}]: Backing up $target_file to $backup_file"
                mkdir -p "$(dirname $backup_file)"
                mv "$target_file" "$backup_file"
            fi
        done

        # Attempt to stow again now that the files have been moved.
        stow -R "$pkg" --ignore '(pre|post).sh'
    fi
done

# --- Run Package Post-Install Scripts ---
for pkg in $PACKAGES; do
    if [ -f "${pkg}/post.sh" ]; then
        /bin/bash "${pkg}/post.sh" | sed "s/^/[${pkg}\/post.sh]: /"
    fi
done

echo "Dotfiles setup complete!"
