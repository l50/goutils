package keeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/l50/goutils/v2/pwmgr"
	"github.com/l50/goutils/v2/sys"
)

// Keeper represents a connection with the Keeper password manager.
type Keeper struct{}

// configPath returns the path of the keeper config file.
func configPath() (string, error) {
	home, err := sys.GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".keeper", "config.json"), nil
}

// CommanderInstalled checks if the Keeper Commander tool is
// installed on the current system.
//
// **Returns:**
//
// bool: True if the Keeper Commander tool is installed, false otherwise.
func (k Keeper) CommanderInstalled() bool {
	return sys.CmdExists("keeper")
}

// LoggedIn checks if the user is logged into their Keeper vault.
//
// **Returns:**
//
// bool: True if the user is logged into their Keeper vault, false otherwise.
func (k Keeper) LoggedIn() bool {
	if !k.CommanderInstalled() {
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

// AddRecord adds a new record to the Keeper vault.
//
// **Parameters:**
//
// fields: A map containing the record fields.
//
// fields.title: The title of the record.
// fields.login: The username or login of the record.
// fields.password: The password of the record.
// fields.notes: Additional notes related to the record.
//
// **Returns:**
//
// error: An error if the record cannot be added.
func (k Keeper) AddRecord(fields map[string]string) error {
	// Ensure keeper commander is installed and
	// there is a valid keeper session.
	if !k.CommanderInstalled() || !k.LoggedIn() {
		return errors.New("error: ensure keeper commander is installed " +
			"and a valid keeper session is established")
	}

	configPath, err := configPath()
	if err != nil {
		return errors.New(color.RedString(
			"failed to retrieve keeper config path"))
	}

	title := fields["title"]
	login := fields["login"]
	password := fields["password"]
	notes := fields["notes"]

	_, err = sys.RunCommand("keeper", "record-add", "--title", title, "--login", login, "--pass", password, "--notes", notes, "--config", configPath)

	return err
}

// RetrieveRecord retrieves a user's Keeper record using the
// provided unique identifier (uid).
//
// **Parameters:**
//
// uid: A string representing the unique identifier of the
// Keeper record to retrieve.
//
// **Returns:**
//
// pwmgr.Record: The retrieved Keeper record. This contains the following
// attributes:
//
// - UID: The unique identifier of the record.
// - Title: The title of the record.
// - Username: The username associated with the record.
// - Password: The password of the record.
// - URL: The URL associated with the record.
// - TOTP: The one-time password (if any) associated with the record.
// - Note: Any additional notes associated with the record.
//
// error: An error if the Keeper record cannot be retrieved.
func (k Keeper) RetrieveRecord(uid string) (pwmgr.Record, error) {
	var record pwmgr.Record

	if !k.CommanderInstalled() || !k.LoggedIn() {
		return record, errors.New("error: ensure keeper commander is installed and a valid keeper session is established")
	}

	configPath, err := configPath()
	if err != nil {
		err := errors.New(color.RedString(
			"failed to retrieve keeper config path"))
		return record, err
	}

	jsonData, err := sys.RunCommand("keeper", "get", uid, "--unmask", "--format", "json", "--config", configPath)
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

// SearchRecords searches the user's Keeper records for records
// that match the provided search term.
//
// **Parameters:**
//
// searchTerm: A string representing the term to search for in the Keeper records.
//
// **Returns:**
//
// string: The unique identifier (UID) of the first Keeper record
// that matches the search term. If multiple records match the
// search term, only the UID of the first record is returned.
//
// error: An error if the Keeper records cannot be searched or if
// the search term does not match any records.
func (k Keeper) SearchRecords(searchTerm string) (string, error) {
	recordUID := ""

	if !k.CommanderInstalled() || !k.LoggedIn() {
		return recordUID, errors.New("error: ensure keeper commander is installed and a valid keeper session is established")
	}

	fmt.Printf("searching keeper for records matching %s, please wait...\n", searchTerm)

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
		fmt.Println("no UID found.")
	}

	return recordUID, nil
}
