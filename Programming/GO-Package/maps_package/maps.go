package main

import "maps"

func main() {
	m1 := make(map[string]int)
	m2 := make(map[string]int)

	maps.Clone(m1)
	maps.Copy(m1, m2)
	maps.Equal(m1, m2)
	maps.Keys(m1)
	maps.Values(m1)
}