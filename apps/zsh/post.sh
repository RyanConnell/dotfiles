#!/bin/bash
set -e

# --- Change Default Shell ---
if [ "$SHELL" != "/bin/zsh" ] && [ "$SHELL" != "/usr/bin/zsh" ]; then
    read -p "Do you want to change your default shell to zsh? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if command -v zsh >/dev/null 2>&1; then
            chsh -s "$(command -v zsh)"
            echo "Default shell changed to zsh. Please log out and log back in for the changes to take effect."
        else
            echo "zsh is not installed. Cannot change default shell."
        fi
    fi
fi
