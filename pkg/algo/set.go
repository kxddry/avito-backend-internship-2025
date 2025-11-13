package algo

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(items ...T) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

func (s Set[T]) Remove(items ...T) {
	for _, item := range items {
		delete(s, item)
	}
}

func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}

func SetFrom[T comparable](items ...T) Set[T] {
	s := make(Set[T])
	s.Add(items...)
	return s
}
