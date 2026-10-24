package main

import "fmt"

func Pass() {
	mp := make(map[string]int)

	// length exceed
	// pass by reference
	function(mp)
	fmt.Println(len(mp))

	// // Error
	// // map returns values, not pointers
	// print(&mp["STR"])
}

func function(mp map[string]int) {
	mp["STR"] = 123
}

func print(n *int) {
	fmt.Println(n)
}
