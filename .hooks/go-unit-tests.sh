#!/bin/bash

set -x

TESTS_TO_RUN=$1
RETURN_CODE=0

if [[ -z "${TESTS_TO_RUN}" ]]; then
    echo "No tests input"
    echo "Example - Run all tests for a specific version: bash go-unit-tests.sh all v1"
    echo "Example - Run all tests for both versions: bash go-unit-tests.sh all all"
    echo "Example - Run coverage for a specific version: bash go-unit-tests.sh coverage v2"
    echo "Example - Run all tests for v1 (default if version is not specified): bash go-unit-tests.sh all"
    exit 1
fi

run_tests()
           {
    local coverage_file=$1

    repo_root=$(git rev-parse --show-toplevel 2> /dev/null) || exit
    pushd "${repo_root}" || exit

    if [[ "${TESTS_TO_RUN}" == 'coverage' ]]; then
        go test -v -race -failfast -tags=integration -coverprofile="${coverage_file}" ./...
    elif [[ "${TESTS_TO_RUN}" == 'all' ]]; then
        go test -v -race -failfast ./...
    elif [[ "${TESTS_TO_RUN}" == 'short' ]] && [[ "${GITHUB_ACTIONS}" != "true" ]]; then
        go test -v -short -failfast -race ./...
    else
        if [[ "${GITHUB_ACTIONS}" != 'true' ]]; then
            go test -v -race -failfast "./.../${TESTS_TO_RUN}"
        fi
    fi

    RETURN_CODE=$?
}

if [[ "${TESTS_TO_RUN}" == 'all' ]]; then
    run_tests '.' 'coverage-all.out'
fi

if [[ "${RETURN_CODE}" -ne 0 ]]; then
    echo "unit tests failed"
    exit 1
fi
