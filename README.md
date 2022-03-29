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
  gvm install go1.18
  ```

- [Install pre-commit](https://pre-commit.com/):

  ```bash
  brew install pre-commit
  ```

- [Install Mage](https://magefile.org/):

  ```bash
  go install github.com/magefile/mage@latest
  ```

- [Optional - install gvm](https://github.com/moovweb/gvm):

  ```bash
  bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
  source "${GVM_BIN}"
  ```

---

## Developer Environment Setup

1. [Fork this project](https://docs.github.com/en/get-started/quickstart/fork-a-repo)

2. (Optional) If you installed gvm, create golang pkgset specifically for this project:

   ```bash
   mkdir "${HOME}/go"
   GVM_BIN="${HOME}/.gvm/scripts/gvm"
   export GOPATH="${HOME}/go"
   VERSION='1.18'
   PROJECT=goutils

   bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
   source $GVM_BIN
   gvm install "go${VERSION}"
   gvm use "go${VERSION}"
   gvm pkgset create "${PROJECT}"
   gvm pkgset use "${PROJECT}"
   ```

3. Generate the `magefile` binary:

   ```bash
   mage -d .mage/ -compile ../magefile
   ```

4. Install pre-commit hooks and dependencies:

   ```bash
   ./magefile installPreCommit
   ```

5. Update and run pre-commit hooks locally:

   ```bash
   ./magefile runPreCommit
   ```

6. Set up `go.mod` for development:

   ```bash
   ./magefile localGoMod
   ```
