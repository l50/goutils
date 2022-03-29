package utils

import (
	"os"
	"testing"
)

func TestCreateLogFile(t *testing.T) {
	_, err := CreateLogFile()
	if err != nil {
		t.Fatalf("failed to create log file: %v", err)
	}

	err = os.RemoveAll("logs")
	if err != nil {
		t.Fatalf("unable to delete the logs directory - DeleteFile() failed: %v", err)
	}
}
