package algo

import (
	"testing"
)

func TestReplaceOnce(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		oldItem  string
		newItem  string
		want     []string
		wantBool bool
	}{
		{
			name:     "replace existing item",
			slice:    []string{"a", "b", "c"},
			oldItem:  "b",
			newItem:  "x",
			want:     []string{"a", "x", "c"},
			wantBool: true,
		},
		{
			name:     "replace first occurrence only",
			slice:    []string{"a", "b", "b", "c"},
			oldItem:  "b",
			newItem:  "x",
			want:     []string{"a", "x", "b", "c"},
			wantBool: true,
		},
		{
			name:     "replace non-existing item",
			slice:    []string{"a", "b", "c"},
			oldItem:  "d",
			newItem:  "x",
			want:     []string{"a", "b", "c"},
			wantBool: false,
		},
		{
			name:     "replace in empty slice",
			slice:    []string{},
			oldItem:  "a",
			newItem:  "x",
			want:     []string{},
			wantBool: false,
		},
		{
			name:     "replace first item",
			slice:    []string{"a", "b", "c"},
			oldItem:  "a",
			newItem:  "x",
			want:     []string{"x", "b", "c"},
			wantBool: true,
		},
		{
			name:     "replace last item",
			slice:    []string{"a", "b", "c"},
			oldItem:  "c",
			newItem:  "x",
			want:     []string{"a", "b", "x"},
			wantBool: true,
		},
		{
			name:     "replace with empty string",
			slice:    []string{"a", "b", "c"},
			oldItem:  "b",
			newItem:  "",
			want:     []string{"a", "", "c"},
			wantBool: true,
		},
		{
			name:     "replace empty string",
			slice:    []string{"a", "", "c"},
			oldItem:  "",
			newItem:  "x",
			want:     []string{"a", "x", "c"},
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy to avoid modifying test data
			sliceCopy := make([]string, len(tt.slice))
			copy(sliceCopy, tt.slice)

			got := ReplaceOnce(sliceCopy, tt.oldItem, tt.newItem)

			if got != tt.wantBool {
				t.Errorf("ReplaceOnce() returned %v, want %v", got, tt.wantBool)
			}

			if len(sliceCopy) != len(tt.want) {
				t.Errorf("ReplaceOnce() resulted in slice length %d, want %d", len(sliceCopy), len(tt.want))
				return
			}

			for i := range sliceCopy {
				if sliceCopy[i] != tt.want[i] {
					t.Errorf("ReplaceOnce() at index %d = %q, want %q", i, sliceCopy[i], tt.want[i])
				}
			}
		})
	}
}

func TestReplaceOnce_Integers(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		oldItem  int
		newItem  int
		want     []int
		wantBool bool
	}{
		{
			name:     "replace integer",
			slice:    []int{1, 2, 3},
			oldItem:  2,
			newItem:  10,
			want:     []int{1, 10, 3},
			wantBool: true,
		},
		{
			name:     "replace with zero",
			slice:    []int{1, 2, 3},
			oldItem:  2,
			newItem:  0,
			want:     []int{1, 0, 3},
			wantBool: true,
		},
		{
			name:     "replace zero",
			slice:    []int{1, 0, 3},
			oldItem:  0,
			newItem:  5,
			want:     []int{1, 5, 3},
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sliceCopy := make([]int, len(tt.slice))
			copy(sliceCopy, tt.slice)

			got := ReplaceOnce(sliceCopy, tt.oldItem, tt.newItem)

			if got != tt.wantBool {
				t.Errorf("ReplaceOnce() returned %v, want %v", got, tt.wantBool)
			}

			for i := range sliceCopy {
				if sliceCopy[i] != tt.want[i] {
					t.Errorf("ReplaceOnce() at index %d = %d, want %d", i, sliceCopy[i], tt.want[i])
				}
			}
		})
	}
}

func TestReplaceOnce_Mutates(t *testing.T) {
	// Test that ReplaceOnce actually mutates the slice
	slice := []string{"a", "b", "c"}
	originalSlice := slice

	ReplaceOnce(slice, "b", "x")

	// Check that the original slice was modified
	if originalSlice[1] != "x" {
		t.Error("ReplaceOnce() should mutate the original slice")
	}
}
