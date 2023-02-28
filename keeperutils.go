package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// KeeperRecord represents a record maintained by Keeper.
type KeeperRecord struct {
	UID      string
	Title    string
	Username string
	Password string
}

// keeperConfigPath returns the path of the keeper config file.
func keeperConfigPath() (string, error) {
	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".keeper", "config.json"), nil
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

	fmt.Println(color.YellowString(
		"Checking if we are logged into Keeper vault"))

	configPath, err := keeperConfigPath()
	if err != nil {
		err := errors.New(color.RedString(
			"failed to retrieve keeper config path"))
		fmt.Println(err)
		return false
	}

	loggedIn := "My Vault>"
	out, err := RunCommandWithTimeout(5, "keeper", "shell", "--config", configPath)
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
		return "", errors.New("error: ensure keeper commander is installed and a valid keeper session is established")
	}

	fmt.Printf("Retrieving password from %s in keeper\n", keeperPath)

	// Get password
	configPath, err := keeperConfigPath()
	if err != nil {
		err := errors.New(color.RedString(
			"failed to retrieve keeper config path"))
		return "", err
	}
	pw, err := RunCommand("keeper", "find-password", keeperPath, "--config", configPath)
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
		return record, errors.New("error: ensure keeper commander is installed and a valid keeper session is established")
	}

	fmt.Printf("Searching keeper for records matching %s, please wait...\n", searchTerm)

	// Get password
	configPath, err := keeperConfigPath()
	if err != nil {
		err := errors.New(color.RedString(
			"failed to retrieve keeper config path"))
		return record, err
	}

	output, err := RunCommand("keeper", "search", searchTerm, "--config", configPath)
	if err != nil {
		return record, err
	}

	// Regular expressions to extract relevant information from the output.
	uidRegex := regexp.MustCompile(`UID:\s+(\S+)`)
	titleRegex := regexp.MustCompile(`Title:\s+(.+)`)
	usernameRegex := regexp.MustCompile(`\(login\):\s(.*)`)

	record.UID = uidRegex.FindStringSubmatch(output)[1]
	record.Title = titleRegex.FindStringSubmatch(output)[1]
	record.Username = usernameRegex.FindStringSubmatch(output)[1]

	return record, nil
}
