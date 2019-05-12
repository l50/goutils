package utils

import (
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	_, err := GetHomeDir()
	if err != nil {
		t.Fatal("Unable to get the user's home directory due to: ", err.Error())
	}
}
