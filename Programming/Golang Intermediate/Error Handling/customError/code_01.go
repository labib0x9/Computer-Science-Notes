package main

import "fmt"

// Custom Error

type NotFoundError struct {
	Resource string
}

// Whay Error() ?
// See error interface
func (n *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", n.Resource)
}

type ValidateError struct {
	Field string
	Message string
}

func (v *ValidateError) Error() string {
	return fmt.Sprintf("%s field %s", v.Field, v.Message)
}

func findUser(id int) error {

	if id > 10 {
		return &ValidateError{Field: fmt.Sprintf("%d"), Message: "Valid"}
	}

	if id != 1 {
		return &NotFoundError{Resource: fmt.Sprint("User id: %d", id)}
	}
	return nil
}

func Main() {

	err := findUser(2)
	if err != nil {
		nfError, ok := err.(*NotFoundError)
		if ok {
			fmt.Println("Handle custom error: ", nfError)
		}

		vError, ok := err.(*ValidateError)
		if ok {
			fmt.Println("Handle custom error: ", vError)
		}
	}
}