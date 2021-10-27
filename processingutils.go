package utils

// StringInSlice determines if an input string exists in an input slice.
// It returns true if the string is found, otherwise it returns false.
func StringInSlice(strToFind string, inputSlice []string) bool {
	for _, value := range inputSlice {
		if value == strToFind {
			return true
		}
	}
	return false
}
