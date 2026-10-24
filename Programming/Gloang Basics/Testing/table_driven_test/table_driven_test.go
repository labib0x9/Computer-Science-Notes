package main

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {

	testCasesInt := []struct {
		a, b, sum int64
	}{
		{10, 20, 30},
		{40, 43, 83},

		// One interesting thing ...
		// By defination of int and int64 it would fail. but passed the test. Why ??
		// 32-bit architect machine : int works like int32
		// 64-bit architext machine : int works like int64
		{1000001212, 1212, 1000002424},
	}

	for i, test := range testCasesInt {
		t.Run(fmt.Sprintf("Int : Case-%d", i), func(t *testing.T) {
			expected := test.sum
			got := Add(int(test.a), int(test.b))
			if expected != int64(got) {
				t.Errorf("Add(%d, %d) = %d, want %d", test.a, test.b, got, expected)
			}
		})
	}

	testCasesFloat := []struct {
		a, b, sum float32
	}{
		{10.5, 20.5, 31},
		{0, 1, 1},
		{12.2, 12.12, 24.22},
	}

	for i, test := range testCasesFloat {
		t.Run(fmt.Sprintf("Float : Case-%d", i), func(t *testing.T) {
			expected := test.sum
			got := Add(int(test.a), int(test.b))
			if expected != float32(got) {
				t.Errorf("Add(%v, %v) = %v, want %v", test.a, test.b, got, expected)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	testCase := []struct {
		msg  string
		flag bool
	}{
		{"aba", true},
		{"ba", false},
		{string([]rune{170, 66, 170}), true},
		{"Aba", true},
		{"a man, nama", true},
	}

	for _, test := range testCase {
		got := IsPalindrome(test.msg)
		if got != test.flag {
			t.Errorf("IsPalindrome(\"%v\") = %v, want %v", test.msg, got, test.flag)
		}
	}
}
