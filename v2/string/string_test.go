package string_test

import (
	"reflect"
	"testing"

	stringutils "github.com/l50/goutils/v2/string"
)

func TestRandomString(t *testing.T) {
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
			got, err := stringutils.RandomString(tc.length)
			if (err != nil) != tc.wantErr {
				t.Errorf("RandomString() error = %v, wantErr %v", err, tc.wantErr)
			}

			if len(got) != tc.length {
				t.Errorf("RandomString() length = %v, want %v", len(got), tc.length)
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
			result := stringutils.InSlice(tc.strToFind, tc.inputSlice)
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
			_, err := stringutils.ToInt64(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("StringToInt64() error = %v, wantErr %v", err, tc.wantErr)
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
			got := stringutils.ToSlice(tc.delimStr, tc.delim)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("StringToSlice() = %v, want %v", got, tc.want)
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
			got := stringutils.SlicesEqual(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("SlicesEqual() = %v, want %v", got, tc.want)
			}
		})
	}
}
