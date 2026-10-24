package main

var a = 10

func add(x, y int) int {
	z := x + y
	return z
}

func main() {

	add(2, 4)
	add(a, 5)

}
