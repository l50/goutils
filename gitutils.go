package utils

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

// CloneRepo clones the repo specified at src to the path
// specified with dst
func CloneRepo(src string, dst string) bool {
	_, err := git.PlainClone(dst, false, &git.CloneOptions{
		URL:      src,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Printf("Failed to clone %s to %s: %v", src, dst, err)
		return false
	}

	return true
}
