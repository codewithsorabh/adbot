#!/usr/bin/env bash

# ==============================================================================
# ADBOT Release Installer
# Repository: https://github.com/codewithsorabh/adbot
# ==============================================================================

set -euo pipefail

REPO="codewithsorabh/adbot"
GITHUB_API="https://api.github.com/repos/${REPO}"
GITHUB_URL="https://github.com/${REPO}"
BIN_NAME="adbot"

# Default install directory: ~/.local/bin (or /usr/local/bin if running as root)
if [ "$(id -u)" -eq 0 ]; then
    INSTALL_DIR="/usr/local/bin"
else
    INSTALL_DIR="${HOME}/.local/bin"
fi

# ==============================================================================
# Helpers
# ==============================================================================

info() {
    printf "[INFO] %s\n" "$*"
}

step() {
    printf "==> %s\n" "$*"
}

error() {
    printf "[ERROR] %s\n" "$*" >&2
}

fatal() {
    error "$*"
    exit 1
}

# ==============================================================================
# Check Requirements
# ==============================================================================

command -v curl >/dev/null 2>&1 || fatal "curl is required to run this installer. Please install curl and try again."

# ==============================================================================
# Detect Operating System & Architecture
# ==============================================================================

step "Detecting system architecture..."

OS_RAW="$(uname -s)"
case "${OS_RAW}" in
    Linux*)     OS="linux" ;;
    Darwin*)    OS="darwin" ;;
    *)          fatal "Unsupported operating system: ${OS_RAW}. Prebuilt binaries support Linux and macOS." ;;
esac

ARCH_RAW="$(uname -m)"
case "${ARCH_RAW}" in
    x86_64|amd64)   ARCH="amd64" ;;
    aarch64|arm64)  ARCH="arm64" ;;
    *)              fatal "Unsupported CPU architecture: ${ARCH_RAW}. Prebuilt binaries support amd64 and arm64." ;;
esac

info "Target platform: ${OS}/${ARCH}"

# ==============================================================================
# Determine Release Version
# ==============================================================================

step "Determining latest ADBOT release..."

TAG="${ADBOT_VERSION:-${VERSION:-}}"

if [ -z "${TAG}" ]; then
    # Query GitHub API for latest release
    API_RESPONSE="$(curl -fsSL -H "Accept: application/vnd.github.v3+json" "${GITHUB_API}/releases/latest" 2>/dev/null || true)"
    
    if [ -n "${API_RESPONSE}" ]; then
        TAG="$(printf '%s' "${API_RESPONSE}" | grep '"tag_name":' | sed -E 's/.*"tag_name":\s*"([^"]+)".*/\1/' | head -n 1 || true)"
    fi

    # Fallback to repository VERSION file if API fails or rate-limited
    if [ -z "${TAG}" ]; then
        RAW_VERSION="$(curl -fsSL "https://raw.githubusercontent.com/${REPO}/main/VERSION" 2>/dev/null || true)"
        if [ -n "${RAW_VERSION}" ]; then
            TAG="v$(printf '%s' "${RAW_VERSION}" | tr -d '[:space:]' | sed 's/^v//')"
        fi
    fi
fi

[ -n "${TAG}" ] || fatal "Could not determine the latest ADBOT release version from GitHub."

# Ensure tag has 'v' prefix
case "${TAG}" in
    v*) ;;
    *)  TAG="v${TAG}" ;;
esac

info "Release version: ${TAG}"

# ==============================================================================
# Construct Download URLs
# ==============================================================================

BINARY_FILENAME="adbot-${TAG}-${OS}-${ARCH}"
DOWNLOAD_URL="${GITHUB_URL}/releases/download/${TAG}/${BINARY_FILENAME}"
CHECKSUMS_URL="${GITHUB_URL}/releases/download/${TAG}/checksums.txt"

# Temporary workspace
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT INT TERM

# ==============================================================================
# Download Binary
# ==============================================================================

step "Downloading ${BINARY_FILENAME}..."

TARGET_BINARY="${TMP_DIR}/${BINARY_FILENAME}"
if ! curl -fSL --progress-bar "${DOWNLOAD_URL}" -o "${TARGET_BINARY}"; then
    fatal "Failed to download binary from ${DOWNLOAD_URL}."
fi

# ==============================================================================
# Checksum Verification
# ==============================================================================

step "Verifying binary checksum..."

CHECKSUMS_FILE="${TMP_DIR}/checksums.txt"
if curl -fsSL "${CHECKSUMS_URL}" -o "${CHECKSUMS_FILE}" 2>/dev/null; then
    EXPECTED_HASH="$(grep -E "${BINARY_FILENAME}\$" "${CHECKSUMS_FILE}" | awk '{print $1}' | head -n 1 || true)"

    if [ -n "${EXPECTED_HASH}" ]; then
        if command -v sha256sum >/dev/null 2>&1; then
            ACTUAL_HASH="$(sha256sum "${TARGET_BINARY}" | awk '{print $1}')"
        elif command -v shasum >/dev/null 2>&1; then
            ACTUAL_HASH="$(shasum -a 256 "${TARGET_BINARY}" | awk '{print $1}')"
        else
            ACTUAL_HASH=""
            info "Neither sha256sum nor shasum is available; skipping checksum verification."
        fi

        if [ -n "${ACTUAL_HASH}" ]; then
            if [ "${EXPECTED_HASH}" != "${ACTUAL_HASH}" ]; then
                fatal "Checksum verification failed! Expected: ${EXPECTED_HASH}, got: ${ACTUAL_HASH}"
            fi
            info "Checksum verified successfully (${ACTUAL_HASH:0:16}...)."
        fi
    else
        info "Checksum for ${BINARY_FILENAME} not found in checksums.txt; proceeding."
    fi
else
    info "checksums.txt not available for this release; proceeding."
fi

# ==============================================================================
# Install Binary
# ==============================================================================

step "Installing ${BIN_NAME} to ${INSTALL_DIR}..."

mkdir -p "${INSTALL_DIR}"
chmod 0755 "${TARGET_BINARY}"
mv -f "${TARGET_BINARY}" "${INSTALL_DIR}/${BIN_NAME}"

# ==============================================================================
# Configure PATH
# ==============================================================================

add_path_to_file() {
    local rc_file="$1"
    local path_entry="$2"

    if [ -f "${rc_file}" ] || [ "${rc_file##*/}" = ".bashrc" ]; then
        touch "${rc_file}"
        if ! grep -Fq "${path_entry}" "${rc_file}" 2>/dev/null; then
            printf '\n# ADBOT\nexport PATH="%s:$PATH"\n' "${path_entry}" >> "${rc_file}"
            info "Added ${path_entry} to ${rc_file}"
        fi
    fi
}

PATH_NEEDS_UPDATE=0
case ":${PATH}:" in
    *:"${INSTALL_DIR}":*) ;;
    *) PATH_NEEDS_UPDATE=1 ;;
esac

if [ "${PATH_NEEDS_UPDATE}" -eq 1 ]; then
    step "Configuring PATH..."
    add_path_to_file "${HOME}/.bashrc" "${INSTALL_DIR}"
    add_path_to_file "${HOME}/.profile" "${INSTALL_DIR}"
    add_path_to_file "${HOME}/.zshrc" "${INSTALL_DIR}"
    export PATH="${INSTALL_DIR}:${PATH}"
fi

# ==============================================================================
# Verification
# ==============================================================================

step "Verifying installation..."

if [ ! -x "${INSTALL_DIR}/${BIN_NAME}" ]; then
    fatal "Installation failed: ${INSTALL_DIR}/${BIN_NAME} is missing or not executable."
fi

echo
echo "=============================================="
echo "        ADBOT INSTALLATION COMPLETE           "
echo "=============================================="
echo
info "Installed version: ${TAG}"
info "Binary location:   ${INSTALL_DIR}/${BIN_NAME}"
echo

if [ "${PATH_NEEDS_UPDATE}" -eq 1 ]; then
    echo "To update your current shell session, run:"
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    echo "or reload your shell profile:"
    echo "  source ~/.bashrc"
    echo
fi

echo "To start ADBOT setup, run:"
echo "  ${BIN_NAME}"
echo