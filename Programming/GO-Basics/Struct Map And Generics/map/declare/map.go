package main

import (
	"fmt"
)

// mp1 variable -> has address
// mp2, mp3 are not variable -> doesn't have address
func main() {
	var mp1 map[string]int	// nil map
	fmt.Println(mp1 == nil)

	mp2 := make(map[string]int)	// empty map
	fmt.Println(mp2 == nil)

	mp3 := map[string]int{}
	fmt.Println(mp3 == nil)

	// Check if a value present in mp1
	if _, err := mp1["Alice"]; err == false {
		fmt.Println("Not Found in mp1")
	}

	delete(mp2, "alice")


	for key, value := range mp3 {
		fmt.Println(key, value)
	}


	// Set
	st := make(map[string]struct{})
	st["a"] = struct{}{}

	clear(st)
	
	st = map[string]struct{}{}
}
