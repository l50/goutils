package utils

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

// Make sure the project utilizes pre-commit,
// otherwise these utilities are not very useful to run.
func checkPCProject() error {

	pcFile := ".pre-commit-config.yaml"
	if !FileExists(pcFile) {
		return errors.New("pre-commit is not configure for the current project")
	}

	return nil
}

// Run pre-commit with input arguments.
func runPCCmd(args ...string) error {
	if err := checkPCProject(); err != nil {
		return err
	}

	for _, arg := range args {
		if err := sh.RunV("pre-commit", arg); err != nil {
			return err
		}

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
		return fmt.Errorf(color.RedString("failed to install pre-commit golang dependencies: %v", err))
	}

	return nil
}

// InstallPCHooks Installs pre-commmit hooks locally.
func InstallPCHooks() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := runPCCmd("install"); err != nil {
		return fmt.Errorf(color.RedString("failed to install pre-commit hooks: %v", err))
	}

	return nil
}

// UpdatePCHooks Updates pre-commmit hooks locally.
func UpdatePCHooks() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := runPCCmd("autoupdate"); err != nil {
		return fmt.Errorf(color.RedString("failed to update the pre-commit hooks: %v", err))
	}

	return nil
}

// ClearPCCache Clears the pre-commit cache.
func ClearPCCache() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := runPCCmd("clean"); err != nil {
		return fmt.Errorf(color.RedString("failed to clear the pre-commit cache: %v", err))
	}

	return nil
}

// RunPCHooks Runs all pre-commit hooks locally.
func RunPCHooks() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := runPCCmd("run", "--all-files"); err != nil {
		return fmt.Errorf(color.RedString("failed to run pre-commit hooks: %v", err))
	}

	return nil
}
