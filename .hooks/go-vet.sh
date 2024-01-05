#!/bin/bash
set -e

pkgs=$(go list ./...)

for pkg in $pkgs; do
    dir="$(basename "$pkg")/"
    if [[ "${dir}" != .*/ ]] && [[ "${dir}" != "magefiles/" ]]; then
        go vet "${pkg}"
    fi
done
