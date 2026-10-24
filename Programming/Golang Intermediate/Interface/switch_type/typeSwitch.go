package main

import "fmt"

// type switch
func detectDataType(i interface{}) {
	// With variable
	switch v := i.(type) {
	case string:
		fmt.Println(v, ": String")
	case int:
		fmt.Println(v, ": Int")
	case float64:
		fmt.Println(v, ": Float")
	}

	// Without variable
	switch i.(type) {
	case string:
		fmt.Println("Without variable string")
	}
}

type A struct {
	a string
}

type B struct {
	b string
}

type C map[string]int

func detectDataTypeInterface(item ...interface{}){
	for _, i := range item {
		switch v := i.(type) {
		case A:
			fmt.Println("A struct:", v.a)
		case B:
			fmt.Println("B struct:", v.b)
		case C:
			fmt.Println("C type")
		}
	}
}

func DataType() {

	a := "ABC"
	b := 12
	c := 12.34

	detectDataType(a)
	detectDataType(b)
	detectDataType(c)

	mp := make(C)

	// slice of interface
	data := []interface{} {
		A{a: "ABC"},
		B{b: "CBA"},
		mp,
		make(C),
	}
	detectDataTypeInterface(data...)
}