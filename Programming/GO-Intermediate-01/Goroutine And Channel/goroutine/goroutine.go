package main

import "fmt"

// This code will run Print() which is a goroutine
// because main() function is also a goroutine and main goroutine, others are child goroutine
// So execution of main() function will terminate other child goroutine
func Goroutine() {
	go Print("wrold")
}

func Print(msg string) {
	fmt.Println("Hello ", msg)
}

func main() {
	Goroutine()
}