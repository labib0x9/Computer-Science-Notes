package main

import (
	"errors"
	"fmt"
	"strings"
)

type MyNewInterface interface {
	Task(name string) (out string, err error)
}

// another type
type FancyName string

func (f FancyName) Task(name string) (out string, err error) {
	return "Fancy Name: " + string(f) + " " + name, nil
}

// this struct implements the interface
type Struct1 struct {
	mni MyNewInterface
}

func (s Struct1) Task(name string) (out string, err error) {
	if len(name) == 0 {
		return "", errors.New("name cannot be empty")
	}
	if strings.HasSuffix(name, "Faisal") {
		fmt.Println("Err = Contains Faisal as last name")
		return "", errors.New("name is not allowd")
	}
	name += " :::"

	return s.mni.Task(name)
}

func Process(i MyNewInterface, s string) {
	out, err := i.Task(s)
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println("[Process] ", out)
}

func main() {

	fname := FancyName("Labib")
	Process(fname, "Al Faisal")

	sfname := Struct1{
		mni: fname,
	}

	Process(sfname, "Al Faisal")
}
