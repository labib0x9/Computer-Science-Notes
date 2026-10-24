package main

import (
	"fmt"
	"strconv"
	"unicode"
)

func ifElse() {

	n := 100
	if n+n >= 200 {
		//
	} else if n >= 200 {
		//
	} else {
		//
	}

	if num, err := strconv.Atoi("324"); err != nil {
		fmt.Println("number: ", num)
	}
}

func switchCase() {
	cmd := "/"

	switch cmd {
	case "+":
		//
	case "-":
		//
	case "/":
		//
		fallthrough
	case "%":
		//
	default:
		//
	}

	r := rune('1')
	// r = rune('Ⅻ')
	switch {
	case unicode.IsLetter(r):
		fmt.Println("Rune is a letter")
	case unicode.IsDigit(r):
		fmt.Println("Rune is a digit")
		fallthrough
	case unicode.IsNumber(r):
		fmt.Println("Rune is a number")
	}

	words := []string{"go", "Hello", "WorldIsInDanger"}
	for _, word := range words {
		switch size := len(word); size {
		case 1, 2, 3:
			fmt.Println("Small")
		case 4, 5, 6:
			fmt.Println("Medium")
		default:
			fmt.Println("Large")
		}
	}
}

func switchInsideLoop() {
	// didn't exit loop
	for i := 0; i < 10; i++ {
		switch i {
		case 0, 2, 4, 6, 8:
			fmt.Println(i, "Even")
		case 7:
			fmt.Println(i, "Exit loop")
			break
		default:
			fmt.Println(i, "Odd")
		}
	}

	// Exit loop
	loop:
	for i := 0; i < 10; i++ {
		switch i {
		case 0, 2, 4, 6, 8:
			fmt.Println(i, "Even")
		case 7:
			fmt.Println(i, "Exit loop")
			break loop
		default:
			fmt.Println(i, "Odd")
		}
	}
}

func main() {

	switchCase()
	switchInsideLoop()

}
