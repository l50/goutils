package utils

import (
	"reflect"
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

func TestStringToSlice(t *testing.T) {
	delimStr := "asasdf\nasdf\nb\ndsfsdf,c"
	delim := "\n"
	ss := StringToSlice(delimStr, delim)
	if reflect.TypeOf(ss).Elem().Kind() != reflect.String {
		t.Fatalf(
			"error - failed to convert %s using %s delimiter to a string slice",
			delimStr, delim)
	}
}

func TestRandomString(t *testing.T) {
	length := 10
	randStr, err := RandomString(10)
	if err != nil {
		t.Fatalf("error creating random string of length %d: %v", length, err)
	}

	if len(randStr) != length {
		t.Fatalf("length of the random string does not match the input length %d", length)
	}
}
