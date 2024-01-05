package str

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// GenRandom generates a random string of a specified length.
//
// **Parameters:**
//
// length: Length of the random string to be generated.
//
// **Returns:**
//
// string: Generated random string.
// error: An error if random string generation fails.
func GenRandom(length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b)[:length], nil
}

// InSlice determines if a specified string exists in a given slice.
//
// **Parameters:**
//
// strToFind: String to search for in the slice.
// inputSlice: Slice of strings to be searched.
//
// **Returns:**
//
// bool: true if string is found in the slice, false otherwise.
func InSlice(strToFind string, inputSlice []string) bool {
	for _, value := range inputSlice {
		if strings.Contains(value, strToFind) {
			return true
		}
	}
	return false
}

// IsNumeric checks if a string is entirely composed of numeric characters.
//
// **Parameters:**
//
// s: String to check for numeric characters.
//
// **Returns:**
//
// bool: true if the string is numeric, false otherwise.
func IsNumeric(s string) bool {
	for _, char := range s {
		if _, err := strconv.Atoi(string(char)); err != nil {
			return false
		}
	}
	return true
}

// ToInt64 converts a string to int64.
//
// **Parameters:**
//
// value: String to be converted to int64.
//
// **Returns:**
//
// int64: int64 equivalent of the string.
// error: An error if the conversion fails.
func ToInt64(value string) (int64, error) {
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, err
	}

	return n, nil
}

// ToSlice converts a string to a slice of strings using a delimiter.
//
// **Parameters:**
//
// delimStr: String to be split into a slice.
// delim: Delimiter to be used for splitting the string.
//
// **Returns:**
//
// []string: Slice of strings from the split input string.
func ToSlice(delimStr string, delim string) []string {
	return strings.Split(delimStr, delim)
}

// SlicesEqual compares two slices of strings for equality.
//
// **Parameters:**
//
// a: First string slice for comparison.
// b: Second string slice for comparison.
//
// **Returns:**
//
// bool: true if slices are equal, false otherwise.
func SlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// StripANSI removes ANSI escape codes from a string.
//
// **Parameters:**
//
// str: String to remove ANSI escape codes from.
//
// **Returns:**
//
// string: String with ANSI escape codes removed.
func StripANSI(str string) string {
	re := regexp.MustCompile(`\x1B\[[0-9;]*[a-zA-Z]`)
	return re.ReplaceAllString(str, "")
}
