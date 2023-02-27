package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// KeeperRecord represents a record maintained by Keeper.
type KeeperRecord struct {
	UID      string
	Title    string
	Username string
	Password string
}

// CommanderInstalled returns true if keeper
// commander is installed on the current system.
func CommanderInstalled() bool {
	return CmdExists("keeper")
}

// KeeperLoggedIn returns true if keeper vault
// is logged in with the input `email`.
// Otherwise, it returns false.
func KeeperLoggedIn() bool {
	if !CommanderInstalled() {
		err := errors.New(color.RedString(
			"keeper commander is not installed - please install and try again"))
		fmt.Println(err)
		return false
	}

	home, err := GetHomeDir()
	cobra.CheckErr(err)

	if err := Cd(filepath.Join(home, ".keeper")); err != nil {
		fmt.Print("failed to change into the keeper config directory: ", err)
		return false
	}

	fmt.Println(color.YellowString(
		"Checking if we are logged into Keeper vault"))
	loggedIn := "My Vault>"

	out, err := RunCommandWithTimeout(5, "keeper", "shell")
	if err != nil {
		fmt.Print("failed to check login state "+
			"for keeper vault: ", err)
		return false
	}

	if strings.Contains(out, loggedIn) {
		return true
	}

	return false
}

// RetrieveKeeperPW returns the password found at
// the input keeperPath.
func RetrieveKeeperPW(keeperPath string) (string, error) {
	// Ensure keeper commander is installed and
	// there is a valid keeper session.
	if !CommanderInstalled() || !KeeperLoggedIn() {
		return "", errors.New(
			color.RedString("error: ensure keeper commander is installed and a valid keeper session is established"))
	}

	fmt.Printf("Retrieving password from %s in keeper", keeperPath)

	// Get password
	pw, err := RunCommand("keeper", "find-password", keeperPath)
	if err != nil {
		return "", err
	}

	// Remove newlines from output
	return strings.ReplaceAll(string(pw), "\n", ""), nil
}

// SearchKeeperRecords searches the logged-in user's
// keeper records for the input query. The searchTerm
// can be a string or regex.
func SearchKeeperRecords(searchTerm string) (KeeperRecord, error) {
	var record KeeperRecord

	// Ensure keeper commander is installed and
	// there is a valid keeper session.
	if !CommanderInstalled() || !KeeperLoggedIn() {
		return record, errors.New(
			color.RedString("error: ensure keeper commander is installed and a valid keeper session is established"))
	}

	fmt.Println(color.YellowString(
		"Searching keeper for records matching %s, please wait...", searchTerm))

	cmd := []string{"keeper", "search", searchTerm}
	output, err := RunCommand(cmd[0], cmd[1:]...)
	if err != nil {
		return record, err
	}

	// Regular expressions to extract relevant information from the output.
	uidRegex := regexp.MustCompile(`UID:\s+(\S+)`)
	titleRegex := regexp.MustCompile(`Title:\s+(.+)`)
	usernameRegex := regexp.MustCompile(`Login:\s+(\S+)`)
	passwordRegex := regexp.MustCompile(`Password:\s+(\S+)`)

	record.UID = uidRegex.FindStringSubmatch(output)[1]
	record.Title = titleRegex.FindStringSubmatch(output)[1]
	record.Username = usernameRegex.FindStringSubmatch(output)[1]
	record.Password = passwordRegex.FindStringSubmatch(output)[1]

	return record, nil
}
