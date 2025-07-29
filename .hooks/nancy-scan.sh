#!/bin/bash
set -e

# Check if nancy is installed
if ! command -v nancy &> /dev/null; then
    echo "Nancy is not installed. Installing..."

    # Detect OS
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')

    case "$OS" in
        darwin)
            # macOS - use Homebrew
            if ! command -v brew &> /dev/null; then
                echo "Error: Homebrew is not installed. Please install Homebrew first:"
                echo "  /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
                exit 1
            fi

            echo "Installing nancy via Homebrew..."
            brew install sonatype-nexus-community/nancy-tap/nancy
            ;;

        linux)
            # Linux - download binary
            ARCH=$(uname -m)

            case "$ARCH" in
                x86_64) ARCH="amd64" ;;
                aarch64 | arm64) ARCH="arm64" ;;
                *)
                    echo "Unsupported architecture: $ARCH"
                                                           exit 1
                                                                  ;;
            esac

            # Download and install nancy
            NANCY_VERSION="v1.0.51"
            NANCY_URL="https://github.com/sonatype-nexus-community/nancy/releases/download/${NANCY_VERSION}/nancy-${NANCY_VERSION}-${OS}-${ARCH}"

            echo "Downloading nancy from ${NANCY_URL}..."
            curl -L "${NANCY_URL}" -o /tmp/nancy
            chmod +x /tmp/nancy

            # Move to a location in PATH (requires sudo)
            echo "Installing nancy to /usr/local/bin (may require sudo)..."
            sudo mv /tmp/nancy /usr/local/bin/nancy
            ;;

        *)
            echo "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
fi

# Run nancy vulnerability scan
echo "Running nancy vulnerability scan..."
go list -json -deps ./... | nancy sleuth

# Capture exit code
exit_code=$?

if [ $exit_code -ne 0 ]; then
    echo ""
    echo "❌ Nancy found vulnerabilities in dependencies!"
    echo "Please fix the vulnerabilities before committing."
    echo ""
    echo "To update vulnerable dependencies, run:"
    echo "  go get -u <package>@<fixed-version>"
    echo "  go mod tidy"
    exit 1
fi

echo "✅ No vulnerabilities found by nancy"
