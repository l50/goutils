package ansible_test

import (
	"os"
	"testing"

	"github.com/l50/goutils/v2/ansible"
)

func TestAnsiblePing(t *testing.T) {
	testCases := []struct {
		name      string
		hostsFile string
		wantErr   bool
	}{
		{
			name:      "valid hosts file",
			hostsFile: "test_inventory.ini",
			wantErr:   false,
		},
		{
			name:      "missing hosts file",
			hostsFile: "missing_hosts.ini",
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a dummy hosts file for the valid case
			if !tc.wantErr {
				dummyData := []byte("[localhost]\n127.0.0.1 ansible_connection=local")
				err := os.WriteFile(tc.hostsFile, dummyData, 0644)
				if err != nil {
					t.Fatalf("Could not create hosts file: %v", err)
				}
				// Clean up after the test
				defer os.Remove(tc.hostsFile)
			}

			// Run the Ping function
			if err := ansible.Ping(tc.hostsFile); (err != nil) != tc.wantErr {
				t.Errorf("Ping(%v) error = %v, wantErr %v", tc.hostsFile, err, tc.wantErr)
			}
		})
	}
}
