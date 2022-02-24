package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	targetDir := "/tmp"
	repo := "https://github.com/l50/helloworld.git"
	cloneLoc := filepath.Clean(filepath.Join(targetDir, "helloworld"))

	if FileExists(cloneLoc) {
		os.RemoveAll(cloneLoc)
	}

	cloned := CloneRepo(repo, cloneLoc)
	if !cloned {
		t.Fatal("Failed to clone ", repo)
	}

	defer os.RemoveAll(cloneLoc)
}
