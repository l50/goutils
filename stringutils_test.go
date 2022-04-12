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

func TestStringToInt64(t *testing.T) {
	_, err := StringToInt64("65")
	if err != nil {
		t.Fatalf("error running StringToInt64(): %v", err)
	}

	_, err = StringToInt64("chicken")
	if err == nil {
		t.Fatalf("error running StringToInt64(): %v", err)
	}
}
