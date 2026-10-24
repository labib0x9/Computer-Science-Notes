package main

func indexOf[T comparable](arr []T, target T) int {
	for i, v := range arr {
		if v == target {
			return i
		}
	}
	return -1
}

func Keys[K comparable, V any](mp map[K]V) []K {
	var keys []K
	for k := range mp {
		keys = append(keys, k)
	}
	return keys
}

func Filter[V any](arr []V, f func (V) bool) []V {
	var result []V
	for _, v := range arr {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}