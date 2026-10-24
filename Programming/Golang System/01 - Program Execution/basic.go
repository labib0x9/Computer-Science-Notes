package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("./sum") // sum is a ext file in the current dir with no arg

	// io.Write implements
	// stores output in buf
	out := strings.Builder{}

	// io.Reader implements
	in := strings.NewReader("1 2\n")

	// I/O mechanism
	cmd.Stdout = &out
	cmd.Stdin = in

	// runs a program and wait to finish
	// Run() = Start() + Wait()
	if err := cmd.Run(); err != nil {
		fmt.Println("Run Error" + err.Error())
		return
	}

	fmt.Println(out.String())

	//		Same - but error because cmd can be run only once.
	// For multiple runnig copy the struct into another	variable
	//
	out.Reset()
	if err := cmd.Start(); err != nil {
		fmt.Println("Start Error: " + err.Error())
		return
	}

	// process is running, do whatever you want....

	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait Error: " + err.Error())
		return
	}

	fmt.Println(out.String())
}
