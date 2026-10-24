package main

import (
	"encoding/json"
	"fmt"
)

type Issue struct {
	TotalCount int        `json:"total_count"`
	HTMLURL    string     `json:"html_url"`
	User       *IssueUser `json:"user"`
}

type IssueUser struct {
	Name string `json:"name"`
	Id   int
}

func Unmarshal() {

	// Design the struct using this data
	data := []byte(`{
		"total_count" : 12,
		"html_url" : "https://api.issue.com/get/labib.faisal.123",
		"user" : {
			"name" : "Labib"
		}
	}`)

	var IssueLabib Issue

	// Json to struct
	if err := json.Unmarshal(data, &IssueLabib); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(IssueLabib)
	fmt.Println(IssueLabib.User.Name)
}
