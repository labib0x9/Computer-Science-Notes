package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Convertion
	n, _ := strconv.Atoi("12")
	fmt.Println(n)

	m := strconv.Itoa(12)
	fmt.Println(m)

	flt, _ := strconv.ParseFloat("123.321", 32)
	flt32 := float32(flt)
	fmt.Println(flt32)

	strconv.ParseInt("12", 2, 64)
	strconv.FormatInt(12, 2)
	strconv.ParseBool("true")
	strconv.FormatFloat(3.1416, 'f', 2, 64)
}
