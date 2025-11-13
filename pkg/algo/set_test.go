package algo

import "testing"

func TestSet_Add(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  []string
	}{
		{
			name:  "add single item",
			items: []string{"a"},
			want:  []string{"a"},
		},
		{
			name:  "add multiple items",
			items: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "add duplicate items",
			items: []string{"a", "a", "b"},
			want:  []string{"a", "b"},
		},
		{
			name:  "add empty string",
			items: []string{""},
			want:  []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make(Set[string])
			s.Add(tt.items...)

			if len(s) != len(tt.want) {
				t.Errorf("Set.Add() resulted in size %d, want %d", len(s), len(tt.want))
			}

			for _, item := range tt.want {
				if !s.Has(item) {
					t.Errorf("Set.Add() missing item %q", item)
				}
			}
		})
	}
}

func TestSet_Remove(t *testing.T) {
	tests := []struct {
		name         string
		initialItems []string
		removeItems  []string
		wantRemain   []string
	}{
		{
			name:         "remove existing item",
			initialItems: []string{"a", "b", "c"},
			removeItems:  []string{"b"},
			wantRemain:   []string{"a", "c"},
		},
		{
			name:         "remove multiple items",
			initialItems: []string{"a", "b", "c", "d"},
			removeItems:  []string{"b", "d"},
			wantRemain:   []string{"a", "c"},
		},
		{
			name:         "remove non-existing item",
			initialItems: []string{"a", "b"},
			removeItems:  []string{"c"},
			wantRemain:   []string{"a", "b"},
		},
		{
			name:         "remove all items",
			initialItems: []string{"a", "b"},
			removeItems:  []string{"a", "b"},
			wantRemain:   []string{},
		},
		{
			name:         "remove from empty set",
			initialItems: []string{},
			removeItems:  []string{"a"},
			wantRemain:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make(Set[string])
			s.Add(tt.initialItems...)
			s.Remove(tt.removeItems...)

			if len(s) != len(tt.wantRemain) {
				t.Errorf("Set.Remove() resulted in size %d, want %d", len(s), len(tt.wantRemain))
			}

			for _, item := range tt.wantRemain {
				if !s.Has(item) {
					t.Errorf("Set.Remove() missing item %q", item)
				}
			}

			for _, item := range tt.removeItems {
				wasInInitial := false
				for _, initial := range tt.initialItems {
					if item == initial {
						wasInInitial = true
						break
					}
				}
				if wasInInitial && s.Has(item) {
					t.Errorf("Set.Remove() still contains removed item %q", item)
				}
			}
		})
	}
}

func TestSet_Has(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		check string
		want  bool
	}{
		{
			name:  "has existing item",
			items: []string{"a", "b", "c"},
			check: "b",
			want:  true,
		},
		{
			name:  "has non-existing item",
			items: []string{"a", "b", "c"},
			check: "d",
			want:  false,
		},
		{
			name:  "has in empty set",
			items: []string{},
			check: "a",
			want:  false,
		},
		{
			name:  "has empty string",
			items: []string{"a", ""},
			check: "",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := make(Set[string])
			s.Add(tt.items...)

			if got := s.Has(tt.check); got != tt.want {
				t.Errorf("Set.Has(%q) = %v, want %v", tt.check, got, tt.want)
			}
		})
	}
}

func TestSetFrom(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  []string
	}{
		{
			name:  "create from multiple items",
			items: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "create from no items",
			items: []string{},
			want:  []string{},
		},
		{
			name:  "create from duplicates",
			items: []string{"a", "a", "b"},
			want:  []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SetFrom(tt.items...)

			if len(s) != len(tt.want) {
				t.Errorf("SetFrom() resulted in size %d, want %d", len(s), len(tt.want))
			}

			for _, item := range tt.want {
				if !s.Has(item) {
					t.Errorf("SetFrom() missing item %q", item)
				}
			}
		})
	}
}

func TestSet_WithIntegers(t *testing.T) {
	s := make(Set[int])
	s.Add(1, 2, 3)

	if !s.Has(2) {
		t.Error("Set[int] should contain 2")
	}

	if s.Has(5) {
		t.Error("Set[int] should not contain 5")
	}

	s.Remove(2)
	if s.Has(2) {
		t.Error("Set[int] should not contain 2 after removal")
	}
}
