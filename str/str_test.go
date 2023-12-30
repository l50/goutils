package str_test

import (
	"reflect"
	"testing"

	str "github.com/l50/goutils/v2/str"
)

func TestGenRandom(t *testing.T) {
	testCases := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "Valid random string",
			length:  10,
			wantErr: false,
		},
		{
			name:    "Zero length",
			length:  0,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := str.GenRandom(tc.length)
			if (err != nil) != tc.wantErr {
				t.Errorf("GenRandom() error = %v, wantErr %v", err, tc.wantErr)
			}

			if len(got) != tc.length {
				t.Errorf("GenRandom() length = %v, want %v", len(got), tc.length)
			}
		})
	}
}

func TestInSlice(t *testing.T) {
	testCases := []struct {
		name       string
		strToFind  string
		inputSlice []string
		expected   bool
	}{
		{
			name:       "Find existing string in slice",
			strToFind:  "sky",
			inputSlice: []string{"sky", "falcon", "rock", "hawk"},
			expected:   true,
		},
		{
			name:       "String not found in slice",
			strToFind:  "cloud",
			inputSlice: []string{"sky", "falcon", "rock", "hawk"},
			expected:   false,
		},
		{
			name:       "Empty slice",
			strToFind:  "sky",
			inputSlice: []string{},
			expected:   false,
		},
		{
			name:       "Partial string match in slice",
			strToFind:  "hawk",
			inputSlice: []string{"skyhawk", "falcon", "rock"},
			expected:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := str.InSlice(tc.strToFind, tc.inputSlice)
			if result != tc.expected {
				t.Errorf("InSlice() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid string to int64",
			input:   "65",
			wantErr: false,
		},
		{
			name:    "Invalid string to int64",
			input:   "chicken",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := str.ToInt64(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ToInt64() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "All numeric",
			input:    "1234",
			expected: true,
		},
		{
			name:     "Alphanumeric",
			input:    "1234abc",
			expected: false,
		},
		{
			name:     "All alphabetic",
			input:    "abcd",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := str.IsNumeric(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	testCases := []struct {
		name     string
		delimStr string
		delim    string
		want     []string
	}{
		{
			name:     "String to slice",
			delimStr: "asasdf\nasdf\nb\ndsfsdf,c",
			delim:    "\n",
			want:     []string{"asasdf", "asdf", "b", "dsfsdf,c"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := str.ToSlice(tc.delimStr, tc.delim)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ToStrSlice() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSlicesEqual(t *testing.T) {
	testCases := []struct {
		name string
		a    []string
		b    []string
		want bool
	}{
		{
			name: "Equal slices",
			a:    []string{"apple", "banana", "cherry"},
			b:    []string{"apple", "banana", "cherry"},
			want: true,
		},
		{
			name: "Unequal slices",
			a:    []string{"apple", "banana", "cherry"},
			b:    []string{"apple", "banana", "grape"},
			want: false,
		},
		{
			name: "Different length slices",
			a:    []string{"apple", "banana", "cherry"},
			b:    []string{"apple", "banana"},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := str.SlicesEqual(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("SlicesEqual() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	testCases := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "No ANSI codes",
			str:  "Hello, world!",
			want: "Hello, world!",
		},
		{
			name: "Single ANSI code",
			str:  "\x1B[31mHello, world!\x1B[0m",
			want: "Hello, world!",
		},
		{
			name: "Multiple ANSI codes",
			str:  "\x1B[1m\x1B[34mBold and blue\x1B[0m text",
			want: "Bold and blue text",
		},
		{
			name: "Nested ANSI codes",
			str:  "Normal \x1B[32mGreen \x1B[1mBold green\x1B[0m Normal",
			want: "Normal Green Bold green Normal",
		},
		{
			name: "Incomplete ANSI code",
			str:  "Hello \x1B[33mYellow",
			want: "Hello Yellow",
		},
		{
			name: "Only ANSI codes",
			str:  "\x1B[4m\x1B[45m",
			want: "",
		},
		{
			name: "Complex ANSI codes",
			str:  "\u001b[34m\u001b[0;32m    docker.ansible-attack-box: The following packages will be upgraded:\u001b[0m\n\u001b[0m",
			want: "    docker.ansible-attack-box: The following packages will be upgraded:\n",
		},
		{
			name: "ANSI codes with text",
			str:  "\u001b[34m\u001b[0;32m    docker.ansible-attack-box:   bash python3-pip\u001b[0m\n\u001b[0m",
			want: "    docker.ansible-attack-box:   bash python3-pip\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := str.StripANSI(tc.str)
			if got != tc.want {
				t.Errorf("StripANSI(%q) = %q, want %q", tc.str, got, tc.want)
			}
		})
	}
}
