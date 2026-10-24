package main

/**
	go build -gcflags="-m" main.go
**/

var a = 10

func add(x, y int) int {
	z := x + y
	// fmt.Println(z) // z escapes to heap, because of Println()
	return z
}

func main() {

	add(2, 4)
	add(a, 5)

}
