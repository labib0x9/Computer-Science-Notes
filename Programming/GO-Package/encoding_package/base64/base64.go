package main

// incomplete

import (
	"encoding/base64"
	"io/ioutil"
)

func main() {
	p := []byte("Hello")

	encoded := base64.StdEncoding.EncodeToString(p)
	decoded, _ := base64.StdEncoding.DecodeString(encoded)

	decoded = decoded

	base64.NewDecoder()
	base64.NewEncoder()

	file, err := ioutil.ReadFile("a.jpeg")
	fileBase64 := base64.StdEncoding.EncodeToString(file)

	base64.URLEncoding.EncodeToString()
	base64.URLEncoding.DecodeString()
}

// Useful for:
//		URLs
//		JWTs
//		Query parameters
