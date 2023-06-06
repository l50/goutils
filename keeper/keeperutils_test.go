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

func TestCommanderInstalled(t *testing.T) {
	assert.True(t, keeper.CommanderInstalled(), "Commander is not installed")
}

func TestLoggedIn(t *testing.T) {
	if os.Getenv("SKIP_KEEPER_TESTS") != "" {
		t.Skip("Skipping test because SKIP_KEEPER_TESTS environment variable is set.")
	}
	assert.True(t, keeper.LoggedIn(), "no valid keeper session. Please log in to keeper before running tests.")
}

func TestRetrieveRecord(t *testing.T) {
	testCases := []struct {
		name      string
		UID       string
		wantError bool
	}{
		{
			name: "Existing UID",
			UID:  testRecord.UID,
		},
		{
			name:      "Non-Existent UID",
			UID:       "non-existant-UID",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if os.Getenv("SKIP_KEEPER_TESTS") != "" {
				t.Skip("Skipping test because SKIP_KEEPER_TESTS environment variable is set.")
			}

			_, err := keeper.RetrieveRecord(tc.UID)
			if tc.wantError {
				assert.Error(t, err, "Expected error but got none.")
			} else {
				assert.NoError(t, err, "Did not expect error but got one.")
			}
		})
	}
}

func TestSearchRecords(t *testing.T) {
	testCases := []struct {
		name       string
		searchTerm string
		wantError  bool
	}{
		{
			name:       "Matching Search Term",
			searchTerm: testRecord.Title,
		},
		{
			name:       "Non-Matching Search Term",
			searchTerm: "BAMSV",
		},
		{
			name:       "Regex Search Term",
			searchTerm: "TEST.*RD",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if os.Getenv("SKIP_KEEPER_TESTS") != "" {
				t.Skip("Skipping test because SKIP_KEEPER_TESTS environment variable is set.")
			}

			_, err := keeper.SearchRecords(tc.searchTerm)
			if tc.wantError {
				assert.Error(t, err, "Expected error but got none.")
			} else {
				assert.NoError(t, err, "Did not expect error but got one.")
			}
		})
	}
}
