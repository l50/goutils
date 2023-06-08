package str_test

import (
	"reflect"
	"testing"

	str "github.com/l50/goutils/str"
)

func TestGenRandom(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "Valid random string",
			length:  10,
			wantErr: false,
		},
	}

	for _, tc := range tests {
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
	tests := []struct {
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := str.InSlice(tc.strToFind, tc.inputSlice)
			if result != tc.expected {
				t.Errorf("StringInSlice() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := str.ToInt64(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ToInt64() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := str.IsNumeric(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	tests := []struct {
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := str.ToSlice(tc.delimStr, tc.delim)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ToStrSlice() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSlicesEqual(t *testing.T) {
	tests := []struct {
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := str.SlicesEqual(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("SlicesEqual() = %v, want %v", got, tc.want)
			}
		})
	}
}
