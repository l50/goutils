package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

var pc = sh.RunCmd("pre-commit")

// Make sure the project utilizes pre-commit,
// otherwise these utilities are not very useful to run.
func checkPCProject() error {
	var pcFile string
	cwd := Gwd()

	if strings.Contains(cwd, ".mage") {
		pcFile = filepath.Join("..", ".pre-commit-config.yaml")
	} else {
		pcFile = ".pre-commit-config.yaml"
	}
	if !FileExists(pcFile) {
		return errors.New(color.RedString(
			"pre-commit is not configured for the current project"))
	}

	return nil
}

// InstallGoPCDeps Installs dependencies used for pre-commit with golang projects.
func InstallGoPCDeps() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	deps := []string{
		"golang.org/x/lint/golint",
		"golang.org/x/tools/cmd/goimports",
		"github.com/fzipp/gocyclo/cmd/gocyclo",
		"github.com/golangci/golangci-lint/cmd/golangci-lint",
		"github.com/go-critic/go-critic/cmd/gocritic",
	}

	if err := InstallGoDeps(deps); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to install pre-commit golang dependencies: %v", err))
	}

	return nil
}

// InstallPCHooks Installs pre-commmit hooks locally.
func InstallPCHooks() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := pc("install"); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to install pre-commit hooks: %v", err))
	}

	return nil
}

// UpdatePCHooks Updates pre-commmit hooks locally.
func UpdatePCHooks() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := pc("autoupdate"); err != nil {
		return fmt.Errorf(color.RedString("failed to update the pre-commit hooks: %v", err))
	}

	return nil
}

// ClearPCCache Clears the pre-commit cache.
func ClearPCCache() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := pc("clean"); err != nil {
		return fmt.Errorf(color.RedString("failed to clear the pre-commit cache: %v", err))
	}

	return nil
}

// Something here is simply not working properly - test with ./magefile runPreCommit
// // RunPCHooks Runs all pre-commit hooks locally.
// func RunPCHooks() error {
// 	// if err := checkPCProject(); err != nil {
// 	// 	return err
// 	// }

// 	// _, err := RunCommand("pre-commit", "run", "--all-files")
// 	// if err != nil {
// 	// 	return fmt.Errorf(color.RedString("failed to run pre-commit hooks: %v", err))
// 	// }
// 	if err := sh.RunV("pre-commit", "run", "--all-files"); err != nil {
// 		return fmt.Errorf(color.RedString("failed to run pre-commit hooks: %v", err))
// 	}

// 	return nil
// }
