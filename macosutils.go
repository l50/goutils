package utils

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

// InstallBrewTFDeps installs dependencies for terraform projects
// with homebrew.
func InstallBrewTFDeps() error {
	brewPackages := []string{
		// Install shell tools
		"shellcheck",
		"shfmt",
		// Install terraform tools
		"terraform-docs",
		"tflint",
		"checkov",
	}

	if err := InstallBrewDeps(brewPackages); err != nil {
		return err
	}

	return nil
}

// InstallBrewDeps runs brew install for the input brew packages.
func InstallBrewDeps(brewPackages []string) error {
	for _, pkg := range brewPackages {
		err := sh.RunV("brew", "install", pkg)
		if err != nil {
			return fmt.Errorf(color.RedString(
				"failed to install dependencies: %v", err))
		}
	}

	return nil
}
