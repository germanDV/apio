package tools

// Map maps over a slice and returns a new slice applying `fn` to every item.
func Map[F any, T any](from []F, fn func(F) T) []T {
	to := make([]T, 0, len(from))
	for _, v := range from {
		to = append(to, fn(v))
	}
	return to
}
