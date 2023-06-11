package keeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/l50/goutils/v2/sys"
)

// Record represents a user's record in the Keeper application.
// Each record includes a unique identifier (UID), title, URL,
// username, password, TOTP, and a note.
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

// CommanderInstalled checks if the Keeper Commander tool is installed on the current system.
//
// Returns:
//
// bool: True if the Keeper Commander tool is installed, false otherwise.
//
// Example:
//
//	if !CommanderInstalled() {
//	  log.Fatal("Keeper Commander is not installed.")
//	}
func CommanderInstalled() bool {
	return sys.CmdExists("keeper")
}

// LoggedIn checks if the user is logged into their Keeper vault.
//
// Returns:
//
// bool: True if the user is logged into their Keeper vault, false otherwise.
//
// Example:
//
//	if !LoggedIn() {
//	  log.Fatal("Not logged into Keeper vault.")
//	}
func LoggedIn() bool {
	if !CommanderInstalled() {
		err := errors.New(color.RedString(
			"keeper commander is not installed - please install and try again"))
		fmt.Println(err)
		return false
	}

	configPath, err := configPath()
	if err != nil {
		err := errors.New("failed to retrieve keeper config path")
		fmt.Println(err)
		return false
	}

	loggedIn := "My Vault>"
	out, err := sys.RunCommandWithTimeout(15, "keeper", "shell", "--config", configPath)
	if err != nil {
		fmt.Print("failed to check login state "+
			"for keeper vault: ", err)
		return false
	}

	if strings.Contains(string(out), loggedIn) {
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

// RetrieveRecord retrieves a user's Keeper record using the provided unique identifier (keeperUID).
//
// Parameters:
//
// keeperUID: A string representing the unique identifier of the Keeper record to retrieve.
//
// Returns:
//
// Record: The retrieved Keeper record.
// error: An error if the Keeper record cannot be retrieved.
//
// Example:
//
// record, err := RetrieveRecord("1234abcd")
//
//	if err != nil {
//	  log.Fatalf("Failed to retrieve record: %v", err)
//	}
//
// log.Printf("Retrieved record: %+v\n", record)
func RetrieveRecord(keeperUID string) (Record, error) {
	var record Record

	// Ensure keeper commander is installed and
	// there is a valid keeper session.
	if !CommanderInstalled() || !LoggedIn() {
		return record, errors.New("error: ensure keeper commander is installed and a valid keeper session is established")
	}

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

// SearchRecords searches the user's Keeper records for records that match the provided search term.
//
// Parameters:
//
// searchTerm: A string representing the term to search for in the Keeper records.
//
// Returns:
//
// string: The unique identifier (UID) of the first Keeper record that matches the search term.
// error: An error if the Keeper records cannot be searched or if the search term does not match any records.
//
// Example:
//
// uid, err := SearchRecords("search term")
//
//	if err != nil {
//	  log.Fatalf("Failed to search records: %v", err)
//	}
//
// log.Printf("Found matching record with UID: %s\n", uid)
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
