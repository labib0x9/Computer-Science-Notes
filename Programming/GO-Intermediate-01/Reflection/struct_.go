package main

import (
	"fmt"
	"reflect"
)

type MyInt int

type User struct {
	Name string
	Age  MyInt `info:"age"`
}

func main() {

	u := User{Name: "Faisal", Age: 12}

	ut := reflect.TypeOf(&u)

	fmt.Println(ut)        // *main.User
	fmt.Println(ut.Name()) // ""
	fmt.Println(ut.Kind()) // ptr

	// fmt.Println(ut.NumField()) // panic: reflect: NumField of non-struct type *main.User
	fmt.Println("----")

	ut = ut.Elem()         // dereference the pointer
	fmt.Println(ut)        // main.User
	fmt.Println(ut.Name()) // User
	fmt.Println(ut.Kind()) // struct
	fmt.Println("----")

	fmt.Println(ut.NumField())           // 2
	fmt.Println(ut.Field(0).Name)        // Name
	fmt.Println(ut.Field(0).Type.Name()) // string
	fmt.Println(ut.Field(1).Name)        // Age
	fmt.Println(ut.Field(1).Type)        // MyInt
	fmt.Println(ut.Field(1).Type.Kind()) // int

	fmt.Println(ut.Field(1).Tag)             // info:"age"
	fmt.Println(ut.Field(1).Tag.Get("info")) // age
	fmt.Println("----")

	// accessing field by name
	f, ok := ut.FieldByName("Age")
	if ok {
		fmt.Println(f.Name) // Age
	}

	fmt.Println("----")

	fieldNum := ut.NumField()
	for i := 0; i < fieldNum; i++ {
		f := ut.Field(i)
		fmt.Println(f.Type)
		fmt.Println(f.Name)
		fmt.Println(f.Type.Kind())
		fmt.Println("---")
	}

	fmt.Println("----")

	uv := reflect.ValueOf(&u)
	uv = uv.Elem() // dereference the pointer

	fmt.Println(uv.CanSet()) // true
	fmt.Println("---")

	nameField := uv.Field(0)
	fmt.Println(nameField.Type())   // string
	fmt.Println(nameField)          // Faisal
	fmt.Println(nameField.CanSet()) // true

	nameField.SetString("Labib Al Faisal")
	fmt.Println(u.Name) // Labib Al Faisal

	fmt.Println("----")

	ageField := uv.FieldByName(ut.Field(1).Name)
	fmt.Println(ageField.Type())   // main.MyInt
	fmt.Println(ageField)          // 12
	fmt.Println(ageField.CanSet()) // true

	ageField.SetInt(25)
	fmt.Println(u.Age) // 25

	fmt.Println("----")
	var output string = "{\n"
	fieldNum = ut.NumField()
	for i := 0; i < fieldNum; i++ {
		f := ut.Field(i)
		v := uv.Field(i)
		output += fmt.Sprintf("\t%s: %v", f.Name, v)
		if i != fieldNum-1 {
			output += ",\n"
		}
	}

	output += "\n}"
	fmt.Println(output) // { Name: Labib Al Faisal, Age: 25 }
}
