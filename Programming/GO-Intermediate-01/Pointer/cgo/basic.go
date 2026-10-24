package main

/*
#include<stdio.h>
#include<stdlib.h>

void Print(int x) {
	printf("%d\n", x);
}

int Add(int a, int b) {
	return a + b;
}

int Mult(int a, int b) {
	return a * b;
}
*/
import "C" // No newline between C code and import

import "fmt"

func main() {

	C.Print(2)
	s := int(C.Add(2, 3))

	fmt.Println(s)
	C.Print(C.int(s))

	C.Print(C.Mult(2, 3))
}
