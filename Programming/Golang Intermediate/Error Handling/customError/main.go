package main

import (
	"errors"
	"fmt"
)

// implementing the error interface. error interface has Error() method
// Error() method's return type is string

// Example - 01
type A struct {
	msg string
}

func (a *A) Error() string {
	return fmt.Sprintf(a.msg)
}

type AA struct {
	msg string
}

func (a AA) Error() string {
	return fmt.Sprintf(a.msg)
}

func getError() error {
	errA := &A{msg: "Error Occurs from A"} // without & -> error ? why ? because A receives a pointer receiver
	errAA := AA{msg: "Error Occurs from AA"}
	return errors.Join(errA, errAA)
}

// Example-02
type B struct {
	msg string
	val int64
}

func (b *B) Error() string {
	return fmt.Sprintf("%d cannot fit in %s data types", b.val, b.msg)
}

// Check if n canbe fit in int8
func checkBoundaryFor8ByteInt(n int64) error {
	if n > (1 << 8) {
		return &B{val: n, msg: "8byte"}
	}
	return nil
}

// Check if n canbe fit in int8
func checkBoundaryFor16ByteInt(n int64) error {
	if n > (1 << 16) {
		return &B{val: n, msg: "16byte"}
	}
	return nil
}

func main() {

	fmt.Println(getError())

	fmt.Println(checkBoundaryFor16ByteInt(1 << 17))
	fmt.Println(checkBoundaryFor8ByteInt(1 << 7))
	fmt.Println(checkBoundaryFor8ByteInt(1 << 9))

	fmt.Println()

	// Type assertion to check which error
	unknownErr2 := checkBoundaryFor8ByteInt(1 << 9)
	if _, ok := unknownErr2.(*B); ok {
		fmt.Println("B type error")
	}
}
