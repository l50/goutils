package str_test

import (
	"fmt"
	"log"

	"github.com/l50/goutils/v2/str"
)

func ExampleGenRandom() {
	randStr, err := str.GenRandom(10)
	if err != nil {
		log.Fatalf("failed to generate random string: %v", err)
	}
	fmt.Printf("Generated random string: %s\n", randStr)
}

func ExampleInSlice() {
	slice := []string{"apple", "banana", "cherry"}
	isFound := str.InSlice("banana", slice)
	fmt.Println(isFound)
	// Output: true
}

func ExampleIsNumeric() {
	isNum := str.IsNumeric("1234")
	fmt.Println(isNum)
	// Output: true
}

func ExampleToInt64() {
	num, err := str.ToInt64("1234567890")
	if err != nil {
		log.Fatalf("failed to convert string to int64: %v", err)
	}
	fmt.Printf("Converted string to int64: %d\n", num)
	// Output: Converted string to int64: 1234567890
}

func ExampleToSlice() {
	slice := str.ToSlice("apple,banana,cherry", ",")
	fmt.Println(slice)
	// Output: [apple banana cherry]
}

func ExampleSlicesEqual() {
	a := []string{"apple", "banana", "cherry"}
	b := []string{"apple", "banana", "cherry"}
	isEqual := str.SlicesEqual(a, b)
	fmt.Println(isEqual)
	// Output: true
}
