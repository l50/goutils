package str

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
)

// GenRandom generates a random string of the specified length.
//
// Parameters:
//
// length: The length of the random string to be generated.
//
// Returns:
//
// string: The generated random string.
// error: An error if the random string generation fails.
//
// Example:
//
// str, err := GenRandom(10)
//
//	if err != nil {
//	  log.Fatalf("failed to generate random string: %v", err)
//	}
//
// log.Printf("Generated random string: %s\n", str)
func GenRandom(length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b)[:length], nil
}

// InSlice determines if a specified string exists in a provided slice.
//
// Parameters:
//
// strToFind: The string to search for in the slice.
// inputSlice: The slice of strings to be searched.
//
// Returns:
//
// bool: true if the string is found in the slice, false otherwise.
//
// Example:
//
// slice := []string{"apple", "banana", "cherry"}
// isFound := InSlice("banana", slice)
//
//	if isFound {
//	  log.Println("Found the string in the slice.")
//	}
func InSlice(strToFind string, inputSlice []string) bool {
	for _, value := range inputSlice {
		if strings.Contains(value, strToFind) {
			return true
		}
	}
	return false
}

// ToInt64 converts a string to an int64.
//
// Parameters:
//
// value: The string to be converted to int64.
//
// Returns:
//
// int64: The int64 equivalent of the string.
// error: An error if the string to int64 conversion fails.
//
// Example:
//
// num, err := ToInt64("1234567890")
//
//	if err != nil {
//	  log.Fatalf("failed to convert string to int64: %v", err)
//	}
//
// log.Printf("Converted string to int64: %d\n", num)
func ToInt64(value string) (int64, error) {
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, err
	}

	return n, nil
}

// ToSlice converts a string to a slice of strings, using the specified delimiter.
//
// Parameters:
//
// delimStr: The string to be split into a slice.
// delim: The delimiter to be used for splitting the string.
//
// Returns:
//
// []string: A slice of strings obtained by splitting the input string.
//
// Example:
//
// slice := ToSlice("apple,banana,cherry", ",")
//
//	for _, str := range slice {
//	  log.Println(str)
//	}
func ToSlice(delimStr string, delim string) []string {
	return strings.Split(delimStr, delim)
}

// SlicesEqual compares two slices of strings for equality.
//
// Parameters:
//
// a: The first string slice to be compared.
// b: The second string slice to be compared.
//
// Returns:
//
// bool: true if the slices are equal (same length and same values), false otherwise.
//
// Example:
//
// a := []string{"apple", "banana", "cherry"}
// b := []string{"apple", "banana", "cherry"}
// isEqual := SlicesEqual(a, b)
//
//	if isEqual {
//	  log.Println("The string slices are equal.")
//	} else {
//
//	  log.Println("The string slices are not equal.")
//	}
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
