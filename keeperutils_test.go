package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestRecord = "TESTRECORD"

func init() {
	if os.Getenv("SKIP_KEEPER_TESTS") != "" {
		fmt.Println("Skipping tests because SKIP_KEEPER_TESTS environment variable is set.")
		return
	}
}

// Added 1 to test name to ensure this test gets run before all others.
func Test1CommanderInstalled(t *testing.T) {
	// Ensure commander is installed
	assert.True(t, CommanderInstalled(), "Commander is not installed")
}

// Added 2 to test name to ensure this test gets run before all others (except for Test1).
func Test2KeeperLoggedIn(t *testing.T) {
	// Ensure a valid keeper session exists
	assert.True(t, KeeperLoggedIn(), "no valid keeper session. Please log in to keeper before running tests.")
}

func TestRetrieveKeeperPW(t *testing.T) {
	// Test case for existing path
	password, err := RetrieveKeeperPW(TestRecord)
	assert.NoError(t, err, "failed to retrieve password")
	assert.Equal(t, "my test password 123!", password, "retrieved password doesn't match expected")

	// Test case for non-existent path
	_, err = RetrieveKeeperPW("non/existent/path")
	assert.Error(t, err, "no error for non-existent path")

	// Simulate a non-existent commander installation
	commanderPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	_, err = RetrieveKeeperPW("my/test/path")
	assert.Error(t, err, "no error when commander is not installed")
	os.Setenv("PATH", commanderPath)
}

func TestSearchKeeperRecords(t *testing.T) {
	_, err := SearchKeeperRecords(TestRecord)
	assert.NoError(t, err, "failed to retrieve test record")

	// Simulate a non-existent commander installation
	commanderPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	_, err = SearchKeeperRecords(TestRecord)
	assert.Error(t, err, "expected error when commander there is no valid keeper session")
	os.Setenv("PATH", commanderPath)
}
