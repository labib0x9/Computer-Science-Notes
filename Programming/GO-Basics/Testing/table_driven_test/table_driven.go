package main

func Add(a, b int) int {
	return a + b
}

func IsPalindrome(msg string) bool {
	n := len(msg)
	for i := 0; i < n / 2; i++ {
		if msg[i] != msg[n - i - 1] {
			return false
		}
	}
	return true
}