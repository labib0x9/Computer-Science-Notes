package main

import "fmt"

func main() {

	/**
	for short declaration := ,
	inner scope declares a new variables,
	and this variables lives only inside that inner scope,
	inner scope does't affect the outer scope variable
	**/

	x := 20

	if true {
		x := 30
		fmt.Println(x) // 30
	}

	fmt.Println(x) // 20

}
