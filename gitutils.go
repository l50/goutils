package utils

import (
	"gopkg.in/src-d/go-git.v4"
	"log"
	"os"
)

// CloneRepo clones the repo specified at src to the path
// specified with dst
func CloneRepo(src string, dst string) bool {
	_, err := git.PlainClone((dst), false, &git.CloneOptions{
		URL:      src,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatal("Failed to clone ", src)
		return false
	}
	return true
}