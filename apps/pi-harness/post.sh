# -------- Helper Functions -------- #

INSTALLED_APPS=$(pi list)

install_if_missing() {
    if echo "$INSTALLED_APPS" | grep -q "$1"; then
        echo "$1 is already installed."
    else
        echo "Installing $1..."
        pi install "$1"
    fi
}

# -------- Plugins -------- #

# Context plugins
install_if_missing npm:pi-observational-memory
install_if_missing git:github.com/elpapi42/pi-fork

# Workflow plugins
install_if_missing npm:@tintinweb/pi-subagents
install_if_missing npm:@juicesharp/rpiv-ask-user-question
install_if_missing npm:@juicesharp/rpiv-todo

# Other utility plugins
install_if_missing npm:@d3ara1n/pi-session-namer
install_if_missing npm:pi-web-access
install_if_missing npm:pi-vimmode
install_if_missing npm:pi-bash-live-view
