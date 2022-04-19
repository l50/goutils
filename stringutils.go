package utils

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
)

// StringInSlice determines if an input string exists in an input slice.
// It returns true if the string is found, otherwise it returns false.
func StringInSlice(strToFind string, inputSlice []string) bool {
	for _, value := range inputSlice {
		if strings.Contains(value, strToFind) {
			return true
		}
	}
	return false
}

// StringToInt64 returns the converted int64 value of an input string.
func StringToInt64(value string) (int64, error) {
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, err
	}

	return n, nil
}

// RandomString returns a random string
// of the specified length.
func RandomString(length int) (string, error) {
	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b)[:length], nil
}
