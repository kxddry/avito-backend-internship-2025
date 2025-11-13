package algo

// Set is a set of comparable items.
type Set[T comparable] map[T]struct{}

// Add adds the given items to the set.
func (s Set[T]) Add(items ...T) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

// Remove removes the given items from the set.
func (s Set[T]) Remove(items ...T) {
	for _, item := range items {
		delete(s, item)
	}
}

// Has checks if the given item is in the set.
func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}

// SetFrom creates a new set from the given items.
func SetFrom[T comparable](items ...T) Set[T] {
	s := make(Set[T])
	s.Add(items...)
	return s
}
