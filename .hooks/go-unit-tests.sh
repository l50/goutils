#!/bin/bash

set -ex

TESTS_TO_RUN=$1
RETURN_CODE=0

TIMESTAMP=$(date +"%Y%m%d%H%M%S")
LOGFILE="/tmp/goutils-unit-test-results-$TIMESTAMP.log"

if [[ -z "${TESTS_TO_RUN}" ]]; then
    echo "No tests input"
    echo "Example - Run all shorter collection of tests: bash go-unit-tests.sh short"
    echo "Example - Run all tests: bash go-unit-tests.sh all"
    echo "Example - Run coverage for a specific version: bash go-unit-tests.sh coverage"
    exit 1
fi

run_tests()
            {
    local coverage_file=$1
    repo_root=$(git rev-parse --show-toplevel 2> /dev/null) || exit
    pushd "${repo_root}" || exit
    echo "Logging output to ${LOGFILE}"
    echo "Run the following command to see the output in real time:"
    echo "tail -f ${LOGFILE}"
    echo "Running tests..."

    if [[ "${TESTS_TO_RUN}" == 'coverage' ]]; then
        go test -v -race -failfast -tags=integration -coverprofile="${coverage_file}" ./... &>> "$LOGFILE"
    elif [[ "${TESTS_TO_RUN}" == 'all' ]]; then
        go test -v -race -failfast ./... &>> "$LOGFILE"
    elif [[ "${TESTS_TO_RUN}" == 'short' ]] && [[ "${GITHUB_ACTIONS}" != "true" ]]; then
        go test -v -short -failfast -race ./... &>> "$LOGFILE"
    else
        if [[ "${GITHUB_ACTIONS}" != 'true' ]]; then
            go test -v -failfast -race "./.../${TESTS_TO_RUN}" &>> "$LOGFILE"
        fi
    fi

    RETURN_CODE=$?
}

if [[ "${TESTS_TO_RUN}" == 'coverage' ]]; then
    run_tests 'coverage-all.out'
else
    run_tests
fi

if [[ "${RETURN_CODE}" -ne 0 ]]; then
    echo "unit tests failed"
    exit 1
fi
