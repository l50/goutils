# goutils

[![License](https://img.shields.io/github/license/l50/goutils?label=License&style=flat&color=blue&logo=github)](https://github.com/l50/goutils/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/l50/goutils/v2)](https://goreportcard.com/report/github.com/l50/goutils/v2)
[![Tests](https://github.com/l50/goutils/actions/workflows/tests.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/tests.yaml)
[![🚨 CodeQL Analysis](https://github.com/l50/goutils/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/codeql-analysis.yaml)
[![🚨 Semgrep Analysis](https://github.com/l50/goutils/actions/workflows/semgrep.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/semgrep.yaml)
[![Coverage Status](https://coveralls.io/repos/github/l50/goutils/badge.svg?branch=main)](https://coveralls.io/github/l50/goutils?branch=main)
[![Renovate](https://github.com/l50/goutils/actions/workflows/renovate.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/renovate.yaml)

This repo is comprised of utilities that I use across various go projects.

## Dependencies

- [Install pre-commit](https://pre-commit.com/):

  ```bash
  python3 -m pip install --upgrade pip
  python3 -m pip install pre-commit
  ```

- [Install Mage](https://magefile.org/):

  ```bash
  go install github.com/magefile/mage@latest
  ```

- [Install Keeper Commander](https://github.com/Keeper-Security/Commander):

  ```bash
  python3 -m pip install --upgrade pip
  python3 -m pip install keepercommander
  ```

---

## For Contributors and Developers

1. [Fork this project](https://docs.github.com/en/get-started/quickstart/fork-a-repo)

1. Install dependencies:

   ```bash
   mage installDeps
   ```

1. Update and run pre-commit hooks locally:

   ```bash
   mage runPreCommit
   ```

---

## Create New Release

- Download and install the [gh cli tool](https://cli.github.com/):

  - [macOS](https://github.com/cli/cli#macos)
  - [Linux](https://github.com/cli/cli/blob/trunk/docs/install_linux.md)
  - [Windows](https://github.com/cli/cli#windows)

- Install changelog extension:

  ```bash
  gh extension install chelnak/gh-changelog
  ```

- Generate changelog:

  ```bash
  NEXT_VERSION=v1.1.3
  gh changelog new --next-version "${NEXT_VERSION}"
  ```

- Create release:

  ```bash
  gh release create "${NEXT_VERSION}" -F CHANGELOG.md
  ```
