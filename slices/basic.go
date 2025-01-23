package slices

func Filter[T comparable](ts []T, keep func(T) bool) []T {
	var result []T
	for i := 0; i < len(ts); i++ {
		if keep(ts[i]) {
			result = append(result, ts[i])
		}
	}
	return result
}

func Map[T, U any](ts []T, f func(T) U) []U {
	result := make([]U, len(ts))
	for i, el := range ts {
		result[i] = f(el)
	}
	return result
}

func Reduce[T, U any](ts []T, accumulator U, f func(U, T) U) U {
	for _, el := range ts {
		accumulator = f(accumulator, el)
	}
	return accumulator
}

func ReduceSame[T any](ts []T, f func(T, T) T) T {
	accumulator := ts[0]
	for i := 1; i < len(ts); i++ {
		accumulator = f(accumulator, ts[i])
	}
	return accumulator
}
