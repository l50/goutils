package macos_test

import (
	"log"

	"github.com/l50/goutils/v2/macos"
)

func ExampleInstallBrewDeps() {
	brewPackages := []string{"shellcheck", "shfmt"}
	err := macos.InstallBrewDeps(brewPackages)
	if err != nil {
		log.Fatalf("failed to install brew dependencies: %v", err)
	}
}

func ExampleInstallBrewTFDeps() {
	err := macos.InstallBrewTFDeps()
	if err != nil {
		log.Fatalf("failed to install terraform brew dependencies: %v", err)
	}
}
