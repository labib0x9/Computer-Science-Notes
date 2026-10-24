package main

import (
	"flag"
	"fmt"
)

func main() {

	syncFlag := flag.Bool("sync", false, "")
	userFlag := flag.String("user", "", "")

	flag.Parse()

	tags := flag.Args()

	if *syncFlag {
		fmt.Println("Not nil")
	} else {
		fmt.Println("nil")
	}

	fmt.Println("userFlag=", *userFlag)
	fmt.Println("tag=", tags)

	var str String
}