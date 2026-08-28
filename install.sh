#!/usr/bin/env bash

set -e

# ============================================================
# ADBOT INSTALLER
# ============================================================

REPO_URL="https://github.com/codewithsorabh/adbot"

DOWNLOAD_URL="${REPO_URL}/raw/refs/heads/main/installer"

BIN_DIR="${HOME}/.local/bin"
INSTALLER="${BIN_DIR}/installer"


# ============================================================
# Check curl
# ============================================================

if ! command -v curl >/dev/null 2>&1; then
    echo "[ERROR] curl is not installed."
    echo
    echo "Install it with:"
    echo
    echo "  sudo apt update"
    echo "  sudo apt install -y curl"
    echo
    exit 1
fi

echo "install..."

mkdir -p "$BIN_DIR"

# ============================================================
# Download installer
# ============================================================

curl -fL \
    --progress-bar \
    "$DOWNLOAD_URL" \
    -o "$INSTALLER"


# ============================================================
# Make executable
# ============================================================


chmod +x "$INSTALLER"

# ============================================================
# Configure PATH
# ============================================================

PATH_LINE='export PATH="$HOME/.local/bin:$PATH"'

add_path() {

    local file="$1"

    if [ -f "$file" ]; then

        if ! grep -Fqx "$PATH_LINE" "$file" 2>/dev/null; then

            echo "" >> "$file"
            echo "# ADBOT" >> "$file"
            echo "$PATH_LINE" >> "$file"

        fi
    fi
}

# Bash
add_path "$HOME/.bashrc"

# Profile
add_path "$HOME/.profile"

# Zsh
add_path "$HOME/.zshrc"

# Current shell
export PATH="$BIN_DIR:$PATH"

# ============================================================
# Verify installer
# ============================================================

if [ ! -f "$INSTALLER" ]; then
    echo
    echo "[ERROR] Installer was not installed."
    exit 1
fi

if [ ! -x "$INSTALLER" ]; then
    echo
    echo "[ERROR] Installer is not executable."
    exit 1
fi

# ============================================================
# Complete
# ============================================================

echo
echo "=============================================="
echo "        ADBOT INSTALLATION COMPLETE"
echo "=============================================="
echo

echo "Installed successully"

echo
echo "Run:"
echo "  installer"