package utils

import (
	"time"

	"github.com/go-git/go-git/v5"
)

var (
	currentTime     time.Time
	cloneDir        = "/tmp"
	err             error
	repo            *git.Repository
	tags            []string
	testFile        string
	testFileContent string
	testRepoURL     = "https://github.com/l50/helloworld.git"
)
