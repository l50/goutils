package utils

import (
	"os"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	repo := "https://github.com/l50/helloworld"
	cloned := CloneRepo(repo, "helloworld")
	if !cloned {
		t.Fatal("Failed to clone ", repo)
	} else {
		os.RemoveAll("helloworld")
	}
}
