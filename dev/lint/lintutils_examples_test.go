package lint_test

import (
	"log"
	"path/filepath"

	lint "github.com/l50/goutils/v2/dev/lint"
	"github.com/l50/goutils/v2/sys"
)

func ExampleInstallGoPCDeps() {
	if err := sys.Cd(filepath.Join("..", "..")); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}
	if err := lint.InstallGoPCDeps(); err != nil {
		log.Fatalf("error installing dependencies: %v", err)
	}
}

func ExampleUpdatePCHooks() {
	if err := sys.Cd(filepath.Join("..", "..")); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}
	if err := lint.UpdatePCHooks(); err != nil {
		log.Fatalf("error updating hooks: %v", err)
	}
}

func ExampleClearPCCache() {
	if err := sys.Cd(filepath.Join("..", "..")); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}
	if err := lint.ClearPCCache(); err != nil {
		log.Fatalf("error clearing cache: %v", err)
	}
}

func ExampleRunPCHooks() {
	if err := sys.Cd(filepath.Join("..", "..")); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}

	// Run with a default timeout of 600.
	if err := lint.RunPCHooks(); err != nil {
		log.Fatalf("failed to run pre-commit hooks: %v", err)
	}

	// Runs with a specified timeout of 300.
	if err := lint.RunPCHooks(300); err != nil {
		log.Fatalf("failed to run pre-commit hooks: %v", err)
	}

	if err := lint.RunPCHooks(300); err != nil {
		log.Fatalf("failed to run pre-commit hooks: %v", err)
	}
}

func ExampleAddFencedCB() {
	if err := lint.AddFencedCB("README.md", "go"); err != nil {
		log.Fatalf("error modifying markdown file: %v", err)
	}
}

func ExampleRunHookTool() {
	if err := lint.RunHookTool("golangci-lint", "run"); err != nil {
		log.Fatalf("error running hook tool: %v", err)
	}
}
