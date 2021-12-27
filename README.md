# goutils
[![Go Report Card](https://goreportcard.com/badge/github.com/l50/goutils)](https://goreportcard.com/badge/github.com/l50/goutils)
[![License](http://img.shields.io/:license-mit-blue.svg)](https://github.com/l50/goutils/blob/master/LICENSE)
[![Test goutils](https://github.com/l50/goutils/actions/workflows/test-build.yml/badge.svg)](https://github.com/l50/goutils/actions/workflows/test-build.yml)
[![Lint test for code](https://github.com/l50/goutils/actions/workflows/lint.yml/badge.svg)](https://github.com/l50/goutils/actions/workflows/lint.yml)
[![🚨 CodeQL Analysis](https://github.com/l50/goutils/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/l50/goutils/actions/workflows/codeql-analysis.yml)
[![Run Pre-Commit hooks](https://github.com/l50/goutils/actions/workflows/pre-commit.yml/badge.svg)](https://github.com/l50/goutils/actions/workflows/pre-commit.yml)

This repo is comprised of utilities that I use across various go projects.

## Development
The following steps can be followed to prepare your environment to hack on `goutils`:

1. Install Mage:
```
go install github.com/magefile/mage@latest
```

2. Run the following command to set up `go.mod` for development with your fork:
```
REPO=github.com/l50/goutils
FORK="${PWD}"

echo -e "\nreplace ${REPO} => ${FORK}" >> go.mod
```

3. Install Pre-Commit by following these instructions: https://pre-commit.com/#install

4. Set up Pre-Commit hooks locally:
```
mage precommit
```
