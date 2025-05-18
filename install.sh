#!/bin/bash
# One-liner install script for qik

# --- Configuration ---
REPO="Andr0id88/qik"
BINARY_NAME="qik"             # The name you want the installed binary to have
INSTALL_DIR_CANDIDATES=(
    "$HOME/.local/bin"        # Common on Linux, good practice
    "$HOME/bin"               # Older common practice
    "/usr/local/bin"          # Common system-wide, might need sudo
)
# --- End Configuration ---

set -e # Exit immediately if a command exits with a non-zero status.
set -o pipefail # Ensure that a pipeline command is treated as failed if any of its components fail.

# Function to detect OS and Architecture
get_os_arch() {
    local os_name arch_name
    case "$(uname -s)" in
        Linux*)     os_name=linux;;
        Darwin*)    os_name=darwin;;
        *)          echo "Error: Unsupported Operating System: $(uname -s)" >&2; exit 1;;
    esac

    case "$(uname -m)" in
        x86_64)     arch_name=amd64;;
        aarch64)    arch_name=arm64;; # Linux ARM64
        arm64)      arch_name=arm64;; # macOS Apple Silicon
        *)          echo "Error: Unsupported Architecture: $(uname -m)" >&2; exit 1;;
    esac
    echo "${os_name}-${arch_name}"
}

# Function to get the latest release tag from GitHub API
# Fetches the tag name of the latest (non-draft, non-prerelease) release.
get_latest_release_tag() {
    local repo_path="$1"
    local latest_tag
    # Try with jq first for robustness
    if command -v jq >/dev/null 2>&1; then
        latest_tag=$(curl --silent --fail -H "Accept: application/vnd.github.v3+json" \
            "https://api.github.com/repos/${repo_path}/releases/latest" | jq -r .tag_name)
    else
        # Fallback if jq is not installed (less robust, relies on GitHub API response order and format)
        # This grep might need adjustment if GitHub API changes significantly.
        echo "Warning: jq is not installed. Attempting to fetch latest release tag with grep (less reliable)." >&2
        latest_tag=$(curl --silent --fail -H "Accept: application/vnd.github.v3+json" \
            "https://api.github.com/repos/${repo_path}/releases" | grep -oP '"tag_name":\s*"\Kv[0-9]+\.[0-9]+\.[0-9]+[^"]*' | head -n 1)
    fi

    if [ -z "$latest_tag" ] || [ "$latest_tag" = "null" ]; then # JQ might return "null" as a string
        echo "Error: Could not fetch the latest release tag for ${repo_path}." >&2
        echo "Please check the repository URL, your network connection, or if releases exist." >&2
        exit 1
    fi
    echo "$latest_tag"
}

# --- Main Script ---
echo "Starting qik installation script..."

OS_ARCH_RAW=$(get_os_arch) # e.g., linux-amd64
OS_NAME=$(echo "$OS_ARCH_RAW" | cut -d'-' -f1)
ARCH_NAME=$(echo "$OS_ARCH_RAW" | cut -d'-' -f2)

echo "Detected System: OS=${OS_NAME}, Arch=${ARCH_NAME}"

LATEST_TAG=$(get_latest_release_tag "$REPO")
if [ -z "$LATEST_TAG" ]; then exit 1; fi # Exit if get_latest_release_tag failed
echo "Latest qik release tag: $LATEST_TAG"

# Extract version number from tag (e.g., v1.2.3 -> 1.2.3)
VERSION_NUMBER="${LATEST_TAG#v}"

# Construct the filename based on the release.yml artifact naming convention
# Example: qik_1.0.0_linux_amd64
DOWNLOAD_FILENAME="${BINARY_NAME}_${VERSION_NUMBER}_${OS_NAME}_${ARCH_NAME}"
if [ "$OS_NAME" = "windows" ]; then # Script primarily targets Unix-like, but handle .exe
    DOWNLOAD_FILENAME="${DOWNLOAD_FILENAME}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${DOWNLOAD_FILENAME}"

echo "Downloading $BINARY_NAME version $LATEST_TAG ($DOWNLOAD_FILENAME) from $DOWNLOAD_URL"

# Create a temporary directory for download and ensure it's cleaned up
TMP_DIR=$(mktemp -d -t qik_install.XXXXXX)
# shellcheck disable=SC2064 # $TMP_DIR is expanded now, not when trap runs
trap "echo 'Cleaning up temporary directory: $TMP_DIR'; rm -rf -- '$TMP_DIR'" EXIT HUP INT QUIT TERM

# Download the binary using curl, -L to follow redirects, -o to specify output file
if ! curl -sfLo "${TMP_DIR}/${BINARY_NAME}" "$DOWNLOAD_URL"; then
    echo "Error: Download failed for $DOWNLOAD_URL." >&2
    echo "Please check the release assets for tag $LATEST_TAG or your network connection." >&2
    exit 1
fi
echo "Download complete."

echo "Making the binary executable..."
chmod +x "${TMP_DIR}/${BINARY_NAME}"

# Find a suitable and writable installation directory
INSTALL_DIR=""
for dir_candidate in "${INSTALL_DIR_CANDIDATES[@]}"; do
    resolved_dir_candidate=$(eval echo "$dir_candidate") # Expand $HOME
    if [ -d "$resolved_dir_candidate" ] && [ -w "$resolved_dir_candidate" ]; then
        INSTALL_DIR="$resolved_dir_candidate"
        echo "Found writable install directory: $INSTALL_DIR"
        break
    elif [ ! -d "$resolved_dir_candidate" ]; then
        # Try to create directory if it's within user's home and doesn't exist
        if [[ "$resolved_dir_candidate" == "$HOME/"* ]]; then
            echo "Directory $resolved_dir_candidate does not exist. Attempting to create..."
            if mkdir -p "$resolved_dir_candidate"; then
                echo "Successfully created $resolved_dir_candidate."
                if [ -w "$resolved_dir_candidate" ]; then
                    INSTALL_DIR="$resolved_dir_candidate"
                    echo "Using newly created writable directory: $INSTALL_DIR"
                    break
                else
                    echo "Warning: Created $resolved_dir_candidate, but it's not writable." >&2
                fi
            else
                echo "Warning: Could not create $resolved_dir_candidate." >&2
            fi
        fi
    fi
done

# Fallback: attempt installation to /usr/local/bin using sudo if no user-writable path was found
if [ -z "$INSTALL_DIR" ]; then
    echo "No user-writable installation directory found in preferred locations."
    usr_local_bin="/usr/local/bin" # Define for clarity
    if [ -d "$usr_local_bin" ]; then
        if command -v sudo >/dev/null 2>&1; then
            echo "Attempting to install to $usr_local_bin using sudo. This may require your password."
            # Check if we can write with sudo (mkdir -p won't fail if it exists, test -w checks writability)
            if sudo mkdir -p "$usr_local_bin" && sudo test -w "$usr_local_bin"; then
                echo "Installing to $usr_local_bin..."
                if sudo mv "${TMP_DIR}/${BINARY_NAME}" "${usr_local_bin}/${BINARY_NAME}"; then
                    echo ""
                    echo "$BINARY_NAME installed successfully to ${usr_local_bin}/${BINARY_NAME}!"
                    echo "You may need to open a new terminal or run 'hash -r' (bash) or 'rehash' (zsh) for the shell to find it."
                    exit 0 # Successfully installed with sudo
                else
                    echo "Error: Failed to move binary to $usr_local_bin with sudo." >&2
                fi
            else
                echo "Error: Cannot write to $usr_local_bin even with sudo, or directory doesn't exist and cannot be created." >&2
            fi
        else
            echo "Error: sudo command not found. Cannot attempt installation to $usr_local_bin." >&2
        fi
    else
        echo "Warning: $usr_local_bin does not exist. Skipping sudo attempt for this location." >&2
    fi
fi

# Final check if an installation directory was determined
if [ -z "$INSTALL_DIR" ]; then
    echo "Error: Could not find or create a suitable writable installation directory from the candidates: ${INSTALL_DIR_CANDIDATES[*]}" >&2
    echo "Please create one of these directories or ensure it's writable, then try again." >&2
    echo "Alternatively, the downloaded binary is at '${TMP_DIR}/${BINARY_NAME}'. You can move it manually to a directory in your PATH." >&2
    # Do not exit 1 here if TMP_DIR is shown, let the trap clean it up.
    # Or, exit 1 and make the trap message more prominent.
    exit 1
fi

echo "Installing $BINARY_NAME to ${INSTALL_DIR}/${BINARY_NAME}..."
if mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"; then
    # Check if the installation directory is in PATH
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        echo ""
        echo "Warning: Installation directory '${INSTALL_DIR}' does not seem to be in your PATH." >&2
        echo "To run '${BINARY_NAME}' directly, you may need to add '${INSTALL_DIR}' to your PATH." >&2
        echo "Example for bash/zsh (add to ~/.bashrc or ~/.zshrc):" >&2
        echo "  export PATH=\"${INSTALL_DIR}:\$PATH\"" >&2
        echo "For fish shell (add to ~/.config/fish/config.fish):" >&2
        echo "  fish_add_path ${INSTALL_DIR}" >&2
        echo "After adding, open a new terminal or source your shell configuration file." >&2
    fi
    echo ""
    echo "$BINARY_NAME installed successfully to ${INSTALL_DIR}/${BINARY_NAME}!"
    echo "You may need to open a new terminal or run 'hash -r' (bash) or 'rehash' (zsh) for the shell to find it."
else
    echo "Error: Failed to move binary to ${INSTALL_DIR}/${BINARY_NAME}." >&2
    echo "The downloaded binary is at '${TMP_DIR}/${BINARY_NAME}'."
    exit 1
fi

echo "Installation complete!"
