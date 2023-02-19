package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var (
	mageCleanupDirs []string
)

func init() {
	// Create test repo and queue it for cleanup
	randStr, _ := RandomString(8)
	clonePath = createTestRepo(fmt.Sprintf("mageutils-%s", randStr))
	mageCleanupDirs = append(mageCleanupDirs, clonePath)
}

func TestGHRelease(t *testing.T) {
	// Call the function with an old version
	newVer := "v1.0.0"
	if err := GHRelease(newVer); err == nil {
		t.Errorf("release %s should not have been created: %v", newVer, err)
	}
}

func TestGoReleaser(t *testing.T) {
	t.Cleanup(func() {
		cleanupMageUtils(t)
	})

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	releaserDir := filepath.Join(cwd, "dist")
	mageCleanupDirs = append(mageCleanupDirs, releaserDir)

	if err := GoReleaser(); err != nil {
		t.Fatal(err)
	}
}

func TestInstallVSCodeModules(t *testing.T) {
	if err := InstallVSCodeModules(); err != nil {
		t.Fatal(err)
	}
}

func TestModUpdate(t *testing.T) {
	// First test
	recursive := false
	verbose := true
	if err := ModUpdate(recursive, verbose); err != nil {
		t.Fatal(err)
	}

	// Second test
	recursive = true
	verbose = false
	if err := ModUpdate(recursive, verbose); err != nil {
		t.Fatal(err)
	}
}

func TestTidy(t *testing.T) {
	if err := Tidy(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateMageDeps(t *testing.T) {
	if err := UpdateMageDeps("magefiles"); err != nil {
		t.Fatal(err)
	}
}

func TestInstallGoDeps(t *testing.T) {
	sampleDeps := []string{
		"golang.org/x/lint/golint",
		"golang.org/x/tools/cmd/goimports",
	}

	if err := InstallGoDeps(sampleDeps); err != nil {
		t.Fatal(err)
	}
}

func cleanupMageUtils(t *testing.T) {
	for _, dir := range mageCleanupDirs {
		if err := RmRf(dir); err != nil {
			fmt.Println("failed to clean up mageUtils: ", err.Error())
		}
	}
}
