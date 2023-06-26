package git_test

import (
	"fmt"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitutils "github.com/l50/goutils/v2/git"
)

func ExampleAddFile() {
	filePath := "/path/to/your/dummy/file"
	if err := gitutils.AddFile(filePath); err != nil {
		log.Fatalf("failed to stage file: %v", err)
	}
	log.Printf("Staged file: %s", filePath)
}

func ExampleCommit() {
	repo, _ := git.PlainOpen("/path/to/dummy/repo")
	msg := "Dummy commit message"
	if err := gitutils.Commit(repo, msg); err != nil {
		log.Fatalf("failed to create commit: %v", err)
	}
}

func ExampleCloneRepo() {
	url := "https://github.com/dummy/repo.git"
	clonePath := "/path/to/dummy/repo"
	auth := &http.BasicAuth{
		Username: "dummy_username",
		Password: "dummy_password",
	}
	_, err := gitutils.CloneRepo(url, clonePath, auth)
	if err != nil {
		log.Fatalf("failed to clone repository: %v", err)
	}
}

// ExampleGetTags demonstrates usage of GetTags function.
func ExampleGetTags() {
	repo, err := git.PlainOpen("/path/to/repo")
	if err != nil {
		log.Fatalf("failed to open repository: %v", err)
	}

	_, err = gitutils.GetTags(repo)
	if err != nil {
		log.Fatalf("failed to get tags: %v", err)
	}
}

// ExampleGetGlobalUserCfg demonstrates usage of GetGlobalUserCfg function.
func ExampleGetGlobalUserCfg() {
	_, err := gitutils.GetGlobalUserCfg()
	if err != nil {
		log.Fatalf("failed to retrieve global git user settings: %v", err)
	}
}

// ExampleCreateTag demonstrates usage of CreateTag function.
func ExampleCreateTag() {
	repo, err := git.PlainOpen("/path/to/repo")
	if err != nil {
		log.Fatalf("failed to open repository: %v", err)
	}

	tag := "v1.0.0"

	if err := gitutils.CreateTag(repo, tag); err != nil {
		log.Fatalf("failed to create tag: %v", err)
	}
}

// ExamplePush demonstrates usage of Push function.
func ExamplePush() {
	repo, err := git.PlainOpen("/path/to/repo")
	if err != nil {
		log.Fatalf("failed to open repository: %v", err)
	}

	auth := &http.BasicAuth{
		Username: "your_username",
		Password: "your_password",
	}

	if err := gitutils.Push(repo, auth); err != nil {
		log.Fatalf("failed to push to remote: %v", err)
	}
}

// ExamplePushTag demonstrates usage of PushTag function.
func ExamplePushTag() {
	repo, err := git.PlainOpen("/path/to/repo")
	if err != nil {
		log.Fatalf("failed to open repository: %v", err)
	}

	tag := "v1.0.0"

	auth := &http.BasicAuth{
		Username: "your_username",
		Password: "your_password",
	}

	if err := gitutils.PushTag(repo, tag, auth); err != nil {
		log.Fatalf("failed to push tag: %v", err)
	}
}

// ExampleDeleteTag demonstrates usage of DeleteTag function.
func ExampleDeleteTag() {
	_, err := git.PlainOpen("/path/to/repo")
	if err != nil {
		log.Fatalf("failed to open repository: %v", err)
	}
}

// ExamplePullRepos demonstrates usage of PullRepos function.
func ExamplePullRepos() {
	dirs := []string{"/path/to/your/directory", "/another/path/to/your/directory"}

	if err := gitutils.PullRepos(dirs...); err != nil {
		log.Fatalf("failed to pull repos: %v", err)
	}
}

// ExampleRepoRoot demonstrates usage of RepoRoot function.
func ExampleRepoRoot() {
	root, err := gitutils.RepoRoot()
	if err != nil {
		log.Fatalf("failed to retrieve root: %v", err)
	}

	fmt.Printf("The root of the current Git repository is: %s\n", root)
}
