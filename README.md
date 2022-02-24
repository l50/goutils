# goutils

[![Go Report Card](https://goreportcard.com/badge/github.com/l50/goutils)](https://goreportcard.com/badge/github.com/l50/goutils)
[![License](http://img.shields.io/:license-mit-blue.svg)](https://github.com/l50/goutils/blob/master/LICENSE)
[![Lint test for code](https://github.com/l50/goutils/actions/workflows/lint.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/lint.yaml)
[![🚨 CodeQL Analysis](https://github.com/l50/goutils/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/codeql-analysis.yaml)
[![🚨 Semgrep Analysis](https://github.com/l50/goutils/actions/workflows/semgrep.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/semgrep.yaml)
[![Run Pre-Commit hooks](https://github.com/l50/goutils/actions/workflows/pre-commit.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/pre-commit.yaml)

This repo is comprised of utilities that I use across various go projects.

## Dependencies

- [Install golang](https://go.dev/):

  ```bash
  gvm install go1.16.4
  ```

- [Install pre-commit](https://pre-commit.com/):

  ```bash
  brew install pre-commit
  ```

- [Install Mage](https://magefile.org/):

  ```bash
  go install github.com/magefile/mage@latest
  ```

- Set up `go.mod` for development:

  ```bash
  REPO=github.com/l50/goutils
  FORK="${PWD}"

  echo -e "\nreplace ${REPO} => ${FORK}" >> go.mod
  ```

- [Optional - install gvm](https://github.com/moovweb/gvm):

  ```bash
  bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
  source $GVM_BIN
  ```

---

## Developer Environment Setup

1. [Fork this project](https://docs.github.com/en/get-started/quickstart/fork-a-repo)

2. (Optional) If you installed gvm, create golang pkgset specifically for this project:

   ```bash
   gvm pkgset create ws
   gvm pkgset use ws
   ```

3. Install pre-commit dependencies:

   ```bash
   go install golang.org/x/lint/golint@latest
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   go install github.com/go-critic/go-critic/cmd/gocritic@latest
   ```

4. Install pre-commit hooks locally

   ```bash
   mage installPreCommit
   ```

5. Update and run pre-commit hooks locally:

   ```bash
    mage runPreCommit
   ```
