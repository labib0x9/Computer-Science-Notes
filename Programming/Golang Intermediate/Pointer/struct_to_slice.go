package main

import (
	"fmt"
	"unsafe"
)

type User struct {
	name string
	age  int
}

func GetUser() User {
	return User{
		name: "Labib Al Faisal",
		age:  25,
	}
}

const (
	nameSize = unsafe.Sizeof(User{}.name)
	ageSize  = unsafe.Sizeof(User{}.age)

	nameOffset = unsafe.Offsetof(User{}.name)
	ageOffset  = unsafe.Offsetof(User{}.age)

	userSize = unsafe.Sizeof(User{})
)

func main() {

	fmt.Println(nameSize, nameOffset)
	fmt.Println(ageSize, ageOffset)

	user_01 := GetUser()

	// access fields
	ptr := unsafe.Pointer(&user_01)
	name := *(*string)(unsafe.Pointer(uintptr(ptr) + nameOffset))
	age := *(*int)(unsafe.Pointer(uintptr(ptr) + ageOffset))

	fmt.Println(name)
	fmt.Println(age)

	age = 36

	fmt.Println(user_01)

	// convert to []byte
	ptrByte := (*byte)(ptr)
	userSlice := unsafe.Slice(ptrByte, userSize)

	nameSlice := userSlice[nameOffset : nameOffset+nameSize]
	ageSlice := userSlice[ageOffset : ageOffset+ageSize]

	fmt.Println(userSlice)
	fmt.Println(nameSlice)
	fmt.Println(ageSlice)

	ageSlice[ageSize-3] = 1 // changes the struct value
	fmt.Println(ageSlice)
	fmt.Println(user_01)

	// if you change, killed signal is shown..
	// you just changed the stringHeader -> idk address
	// nameSlice[nameSize - 3] = 11
	// fmt.Println(user_01)

	// convert a single field to []byte
	newAgeSlice := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr)+ageOffset)), ageSize)
	fmt.Println(newAgeSlice)

	myName := "labib0x0hunter@9xzer0"
	fmt.Println(unsafe.Sizeof(myName))
	newPtr := (*byte)(unsafe.Pointer(&myName))
	newS := unsafe.Slice(newPtr, nameSize)

	fmt.Println(user_01)
	fmt.Println(nameSlice)
	copy(nameSlice, newS)
	fmt.Println(nameSlice)
	fmt.Println(user_01)

}
