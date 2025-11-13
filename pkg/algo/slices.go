package algo

func ReplaceOnce[T comparable](slice []T, oldItem T, newItem T) bool {
	for i := range slice {
		if slice[i] == oldItem {
			slice[i] = newItem
			return true
		}
	}
	return false
}
