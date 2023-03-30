# goutils

[![License](https://img.shields.io/github/license/l50/goutils?label=License&style=flat&color=blue&logo=github)](https://github.com/l50/goutils/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/l50/goutils)](https://goreportcard.com/report/github.com/l50/goutils)
[![Tests](https://github.com/l50/goutils/actions/workflows/tests.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/tests.yaml)
[![🚨 CodeQL Analysis](https://github.com/l50/goutils/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/codeql-analysis.yaml)
[![🚨 Semgrep Analysis](https://github.com/l50/goutils/actions/workflows/semgrep.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/semgrep.yaml)
[![Coverage Status](https://coveralls.io/repos/github/l50/goutils/badge.svg?branch=main)](https://coveralls.io/github/l50/goutils?branch=main)
[![Renovate](https://github.com/l50/goutils/actions/workflows/renovate.yaml/badge.svg)](https://github.com/l50/goutils/actions/workflows/renovate.yaml)

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

1. (Optional) If you installed gvm:

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

1. Install keeper commander:

   ```bash
   python3 -m pip install --upgrade pip
   python3 -m pip install keepercommander
   ```

1. Install `rvm` and `ruby` (for markdownlint):

   ```bash
   gpg --keyserver keyserver.ubuntu.com \
     --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB
   \curl -sSL https://get.rvm.io | bash -s stable --ruby
   rvm install ruby-3.2.1
   ```

---

## Create New Release

This requires the [GitHub CLI](https://github.com/cli/cli#installation)
and [gh-changelog GitHub CLI extension](https://github.com/chelnak/gh-changelog).

Install changelog extension:

```bash
gh extension install chelnak/gh-changelog
```

Generate changelog:

```bash
NEXT_VERSION=v1.1.3
gh changelog new --next-version "${NEXT_VERSION}"
```

Create release:

```bash
gh release create "${NEXT_VERSION}" -F CHANGELOG.md
```
