package main

import "encoding/hex"

func main() {

	hex.NewDecoder()
	hex.NewEncoder()

	hex.EncodeToString()
	hex.DecodeString()

	hex.Dump()
	hex.Dumper()

	hex.Encode()
	hex.Decode()
}

// Use Cases:

//     Display hash values (e.g., SHA256, MD5).

//     Debugging byte sequences

//     Network packet dumps

//     Working with MAC/IP addresses