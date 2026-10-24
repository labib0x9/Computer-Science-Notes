package main

import "fmt"

type MyInt int

func main() {

	var r any
	var m MyInt = 10

	r = m

	i, ok := r.(MyInt)
	fmt.Println(i + 10)

	_, ok = r.(string)
	_, ok = r.(int)

	if ok {
		
	}

}