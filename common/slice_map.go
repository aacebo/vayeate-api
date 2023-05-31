package common

func SliceMap[T any, R any](s []T, fn func(T) R) []R {
	res := make([]R, len(s))

	for i := 0; i < len(s); i++ {
		res[i] = fn(s[i])
	}

	return res
}
