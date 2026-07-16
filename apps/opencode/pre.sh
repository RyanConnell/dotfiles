#!/bin/bash
set -e

if ! command -v opencode >/dev/null 2>&1; then
    echo "Installing OpenCode CLI..."
    curl -fsSL https://opencode.ai/install | bash
else
    echo "OpenCode CLI is already installed."
fi
