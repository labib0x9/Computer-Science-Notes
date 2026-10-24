package main

import (
	"encoding/json"
	"fmt"
)

type Movie struct {
	Title string
	Year int `json:"released"`
}

func Marshal() {

	movieList := []Movie{
		{Title: "Lord of the rings", Year: 2011},
		{Title: "The hobbits", Year: 2013},
	}

	// struct to json
	// []byte, error
	data, err := json.Marshal(movieList)
	if err != nil {
		fmt.Println(err)
		return
	}

	// What is a byte array ??
	fmt.Println(data)


	// Json Output
	fmt.Println(string(data))

	data, err = json.MarshalIndent(movieList, "", " ")
	if err != nil {
		fmt.Println(err)
		return;
	}

	fmt.Println(string(data))
}