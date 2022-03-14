package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	targetDir := "/tmp"
	repo := "https://github.com/l50/helloworld.git"
	cloneLoc := filepath.Join(targetDir, "helloworld")

	if FileExists(cloneLoc) {
		os.RemoveAll(cloneLoc)
	}

	cloned := CloneRepo(repo, cloneLoc)
	if !cloned {
		t.Fatalf("failed to clone %s", repo)
	}

	defer os.RemoveAll(cloneLoc)
}
