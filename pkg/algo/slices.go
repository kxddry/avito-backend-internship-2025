package algo

// ReplaceOnce replaces the first occurrence of oldItem with newItem in the slice.
func ReplaceOnce[T comparable](slice []T, oldItem, newItem T) bool {
	for i := range slice {
		if slice[i] == oldItem {
			slice[i] = newItem
			return true
		}
	}
	return false
}
