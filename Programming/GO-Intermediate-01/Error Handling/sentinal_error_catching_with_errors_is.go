package main

import (
	"errors"
	"fmt"
)

var Err = errors.New("error")

func main() {
	e := fmt.Errorf("Error Cause: %w", Err)

	// should catch the sentinal error
	if errors.Is(e, Err) {
		fmt.Println("YEAP, I GOT IT")
	} else {
		fmt.Println("COULDNOT CATCH")
	}
}
