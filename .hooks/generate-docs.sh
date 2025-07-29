#!/bin/bash
#
# This script is a pre-commit hook that checks if the mage command is
# installed and if not, prompts the user to install it. If mage is
# installed, the script navigates to the repository root by calling the
# `rr` function and runs the `mage generatepackagedocs` command.
# This command generates documentation for all Go packages in the current
# directory and its subdirectories by traversing the file tree and creating
# a new README.md file in each directory containing a Go package.
# The script also ensures the presence of a specific utility bash file
# (bashutils.sh) and sources it if found. If the command fails, the commit
# is stopped, and an error message is shown.
set -e

# Define the URL of bashutils.sh
bashutils_url="https://raw.githubusercontent.com/l50/dotfiles/refs/heads/main/bashutils.sh"

# Define the local path of bashutils.sh
bashutils_path="/tmp/bashutils"

# Remove existing file if it exists to ensure fresh download
rm -f "${bashutils_path}"

# Download with error checking
echo "Downloading bashutils.sh..."
if ! curl -fsSL "${bashutils_url}" -o "${bashutils_path}"; then
    curl_exit_code=$?
    echo "Failed to download bashutils.sh from ${bashutils_url}"
    echo "HTTP response code: ${curl_exit_code}"
    exit 1
fi

# Verify the file starts with a shebang
if ! head -n 1 "${bashutils_path}" | grep -q "^#!"; then
    echo "Downloaded file does not appear to be a valid script"
    echo "First line of file:"
    head -n 1 "${bashutils_path}"
    exit 1
fi

# Make it executable
chmod +x "${bashutils_path}"

# Source bashutils
# shellcheck source=/dev/null
source "${bashutils_path}"

rr || exit 1

# Check if mage is installed
if ! command -v mage > /dev/null 2>&1; then
    echo -e "mage is not installed\n"
    echo -e "Please install mage by running the following command:\n"
    echo -e "go install github.com/magefile/mage@latest\n"
    exit 1
fi

# Run the mage generatepackagedocs command
mage generatepackagedocs
# Catch the exit code
exit_status=$?

# If the exit code is not zero (i.e., the command failed),
# then stop the commit and show an error message
if [ $exit_status -ne 0 ]; then
    echo "failed to generate package docs"
    exit 1
fi
