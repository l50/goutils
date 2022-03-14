package utils

import (
	"testing"
)

func TestStringInSlice(t *testing.T) {
	words := []string{"sky", "falcon", "rock", "hawk"}
	if !StringInSlice("sky", words) {
		t.Fatal("unable to find a string that exists in the test slice - StringInSlice() failed")
	}
}
