package macos

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

// InstallBrewDeps installs the input brew packages by running brew install.
// If any installation fails, it returns an error.
//
// Parameters:
//
// brewPackages: A slice of strings representing the brew packages to be installed.
//
// Returns:
//
// error: An error if any brew package fails to install.
//
// Example:
//
// brewPackages := []string{"shellcheck", "shfmt"}
// err := macos.InstallBrewDeps(brewPackages)
//
//	if err != nil {
//	  log.Fatalf("failed to install brew dependencies: %v", err)
//	}
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

// InstallBrewTFDeps installs dependencies for terraform projects using homebrew.
// The dependencies include several shell and terraform tools.
// If any installation fails, it returns an error.
//
// Returns:
//
// error: An error if any brew package fails to install.
//
// Example:
//
// err := macos.InstallBrewTFDeps()
//
//	if err != nil {
//	  log.Fatalf("failed to install terraform brew dependencies: %v", err)
//	}
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
