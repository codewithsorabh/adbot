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

echo "Installing ADBOT..."
echo

# ============================================================
# Create directories
# ============================================================

mkdir -p "$BIN_DIR"

# ============================================================
# Download installer
# ============================================================

echo "[1/4] Downloading installer..."

curl -fL \
    --progress-bar \
    "$DOWNLOAD_URL" \
    -o "$INSTALLER"

# ============================================================
# Make executable
# ============================================================

echo "[2/4] Making installer executable..."

chmod +x "$INSTALLER"

# ============================================================
# Configure PATH
# ============================================================

echo "[3/4] Configuring PATH..."

PATH_LINE='export PATH="$HOME/.local/bin:$PATH"'

add_path() {

    local file="$1"

    # Create file if it does not exist
    if [ ! -f "$file" ]; then
        echo "[INFO] Creating $file..."
        touch "$file"
    fi

    # Add PATH only if it is not already present
    if ! grep -Fqx "$PATH_LINE" "$file" 2>/dev/null; then

        echo "" >> "$file"
        echo "# ADBOT" >> "$file"
        echo "$PATH_LINE" >> "$file"

        echo "[OK] PATH added to $file"
    else
        echo "[OK] PATH already configured in $file"
    fi
}

# Bash
add_path "$HOME/.bashrc"

# Profile
add_path "$HOME/.profile"

# Zsh
add_path "$HOME/.zshrc"

# ============================================================
# Current shell PATH
# ============================================================

export PATH="$BIN_DIR:$PATH"

# ============================================================
# Reload Bash configuration
# ============================================================

if [ -f "$HOME/.bashrc" ]; then
    source "$HOME/.bashrc"
fi
# ============================================================
# Verify installer
# ============================================================

echo "[4/4] Verifying installation..."

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

# Verify command
if command -v installer >/dev/null 2>&1; then
    echo "[OK] installer command is available."
else
    echo "[INFO] installer installed successfully."
    echo "[INFO] Restart your shell or run:"
    echo
    echo "  source ~/.bashrc"
fi

# ============================================================
# Complete
# ============================================================

echo
echo "=============================================="
echo "       ADBOT INSTALLATION COMPLETE"
echo "=============================================="
echo

echo "Installed successfully."

echo
echo "Run:"
echo "  installer"
echo