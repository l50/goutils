package keeper_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/l50/goutils/v2/pwmgr"
	"github.com/l50/goutils/v2/pwmgr/keeper"
	"github.com/stretchr/testify/assert"
)

var testRecord pwmgr.Record
var note pwmgr.Record

func init() {
	testRecord.UID = "hfLu-IbhTTVhE3DjWsS-Eg"
	testRecord.Title = "TESTRECORD"
	note.UID = "d2MxKXQpWWhjEPCDz6JKOQ"
}

func TestCommanderInstalled(t *testing.T) {
	k := keeper.Keeper{}
	assert.True(t, k.CommanderInstalled(), "Commander is not installed")
}

func TestLoggedIn(t *testing.T) {
	if os.Getenv("SKIP_KEEPER_TESTS") != "" {
		t.Skip("Skipping test because SKIP_KEEPER_TESTS environment variable is set.")
	}

	k := keeper.Keeper{}
	assert.True(t, k.LoggedIn(), "no valid keeper session. Please log in to keeper before running tests.")
}

func TestRetrieveRecord(t *testing.T) {
	testCases := []struct {
		name      string
		UID       string
		wantError bool
	}{
		{
			name: "Existing record",
			UID:  testRecord.UID,
		},
		{
			name:      "Non-existent record",
			UID:       "Non-existent UID",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if os.Getenv("SKIP_KEEPER_TESTS") != "" {
				t.Skip("Skipping test because SKIP_KEEPER_TESTS environment variable is set.")
			}

			k := keeper.Keeper{}
			jsonRecord, err := k.RetrieveRecord(tc.UID)
			if tc.wantError {
				assert.Error(t, err, "Expected error but got none.")
			} else {
				assert.NoError(t, err, "Did not expect error but got one.")

				// Assert that the returned string is a valid JSON representation of a pwmgr.Record
				var record pwmgr.Record
				err := json.Unmarshal([]byte(jsonRecord), &record)
				assert.NoError(t, err, "Failed to unmarshal returned JSON into a pwmgr.Record")

				if !tc.wantError {
					// Assert specific attributes of the record
					assert.Equal(t, "hfLu-IbhTTVhE3DjWsS-Eg", record.UID, "UID mismatch")
					assert.Equal(t, "TESTRECORD", record.Title, "Title mismatch")
					assert.Equal(t, "https://evil.com", record.URL, "URL mismatch")
					assert.Equal(t, "test", record.Username, "Username mismatch")
					assert.Equal(t, "my test password 123!", record.Password, "Password mismatch")
					assert.Equal(t, "", record.TOTP, "TOTP mismatch")
					assert.Equal(t, "", record.Note, "Note mismatch")
				}
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

			k := keeper.Keeper{}
			_, err := k.SearchRecords(tc.searchTerm)
			if tc.wantError {
				assert.Error(t, err, "Expected error but got none.")
			} else {
				assert.NoError(t, err, "Did not expect error but got one.")
			}
		})
	}
}
