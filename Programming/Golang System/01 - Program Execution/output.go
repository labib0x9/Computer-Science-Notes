package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("./sum") // sum is a ext file in the current dir with no arg

	// io.Reader implements
	in := strings.NewReader("1 2\n")

	// // I/O mechanism
	cmd.Stdin = in

	// runs the cmd and stores the output in out
	// out is []byte
	// io.Writer = bytes.Buffer
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Output Error" + err.Error())
		return
	}

	fmt.Println(string(out))
}
