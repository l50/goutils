package keeper_test

import (
	"os"
	"testing"

	"github.com/l50/goutils/keeper"
	"github.com/stretchr/testify/assert"
)

var testRecord keeper.Record
var note keeper.Record

func init() {
	testRecord.UID = "hfLu-IbhTTVhE3DjWsS-Eg"
	testRecord.Title = "TESTRECORD"
	note.UID = "d2MxKXQpWWhjEPCDz6JKOQ"
}

// Added 1 to test name to ensure this test gets run before all others.
func Test1CommanderInstalled(t *testing.T) {
	// Ensure commander is installed
	assert.True(t, keeper.CommanderInstalled(), "Commander is not installed")
}

// Added 2 to test name to ensure this test gets run before all others (except for Test1).
func Test2LoggedIn(t *testing.T) {
	if os.Getenv("SKIP_KEEPER_TESTS") != "" {
		t.Skip("Skipping test because SKIP_KEEPER_TESTS environment variable is set.")
	}
	// Ensure a valid keeper session exists
	assert.True(t, keeper.LoggedIn(), "no valid keeper session. Please log in to keeper before running tests.")
}

func TestRetrieveRecord(t *testing.T) {
	if os.Getenv("SKIP_KEEPER_TESTS") != "" {
		t.Skip("Skipping test because SKIP_KEEPER_TESTS environment variable is set.")
	}
	// Test case for existing path
	record, err := keeper.RetrieveRecord(testRecord.UID)
	assert.NoError(t, err, "failed to retrieve password")
	assert.Equal(t, "my test password 123!", record.Password, "retrieved password doesn't match expected value")

	// Test case for non-existent path
	_, err = keeper.RetrieveRecord("non/existent/path")
	assert.Error(t, err, "no error for non-existent path")

	// Simulate a non-existent commander installation
	commanderPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	_, err = keeper.RetrieveRecord("my/test/path")
	assert.Error(t, err, "no error when commander is not installed")
	os.Setenv("PATH", commanderPath)

	// Retrieve encryptedNote
	record, err = keeper.RetrieveRecord(note.UID)
	assert.NoError(t, err, "failed to retrieve note")
	assert.Equal(t, "SWEETSECRET!", record.Note, "retrieved note doesn't match expected value")
}

func TestSearchRecords(t *testing.T) {
	var err error

	if os.Getenv("SKIP_KEEPER_TESTS") != "" {
		t.Skip("Skipping test because SKIP_KEEPER_TESTS environment variable is set.")
	}
	_, err = keeper.SearchRecords(testRecord.Title)
	assert.NoError(t, err, "fails to retrieve test record")

	_, err = keeper.SearchRecords("BAMSV")
	assert.NoError(t, err, "does not handle non-matching searchTerm")

	_, err = keeper.SearchRecords("TEST.*RD")
	assert.NoError(t, err, "does not handle regex searchTerm")

	// Simulate a non-existent commander installation
	commanderPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	_, err = keeper.SearchRecords(testRecord.Title)
	assert.Error(t, err, "expected error when commander there is no valid keeper session")
	os.Setenv("PATH", commanderPath)
}
