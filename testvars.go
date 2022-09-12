package utils

import (
	"time"

	"github.com/go-git/go-git/v5"
)

var (
	cloneDir    string
	clonePath   string
	currentTime time.Time
	err         error
	repo        *git.Repository
	repoURL     string
	tags        []string
)
