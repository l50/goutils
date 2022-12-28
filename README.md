# goutils

[![Go Report Card](https://goreportcard.com/badge/github.com/l50/goutils)](https://goreportcard.com/report/github.com/l50/goutils)
[![License](https://img.shields.io/github/license/l50/goutils?label=License&style=flat&color=blue&logo=github)](https://github.com/l50/goutils/blob/main/LICENSE)
[![Tests](https://github.com/l50/goutils/actions/workflows/tests.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/tests.yaml)
[![🚨 CodeQL Analysis](https://github.com/l50/goutils/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/codeql-analysis.yaml)
[![🚨 Semgrep Analysis](https://github.com/l50/goutils/actions/workflows/semgrep.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/semgrep.yaml)
[![Pre-Commit](https://github.com/l50/goutils/actions/workflows/pre-commit.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/pre-commit.yaml)
[![Coverage Status](https://coveralls.io/repos/github/l50/goutils/badge.svg?branch=main)](https://coveralls.io/github/l50/goutils?branch=main)

This repo is comprised of utilities that I use across various go projects.

## Dependencies

- [Install gvm](https://github.com/moovweb/gvm):

  ```bash
  bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
  ```

- [Install golang](https://go.dev/):

  ```bash
  source .gvm
  ```

- [Install pre-commit](https://pre-commit.com/):

  ```bash
  brew install pre-commit
  ```

- [Install Mage](https://magefile.org/):

  ```bash
  go install github.com/magefile/mage@latest
  ```

---

## Developer Environment Setup

1. [Fork this project](https://docs.github.com/en/get-started/quickstart/fork-a-repo)

1. (Optional) If you installed gvm, create golang pkgset specifically for this project:

   ```bash
   source "${HOME}/.gvm"
   ```

1. Install pre-commit hooks and dependencies:

   ```bash
   mage installPreCommitHooks
   ```

1. Update and run pre-commit hooks locally:

   ```bash
   mage runPreCommit
   ```

---

## Create New Release

This requires the [GitHub CLI](https://github.com/cli/cli#installation)

```bash
gh release create
```
