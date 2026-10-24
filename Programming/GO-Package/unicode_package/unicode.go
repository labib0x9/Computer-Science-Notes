package main

import "unicode"

func main() {
	unicode.Is(unicode.Letter, 'A')
	unicode.Is(unicode.Digit, '4')

	unicode.IsDigit('1')
	unicode.IsLetter('r')
	unicode.IsSpace('\n')
	unicode.IsSymbol(',')
	unicode.IsPunct('/')

	unicode.ToLower('R')
	unicode.ToUpper('r')
	unicode.ToTitle('r')
}