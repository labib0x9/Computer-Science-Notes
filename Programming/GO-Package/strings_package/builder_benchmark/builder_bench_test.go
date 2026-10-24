package main

import (
	"strings"
	"testing"
)

func BenchmarkConcatWithBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		for j := 0; j < 10000; j++ {
			sb.WriteString("Heeloo")
		}
		_ = sb.String()
	}
}

func BenchmarkConcatWithPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < 10000; j++ {
			s += "Heeloo"
		}
	}
}