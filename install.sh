#!/bin/bash
set -e

# --- Dependency Checks ---
if ! command -v stow >/dev/null 2>&1; then
	echo "Error: stow is not installed. Please install it first." >&2
	exit 1
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
