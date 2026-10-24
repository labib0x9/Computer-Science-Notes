package main

import (
	"fmt"
	"sync"
)

// singleton design pattern

type DbInstance struct {
	flag string
}

func NewDbInstance() *DbInstance {
	return &DbInstance{
		flag: "true",
	}
}

var db *DbInstance
var once sync.Once

func GetDbInstance() *DbInstance {
	if db != nil {
		return db
	}
	once.Do(func ()  {
		db = NewDbInstance()	
	})
	return db
}

func main() {

	GetDbInstance()
	fmt.Printf("%p\n", db)

	GetDbInstance()
	fmt.Printf("%p\n", db)

	GetDbInstance()
	fmt.Printf("%p\n", db)

}
