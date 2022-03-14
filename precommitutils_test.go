package utils

import "testing"

func TestInstallGoPCDeps(t *testing.T) {
	if err := InstallGoPCDeps(); err != nil {
		t.Fatal(err)
	}
}

func TestInstallPCHooks(t *testing.T) {
	if err := InstallPCHooks(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdatePCHooks(t *testing.T) {
	if err := UpdatePCHooks(); err != nil {
		t.Fatal(err)
	}
}

func TestClearPCCache(t *testing.T) {
	if err := ClearPCCache(); err != nil {
		t.Fatal(err)
	}
}

// This isn't worth running - it takes forever and the function gets plenty
// of testing with the numerous magefiles I have across projects.
// func TestRunPCHooks(t *testing.T) {
// 	if err := RunPCHooks(); err != nil {
// 		t.Fatal(err)
// 	}
// }
