package main

import (
	"fmt"
	"reflect"
)

// printTypeAndValue prints the type and value of the given interface{}.
// v is a rune that determines the format of the output.
// If v is 'T', it prints the type only.
// If v is 'V', it prints the value only.
// If v is 'B', it prints both the type and value.
func printTypeAndValue(x interface{}, verb rune) {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	switch verb {
	case 'T':
		fmt.Println("Type:", t)
	case 'V':
		fmt.Println("Value:", v)
	case 'B':
		fmt.Println("Type:", t, ", Value:", v)
	default:
		fmt.Println("Invalid format specifier")
	}
}

func main() {

	printTypeAndValue(42, 'T')
	printTypeAndValue(42, 'V')
	printTypeAndValue(42, 'B')

	printTypeAndValue("Hello, World!", 'T')
	printTypeAndValue("Hello, World!", 'V')
	printTypeAndValue("Hello, World!", 'B')

	printTypeAndValue(3.14, 'T')
	printTypeAndValue(3.14, 'V')
	printTypeAndValue(3.14, 'B')
}
