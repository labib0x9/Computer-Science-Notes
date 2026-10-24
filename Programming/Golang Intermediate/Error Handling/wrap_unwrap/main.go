package main

import (
	"errors"
	"fmt"
)

// wrapping, Unwrapping

type NotFoundError struct {
	Resource string
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", n.Resource)
}

// type ValidateError struct {
// 	Field string
// 	Message string
// }

// func (v *ValidateError) Error() string {
// 	return fmt.Sprintf("%s field %s", v.Field, v.Message)
// }

func findUser(id int) error {

	// if id > 10 {
	// 	return &ValidateError{Field: fmt.Sprintf("%d"), Message: "Valid"}
	// }

	// if id != 1 {
	return &NotFoundError{Resource: fmt.Sprint("User id: %d", id)}
	// }
	// return nil
}

// wrapping
func getUser(id int) error {
	err := findUser(id)
	if err != nil {
		return fmt.Errorf("find user failed: %w", err)
	}
	return nil
}

func main() {

	err := findUser(12)
	if err != nil {
		var nfError *NotFoundError
		ok := errors.As(err, &nfError) // unwrapped
		if ok {
			fmt.Println("Handle custom err: ", nfError)
		} else {
			fmt.Println("Other erroe")
		}
	}

}
