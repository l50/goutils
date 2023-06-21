package macos

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// InstallBrewDeps executes brew install for the input packages.
// If any installation fails, it returns an error.
//
// **Parameters:**
//
// brewPackages: Slice of strings representing the packages to install.
//
// **Returns:**
//
// error: An error if any package fails to install.
func InstallBrewDeps(brewPackages []string) error {
	for _, pkg := range brewPackages {
		err := sh.RunV("brew", "install", pkg)
		if err != nil {
			return fmt.Errorf("failed to install dependencies: %v", err)
		}
	}

	return nil
}

// InstallBrewTFDeps installs dependencies for terraform projects
// using homebrew. The dependencies include several shell and
// terraform tools. If any installation fails, it returns an error.
//
// **Returns:**
//
// error: An error if any package fails to install.
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
