# -------- Plugins -------- #

# Context plugins
pi install npm:pi-observational-memory
pi install git:github.com/elpapi42/pi-fork

# Workflow plugins
pi install npm:@tintinweb/pi-subagents
pi install npm:@juicesharp/rpiv-ask-user-question
pi install npm:@juicesharp/rpiv-todo

# Other utility plugins
pi install npm:@d3ara1n/pi-session-namer
pi install npm:pi-web-access
pi install npm:pi-vimmode
pi install npm:pi-bash-live-view

# -------- Models -------- #

# Customise Gemma4:26b with a 64k context window
if ! ollama list | grep -q "gemma4-64k"; then
    echo "Creating gemma4-64k model in Ollama..."
    TEMP_DIR=$(mktemp -d)
    cat <<EOF > "$TEMP_DIR/Modelfile"
FROM gemma4:26b
PARAMETER num_ctx 65536
EOF
    ollama create gemma4-64k -f "$TEMP_DIR/Modelfile"
    rm -rf "$TEMP_DIR"
else
    echo "gemma4-64k is already built in Ollama."
fi
