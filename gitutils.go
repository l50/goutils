package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/magefile/mage/sh"
)

var (
	errMsg string
	tags   []string
)

// GitConfigUserInfo holds a username and
// email to use for user.name and user.email.
type GitConfigUserInfo struct {
	User  string
	Email string
}

// SSHKeyInfo is used to hold the name of an SSH Key
// file and a password for decryption.
type SSHKeyInfo struct {
	Name string
	PW   string
}

// GetSSHPubKey returns the public SSH key for the input
// `keyName` and uses the input `password` (if provided)
// to decrypt the associated private key.
func GetSSHPubKey(keyName string, password string) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys

	sshPath := os.Getenv("HOME") + "/.ssh/" + keyName
	publicKey, err := ssh.NewPublicKeysFromFile("git", sshPath, password)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

// CloneRepo clones the repo specified with the input `url` to
// `clonePath`.
//
// To clone a repo using an SSH key, provide
// the name of the key file for `sshKey.Name`.
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
			errMsg = fmt.Sprint(color.RedString(
				"%s was already cloned to %s", url, clonePath))
			return nil, errors.New(errMsg)
		}
		errMsg = fmt.Sprint(color.RedString(
			"failed to clone %s to %s: %v", url, clonePath, err))
		return nil, errors.New(errMsg)
	}

	return repo, nil
}

// GetTags returns the tags for an input `repo`.
func GetTags(repo *git.Repository) ([]string, error) {
	tagObjects, err := repo.TagObjects()
	if err != nil {
		errMsg = fmt.Sprint(color.RedString(
			"failed to retrieve repo tags: %v", err))
		return tags, errors.New(errMsg)
	}

	err = tagObjects.ForEach(func(t *object.Tag) error {
		tags = append(tags, t.Name)
		return nil
	})

	if err != nil {
		errMsg = fmt.Sprint(color.RedString(
			"failed to retrieve repo tags: %v", err))
		return tags, errors.New(errMsg)
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
func GetGlobalUserCfg() (GitConfigUserInfo, error) {
	userInfo := GitConfigUserInfo{}
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

// PushTag is used to push a tag to remote.
func PushTag(repo *git.Repository, tag string, auth transport.AuthMethod) error {
	pushOptions := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs: []config.RefSpec{config.RefSpec(
			"refs/tags/*:refs/tags/*")},
		Auth: auth,
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
func DeleteTag(repo *git.Repository, tag string) error {
	if err := repo.DeleteTag(tag); err != nil {
		return fmt.Errorf(color.RedString(
			"error deleting local %s tag: %v", tag, err))
	}

	return nil
}

// DeletePushedTag deletes an input `tag` that has been pushed
// to remote.
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
