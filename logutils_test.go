package utils

import (
	"testing"
)

func TestCreateLogFile(t *testing.T) {
	if _, err := CreateLogFile(); err != nil {
		t.Fatalf("failed to create log file: %v", err)
	}

	// Remove the temporary file after the test completes.
	defer func() {
		dir := "./logs"
		if err := RmRf(dir); err != nil {
			t.Fatalf("unable to delete %s, RmRf() failed", dir)
		}
	}()
}
