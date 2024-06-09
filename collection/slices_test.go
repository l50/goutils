package collection_test

import (
	"testing"

	"github.com/l50/goutils/v2/collection"
)

func TestContains(t *testing.T) {
	testCases := []struct {
		name  string
		slice []int
		value int
		want  bool
	}{
		{
			name:  "element in slice",
			slice: []int{1, 2, 3, 4, 5},
			value: 3,
			want:  true,
		},
		{
			name:  "element not in slice",
			slice: []int{1, 2, 3, 4, 5},
			value: 6,
			want:  false,
		},
		{
			name:  "empty slice",
			slice: []int{},
			value: 1,
			want:  false,
		},
		{
			name:  "nil slice",
			slice: nil,
			value: 1,
			want:  false,
		},
		{
			name:  "multiple occurrences",
			slice: []int{1, 2, 3, 3, 4},
			value: 3,
			want:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := collection.Contains(tc.slice, tc.value)
			if got != tc.want {
				t.Errorf("Contains(%v, %v) = %v; want %v", tc.slice,
					tc.value, got, tc.want)
			}
		})
	}
}
