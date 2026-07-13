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

echo "Pulling local LLM models (this may take a while)..."
ollama pull gemma4:26b || echo "Warning: Failed to pull gemma4:26b"
ollama pull qwen3.6:27b-q4_K_M || echo "Warning: Failed to pull qwen3.6:27b-q4_K_M"

# Customise gemma4:26b with a 64k context window and 4k predict limit
if ! ollama list | grep -q "gemma4-64k"; then
    echo "Creating gemma4-64k model in Ollama..."
    TEMP_DIR=$(mktemp -d)
    cat <<EOF > "$TEMP_DIR/Modelfile"
FROM gemma4:26b
PARAMETER num_ctx 65536
PARAMETER num_predict 4096
EOF
    ollama create gemma4-64k -f "$TEMP_DIR/Modelfile"
    rm -rf "$TEMP_DIR"
fi

# Customise qwen3.6:27b-q4_K_M with a 128k context window and 4k predict limit
if ! ollama show --modelfile qwen3.6:27b-q4_K_M 2>/dev/null | grep -q "num_ctx 131072"; then
    echo "Creating qwen3.6:27b-q4_K_M custom model with 128k context window in Ollama..."
    TEMP_DIR=$(mktemp -d)
    cat <<EOF > "$TEMP_DIR/Modelfile"
FROM qwen3.6:27b-q4_K_M
PARAMETER num_ctx 131072
PARAMETER num_predict 4096
EOF
    ollama create qwen3.6:27b-q4_K_M -f "$TEMP_DIR/Modelfile"
    rm -rf "$TEMP_DIR"
fi
