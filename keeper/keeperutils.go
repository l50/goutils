package keeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/l50/goutils/sys"
)

// Record represents a record maintained by Keeper.
type Record struct {
	UID      string
	Title    string
	URL      string
	Username string
	Password string
	TOTP     string
	Note     string
}

// configPath returns the path of the keeper config file.
func configPath() (string, error) {
	home, err := sys.GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".keeper", "config.json"), nil
}

// CommanderInstalled returns true if keeper
// commander is installed on the current system.
func CommanderInstalled() bool {
	return sys.CmdExists("keeper")
}

// LoggedIn returns true if keeper vault
// is logged in with the input `email`.
// Otherwise, it returns false.
func LoggedIn() bool {
	if !CommanderInstalled() {
		err := errors.New(color.RedString(
			"keeper commander is not installed - please install and try again"))
		fmt.Println(err)
		return false
	}

	configPath, err := configPath()
	if err != nil {
		err := errors.New(color.RedString(
			"failed to retrieve keeper config path"))
		fmt.Println(err)
		return false
	}

	loggedIn := "My Vault>"
	out, err := sys.RunCommandWithTimeout(5, "keeper", "shell", "--config", configPath)
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

type rawRecord struct {
	RecordUID string `json:"record_uid"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Fields    []struct {
		Type  string   `json:"type"`
		Value []string `json:"value"`
	} `json:"fields"`
}

// RetrieveRecord returns the record found with the input keeperPath.
func RetrieveRecord(keeperUID string) (Record, error) {
	var record Record

	// Ensure keeper commander is installed and
	// there is a valid keeper session.
	if !CommanderInstalled() || !LoggedIn() {
		return record, errors.New("error: ensure keeper commander is installed and a valid keeper session is established")
	}

	// fmt.Printf("Retrieving record with ID %s from keeper\n", keeperUID)

	configPath, err := configPath()
	if err != nil {
		err := errors.New(color.RedString(
			"failed to retrieve keeper config path"))
		return record, err
	}

	jsonData, err := sys.RunCommand("keeper", "get", keeperUID, "--unmask", "--format", "json", "--config", configPath)
	if err != nil {
		return record, err
	}

	var r rawRecord
	if err := json.Unmarshal([]byte(jsonData), &r); err != nil {
		return record, err
	}

	record.UID = r.RecordUID
	record.Title = r.Title
	for _, field := range r.Fields {
		if len(field.Value) > 0 {
			switch field.Type {
			case "login":
				record.Username = field.Value[0]
			case "password":
				record.Password = field.Value[0]
			case "url":
				record.URL = field.Value[0]
			case "oneTimeCode":
				record.TOTP = field.Value[0]
			case "note":
				record.Note = field.Value[0]
			}
		}
	}

	return record, nil
}

// SearchRecords searches the logged-in user's
// keeper records for records matching the input searchTerm.
//
// The searchTerm can be a string or regex.
//
// Example Inputs:
//
// SearchKeeperRecords("TESTING")
// SearchKeeperRecords("TEST.*RD")
//
// If a searchTerm matches a record, the associated UID is returned.
func SearchRecords(searchTerm string) (string, error) {
	recordUID := ""

	// Ensure keeper commander is installed and
	// there is a valid keeper session.
	if !CommanderInstalled() || !LoggedIn() {
		return recordUID, errors.New("error: ensure keeper commander is installed and a valid keeper session is established")
	}

	fmt.Printf("Searching keeper for records matching %s, please wait...\n", searchTerm)

	configPath, err := configPath()
	if err != nil {
		err := errors.New(color.RedString(
			"failed to retrieve keeper config path"))
		return recordUID, err
	}

	searchResults, err := sys.RunCommand("keeper", "search", searchTerm, "--config", configPath)
	if err != nil {
		return recordUID, err
	}

	regex := regexp.MustCompile(`UID:\s+([a-zA-Z0-9-_]+)`)
	match := regex.FindStringSubmatch(searchResults)

	if len(match) == 2 {
		recordUID = match[1]
	} else {
		fmt.Println("No UID found.")
	}

	return recordUID, nil
}
