package collection

// Contains checks if a value is present in a slice.
//
// **Parameters:**
//
// slice: Slice to check for the value.
// value: Value to check for in the slice.
//
// **Returns:**
//
// bool: true if the value is present in the slice, false otherwise.
func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
