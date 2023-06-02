package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	
	"github.com/l50/goutils/sys"
	"github.com/magefile/mage/sh"
)

// ConfigUserInfo holds a username and
// email to use for user.name and user.email.
type ConfigUserInfo struct {
	User  string
	Email string
}


// AddFile stages the file at the given file path in its associated Git repository. Returns an error if one occurs.
//
// Parameters:
//
// filePath: A string representing the path to the file to be staged.
//
// Returns:
//
// error: An error if one occurs during the staging process.
//
// Example:
//
// filePath := "/path/to/your/file"
// err := AddFile(filePath)
//
//	if err != nil {
//	  log.Fatalf("failed to stage file: %v", err)
//	}
//
// log.Printf("Staged file: %s", filePath)
func AddFile(filePath string) error {
	repo, err := git.PlainOpen(filepath.Dir(filePath))
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed to open %s repo: %v", repo, err))
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed to retrieve worktree: %v", err))
	}

	_, err = w.Add(filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed to run `git add` on %s: %v", filePath, err))
	}

	status, err := w.Status()
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed to run `git status` on %s: %v", filePath, err))
	}

	if status.IsClean() {
		return fmt.Errorf(color.RedString(
			"status is clean - failed to run `git add` "+
				"on %s: %v", filePath, err))
	}

	return nil
}

// Commit creates a new commit in the given repository with the provided message. The author of the commit is retrieved from the global Git user settings.
//
// Parameters:
//
// repo: A pointer to the Repository struct representing the repository where the commit should be created.
//
// msg: A string representing the commit message.
//
// Returns:
//
// error: An error if the commit cannot be created.
//
// Example:
//
// repo, err := git.PlainOpen("/path/to/repo")
//
//	if err != nil {
//	  log.Fatalf("failed to open repository: %v", err)
//	}
//
// msg := "Commit message"
// err = Commit(repo, msg)
//
//	if err != nil {
//	  log.Fatalf("failed to create commit: %v", err)
//	}
func Commit(repo *git.Repository, msg string) error {
	cfg, err := GetGlobalUserCfg()
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed get repo config: %v", err))
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed to retrieve worktree: %v", err))
	}

	commit, err := w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  cfg.User,
			Email: cfg.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed to commit current staging area`: %v",
			err))
	}

	obj, err := repo.CommitObject(commit)
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed to run `git show`: %v", err))
	}

	if obj.Author.Email != cfg.Email {
		return fmt.Errorf(color.RedString(
			"author email in commit doesn't match "+
				"global git config email - Commit() failed: %v",
			err))
	}

	return nil
}

// CloneRepo clones the repo specified with the input `url` to
// clonePath.
func CloneRepo(url string, clonePath string, auth transport.AuthMethod) (
	*git.Repository, error) {
	var err error
	var repo *git.Repository
	var cloneOptions *git.CloneOptions

	if auth != nil {
		cloneOptions = &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
			Auth:     auth,
		}
	} else {
		cloneOptions = &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		}
	}

	repo, err = git.PlainClone(clonePath, false, cloneOptions)
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			return nil, fmt.Errorf(color.RedString(
				"%s was already cloned to %s", url, clonePath))
		}
		return nil, fmt.Errorf(color.RedString(
			"failed to clone %s to %s: %v", url, clonePath, err))
	}

	return repo, nil
}

// GetSSHPubKey retrieves the public SSH key for the given key name, decrypting the associated private key if a password
// is provided. It returns a pointer to the public key object, or an error if one occurs.
//
// Parameters:
//
// keyName: A string representing the name of the key to retrieve.
// password: A string representing the password used to decrypt the private key.
//
// Returns:
//
// *ssh.PublicKeys: A pointer to a PublicKeys object representing the retrieved public key.
// error: An error if one occurs during key retrieval or decryption.
//
// Example:
//
// keyName := "id_rsa"
// password := "mypassword"
// publicKey, err := GetSSHPubKey(keyName, password)
//
//	if err != nil {
//	  log.Fatalf("failed to get SSH public key: %v", err)
//	}
//
// log.Printf("Retrieved public key: %v", publicKey)
func GetSSHPubKey(keyName string, password string) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys

	sshPath := filepath.Join(os.Getenv("HOME"), ".ssh", keyName)
	publicKey, err := ssh.NewPublicKeysFromFile("git", sshPath, password)
	if err != nil {
		return nil,
			fmt.Errorf(color.RedString(
				"failed to retrieve public SSH key %s: %v",
				keyName, err))
	}

	return publicKey, nil
}

// GetTags returns the tags for an input repo.
func GetTags(repo *git.Repository) ([]string, error) {
	var tags []string
	tagObjects, err := repo.TagObjects()
	if err != nil {
		return []string{}, fmt.Errorf(color.RedString(
			"failed to retrieve repo tags: %v", err))
	}

	err = tagObjects.ForEach(func(t *object.Tag) error {
		tags = append(tags, t.Name)
		return nil
	})

	if err != nil {
		return tags, fmt.Errorf(color.RedString(
			"failed to retrieve repo tags: %v", err))
	}

	return tags, err
}

func tagExists(repo *git.Repository, tag string) (bool, error) {
	tags, err := GetTags(repo)
	if err != nil {
		return false, err
	}

	for _, t := range tags {
		if t == tag {
			return true, nil
		}
	}

	return false, nil
}

// GetGlobalUserCfg returns the username and email from
// the global git user settings.
func GetGlobalUserCfg() (ConfigUserInfo, error) {
	userInfo := ConfigUserInfo{}
	var err error

	userInfo.User, err = sh.Output("git", "config", "user.name")
	if err != nil {
		return userInfo, fmt.Errorf(color.RedString(
			"failed to retrieve global git username: %v", err))
	}

	userInfo.Email, err = sh.Output("git", "config", "user.email")
	if err != nil {
		return userInfo, fmt.Errorf(color.RedString(
			"failed to retrieve global git email: %v", err))
	}

	return userInfo, nil
}

// CreateTag is used to create an input `tag` in the
// specified `repo` if it doesn't already exist.
// Resource: https://github.com/go-git/go-git/blob/bf3471db54b0255ab5b159005069f37528a151b7/_examples/tag-create-push/main.go
func CreateTag(repo *git.Repository, tag string) error {
	exists, err := tagExists(repo, tag)
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed to retrieve repo tags: %v", err))
	}

	if exists {
		return fmt.Errorf(color.RedString(
			"error creating input tag %s: it already exists", tag))
	}

	cfg, err := GetGlobalUserCfg()
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed get repo config: %v", err))
	}

	h, err := repo.Head()
	if err != nil {
		return fmt.Errorf(color.RedString(
			"failed to get repo head: %v", err))
	}

	_, err = repo.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  cfg.User,
			Email: cfg.Email,
			When:  time.Now(),
		},
		Message: tag,
	})

	if err != nil {
		return fmt.Errorf(color.RedString(
			"error creating input tag %s: %v", tag, err))
	}

	return nil
}

// Push pushes the contents of the input
// repo to the default remote (origin).
func Push(repo *git.Repository, auth transport.AuthMethod) error {
	var pushOptions *git.PushOptions

	if auth != nil {
		pushOptions = &git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			Auth:       auth,
		}
	} else {
		pushOptions = &git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
		}
	}

	err := repo.Push(pushOptions)

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Print(color.YellowString(
				"origin remote is up-to-date, no push was executed."))
			return nil
		}
		return fmt.Errorf(color.RedString(
			"error pushing to remote origin: %v", err))
	}

	return nil
}

// PushTag is used to push a tag to remote.
func PushTag(repo *git.Repository, tag string, auth transport.AuthMethod) error {
	var pushOptions *git.PushOptions

	if auth != nil {
		pushOptions = &git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			RefSpecs: []config.RefSpec{config.RefSpec(
				"refs/tags/*:refs/tags/*")},
			Auth: auth,
		}
	} else {
		pushOptions = &git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			RefSpecs: []config.RefSpec{config.RefSpec(
				"refs/tags/*:refs/tags/*")},
		}
	}

	err := repo.Push(pushOptions)

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Print(color.YellowString(
				"origin remote is up-to-date, no push was executed."))
			return nil
		}

		return fmt.Errorf(color.RedString(
			"error pushing %s tag to remote origin: %v", tag, err))
	}

	return nil
}

// DeleteTag deletes the local input `tag` from the
// specified repo.
// DeletePushedTag deletes a tag from a given repository that has been pushed to a remote. The tag is deleted from both the local repository and the remote repository.
//
// Parameters:
//
// repo: A pointer to the Repository struct representing the repository where the tag should be deleted.
//
// tag: A string representing the tag that should be deleted.
//
// auth: An AuthMethod representing the method used to authenticate with the remote repository.
//
// Returns:
//
// error: An error if the tag cannot be deleted.
//
// Example:
//
// repo, err := git.PlainOpen("/path/to/repo")
//
//	if err != nil {
//	  log.Fatalf("failed to open repository: %v", err)
//	}
//
// tag := "v1.0.0"
// auth, err := ssh.NewSSHAgentAuth("git")
//
//	if err != nil {
//	  log.Fatalf("failed to create authentication method: %v", err)
//	}
//
// err = DeletePushedTag(repo, tag, auth)
//
//	if err != nil {
//	  log.Fatalf("failed to delete pushed tag: %v", err)
//	}
func DeleteTag(repo *git.Repository, tag string) error {
	if err := repo.DeleteTag(tag); err != nil {
		return fmt.Errorf(color.RedString(
			"error deleting local %s tag: %v", tag, err))
	}

	return nil
}

// DeletePushedTag deletes a tag from a given repository that has been pushed remote.
// The tag is deleted from both the local repository and the remote repository.
//
// Parameters:
//
// repo: A pointer to the Repository struct representing the repository where the tag should be deleted.
//
// tag: A string representing the tag that should be deleted.
//
// auth: An AuthMethod representing the method used to authenticate with the remote repository.
//
// Returns:
//
// error: An error if the tag cannot be deleted.
//
// Example:
//
// repo, err := git.PlainOpen("/path/to/repo")
//
//	if err != nil {
//	  log.Fatalf("failed to open repository: %v", err)
//	}
//
// tag := "v1.0.0"
// auth, err := ssh.NewSSHAgentAuth("git")
//
//	if err != nil {
//	  log.Fatalf("failed to create authentication method: %v", err)
//	}
//
// err = DeletePushedTag(repo, tag, auth)
//
//	if err != nil {
//	  log.Fatalf("failed to delete pushed tag: %v", err)
//	}
func DeletePushedTag(repo *git.Repository, tag string, auth transport.AuthMethod) error {
	err := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs: []config.RefSpec{config.RefSpec(
			"refs/tags/" + tag)},
		Auth: auth,
	})

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Print(color.YellowString(
				"origin remote is up-to-date, no delete was executed."))
			return nil
		}

		return fmt.Errorf(color.RedString(
			"error deleting pushed tag %s: %v", tag, err))
	}

	return nil
}

// PullRepos updates all git repositories found in the given directories by pulling changes from the upstream branch.
// It looks for repositories by finding directories with a ".git" subdirectory.
// If a repository is not on the default branch, it will switch to the default branch before pulling changes.
// Returns an error if any step of the process fails.
func PullRepos(dirs ...string) error {
	// Save the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	defer func() {
		// Change the working directory back to the original directory
		if err := sys.Cd(wd); err != nil {
			fmt.Printf("failed to change directory back to %s: %v\n", wd, err)
		}
	}()

	for _, dir := range dirs {
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == ".git" {
				// Get the path to the parent directory of the .git directory
				repoDir := filepath.Dir(path)

				// Change to the repository directory
				if err := sys.Cd(repoDir); err != nil {
					return fmt.Errorf("failed to change directory to %s: %w", repoDir, err)
				}

				// Get the current branch
				refOutput, err := sys.RunCommand("git", "symbolic-ref", "HEAD")
				if err != nil {
					// If no branch is checked out, get the default branch
					defaultBranchOutput, defaultBranchErr := sys.RunCommand("git", "config", "--get", "init.defaultBranch")
					if defaultBranchErr != nil {
						return fmt.Errorf("failed to get current or default branch for %s: %v, %v", repoDir, err, defaultBranchErr)
					}
					refOutput = defaultBranchOutput
				}

				// Pull changes in the current branch
				ref := strings.TrimSpace(strings.TrimPrefix(refOutput, "refs/heads/"))
				res, err := sys.RunCommand("git", "pull", "origin", ref)
				if err != nil {
					fmt.Printf("failed to update %s: %s\n", repoDir, res)
				} else if strings.TrimSpace(res) != "Already up to date." {
					fmt.Printf("Now pulling the latest from upstream for %s\n", repoDir)
				}
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed to find git repositories in %s: %w", dir, err)
		}
	}
	return nil
}

// RepoRoot returns the root of the git repo a user is currently in.
func RepoRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(cwd, ".git")); err == nil {
			return cwd, nil
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			return "", fmt.Errorf("git root not found")
		}
		cwd = parent
	}
}
