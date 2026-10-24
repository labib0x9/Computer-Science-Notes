package main

import "fmt"

func Print[T any] (arr []T) {
	for _, value := range arr {
		fmt.Println(value)
	}
}
