#!/bin/bash
set -e

# Install Pi Coding Agent if not present
if ! command -v pi >/dev/null 2>&1; then
    echo "Installing Pi Coding Agent..."
    curl -fsSL https://pi.dev/install.sh | sh
else
    echo "Pi Coding Agent is already installed."
fi
