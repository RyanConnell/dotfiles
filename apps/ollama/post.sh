#!/bin/bash
set -e

# Ensure Ollama service is running and enabled on boot
if command -v systemctl >/dev/null 2>&1; then
    echo "Enabling and starting Ollama service..."
    sudo systemctl enable ollama || true
    sudo systemctl start ollama || true
fi

# Wait for Ollama to become responsive
echo "Waiting for Ollama to initialize..."
attempts=0
max_attempts=30
until curl -s http://localhost:11434/ >/dev/null; do
    attempts=$((attempts + 1))
    if [ "$attempts" -ge "$max_attempts" ]; then
        echo "Error: Ollama did not initialize within 30 seconds." >&2
        exit 1
    fi
    sleep 1
done

echo "Pulling local LLM models (this may take a few minutes)..."
ollama pull gemma4:26b || echo "Warning: Failed to pull gemma4:26b. Please pull manually later."
ollama pull qwen3.5:35b || echo "Warning: Failed to pull qwen3.5:35b. Please pull manually later."

# Customise gemma4:26b with a 64k context window
if ! ollama list | grep -q "gemma4-64k"; then
    echo "Creating gemma4-64k model in Ollama..."
    TEMP_DIR=$(mktemp -d)
    cat <<EOF > "$TEMP_DIR/Modelfile"
FROM gemma4:26b
PARAMETER num_ctx 65536
EOF
    ollama create gemma4-64k -f "$TEMP_DIR/Modelfile"
    rm -rf "$TEMP_DIR"
fi
