#!/bin/bash
set -ex

RETURN_CODE=0
TIMESTAMP=$(date +"%Y%m%d%H%M%S")
LOGFILE="/tmp/go-vet-results-$TIMESTAMP.log"
MODULE_ROOT=$(go list -m -f "{{.Dir}}")

function run_go_vet() {
    local input_files=("$@")
    local pkg_dirs=()

    for file in "${input_files[@]}"; do
        if [[ -f "$file" && "$file" == *.go ]]; then
            local pkg_dir
            pkg_dir=$(dirname "$file")
            if [[ "$pkg_dir" != "magefiles" ]]; then
                # Convert to package import path relative to module root
                pkg_dir=${pkg_dir#"$MODULE_ROOT/"}
                pkg_dirs+=("$pkg_dir")
            fi
        fi
    done

    # Remove duplicate package directories
    IFS=$'\n' read -r -a pkg_dirs <<< "$(sort -u <<< "${pkg_dirs[*]}")"
    unset IFS

    for dir in "${pkg_dirs[@]}"; do
        go vet "./$dir/..."
    done

    RETURN_CODE=$?
}

# Check if any filenames are passed
if [ "$#" -eq 0 ]; then
    # No filenames passed, run go vet on all packages
    go vet ./...
    RETURN_CODE=$?
else
    # Filenames passed, run go vet on those files' packages
    run_go_vet "$@"
fi

if [[ "${RETURN_CODE}" -ne 0 ]]; then
    echo "go vet failed" | tee -a "$LOGFILE"
    exit 1
fi
