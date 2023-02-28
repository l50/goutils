package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TESTRECORD = "TESTRECORD"

func init() {
	if os.Getenv("SKIP_TESTS") != "" {
		fmt.Println("Skipping tests because SKIP_TESTS environment variable is set.")
		return
	}

	if !CommanderInstalled() {
		fmt.Println("commander is not installed. Please install commander before running tests.")
		os.Exit(1)
	}

	if !KeeperLoggedIn() {
		fmt.Println("no valid keeper session. Please log in to keeper before running tests.")
		os.Exit(1)
	}
}

func TestCommanderInstalled(t *testing.T) {
	// Ensure commander is installed
	assert.True(t, CommanderInstalled(), "Commander is not installed")
}

func TestRetrieveKeeperPW(t *testing.T) {
	// Test case for existing path
	password, err := RetrieveKeeperPW(TESTRECORD)
	assert.NoError(t, err, "failed to retrieve password")
	assert.Equal(t, "my test password 123!", password, "retrieved password doesn't match expected")

	// Test case for non-existent path
	_, err = RetrieveKeeperPW("non/existent/path")
	assert.Error(t, err, "no error for non-existent path")

	// Simulate not logged in to vault
	if CommanderInstalled() {
		home, _ := GetHomeDir()
		if err := Cd(filepath.Join(home, ".keeper")); err != nil {
			assert.NoError(t, err, "error using Cd() to change into the keeper conf directory")
		}

		os.Remove("config.json")

		_, err = RetrieveKeeperPW(TESTRECORD)
		assert.Error(t, err, "failed to detect missing keeper session")
	}

	// Simulate a non-existent commander installation
	commanderPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	_, err = RetrieveKeeperPW("my/test/path")
	assert.Error(t, err, "no error when commander is not installed")
	os.Setenv("PATH", commanderPath)
}

func TestSearchKeeperRecords(t *testing.T) {
	// Create a temporary file with a test keeper record
	tmpFile, err := os.CreateTemp("", "test_keeper_record_")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("UID: 123\nTitle: test title\nLogin: testusername\nPassword: testpassword"))
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	record, err := SearchKeeperRecords("test title")
	assert.NoError(t, err, "Error in searching keeper records")
	assert.Equal(t, "123", record.UID, "Record UID doesn't match")
	assert.Equal(t, "test title", record.Title, "Record title doesn't match")
	assert.Equal(t, "testusername", record.Username, "Record username doesn't match")
	assert.Equal(t, "testpassword", record.Password, "Record password doesn't match")

	// Simulate a non-existent commander installation
	commanderPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	_, err = SearchKeeperRecords("test title")
	assert.Error(t, err, "No error when commander is not installed")
	os.Setenv("PATH", commanderPath)

	// Simulate not logged in to vault
	if CommanderInstalled() {
		home, _ := GetHomeDir()
		if err := Cd(filepath.Join(home, ".keeper")); err != nil {
			assert.NoError(t, err, "error using Cd() to change into the keeper conf directory")
		}

		os.Remove("config.json")

		_, err = SearchKeeperRecords("test title")
		assert.Error(t, err, "No error when not logged in")
	}
}
