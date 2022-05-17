package utils

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

// InstallVSCodeModules installs the modules used by the vscode-go extension in VSCode.
func InstallVSCodeModules() error {
	fmt.Println(color.YellowString("Installing vscode-go dependencies."))
	vscodeDeps := []string{
		"github.com/uudashr/gopkgs/v2/cmd/gopkgs",
		"github.com/ramya-rao-a/go-outline",
		"github.com/cweill/gotests/gotests",
		"github.com/fatih/gomodifytags",
		"github.com/josharian/impl",
		"github.com/haya14busa/goplay/cmd/goplay",
		"github.com/go-delve/delve/cmd/dlv",
		"honnef.co/go/tools/cmd/staticcheck",
		"golang.org/x/tools/gopls",
		"github.com/rogpeppe/godef",
	}

	if err := InstallGoDeps(vscodeDeps); err != nil {
		return fmt.Errorf(
			color.RedString("failed to install vscode-go dependencies: %v", err))
	}

	return nil
}

// Tidy runs go mod tidy.
func Tidy() error {
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return fmt.Errorf(
			color.RedString("failed to run go mod tidy: %v", err))
	}

	return nil
}

// UpdateMageDeps updates mage-specific dependencies
// using the input path to the associated go.mod.
func UpdateMageDeps(magedir string) error {
	cwd := Gwd()
	if err := Cd(magedir); err != nil {
		return fmt.Errorf(
			color.RedString(
				"failed to cd from %s to %s: %v", cwd, magedir, err))
	}

	if err := Tidy(); err != nil {
		return fmt.Errorf(
			color.RedString(
				"failed to update mage dependencies in %s: %v", magedir, err))
	}

	if err := Cd(cwd); err != nil {
		return fmt.Errorf(
			color.RedString(
				"failed to cd from %s to %s: %v", magedir, cwd, err))
	}

	return nil
}

// InstallGoDeps runs go install for the input dependencies.
func InstallGoDeps(deps []string) error {
	var err error
	failed := false

	for _, dep := range deps {
		if err := sh.RunV("go", "install", dep+"@latest"); err != nil {
			failed = true
		}
	}

	if failed {
		return fmt.Errorf(
			color.RedString("failed to install input go dependencies: %w", err))
	}

	return nil
}

// Create artifacts for upload to github
func CreateArtifacts(os []string, binPath string) error {
	operatingSystems := os
	for _, os := range operatingSystems {
		err := sh.Run("mage", "-d", ".mage", "-compile", binPath+"-"+os)
		if err != nil {
			return fmt.Errorf(
				color.RedString("failed to create artifacts at %s: %v", binPath, err))
		}
	}

	return nil
}
