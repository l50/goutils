#!/bin/bash
set -ex

pkg=$(go list ./...)
for dir in */; do
    if [[ "${dir}" != ".mage" ]] \
                              && [[ "${dir}" != ".hooks/" ]] \
                              && [[ "${dir}" != "config/" ]] \
                              && [[ "${dir}" != "magefiles/" ]] \
                              && [[ "${dir}" != "cmd/" ]] \
                              && [[ "${dir}" != "bin/" ]] \
                              && [[ "${dir}" != "images/" ]] \
                              && [[ "${dir}" != "resources/" ]] \
                              && [[ "${dir}" != "files/" ]] \
                              && [[ "${dir}" != "logs/" ]]; then
        go vet "${pkg}/${dir}"
    fi
done
