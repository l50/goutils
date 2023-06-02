#!/bin/bash
set -ex

pkgs=$(go list ./...)

for pkg in $pkgs; do
    dir="$(basename "$pkg")/"
    if [[ "${dir}" != ".hooks/" ]] && [[ "${dir}" != "magefiles/" ]]; then
        go vet "${pkg}/${dir}"
    fi
done
