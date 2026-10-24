package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type execTemplate struct {
	cmd *exec.Cmd
	r   io.Reader
	w   io.Writer
}

func (e *execTemplate) Run() error {
	return e.cmd.Run()
}

func (e *execTemplate) Start() error {
	return e.cmd.Start()
}

func (e *execTemplate) Wait() error {
	return e.cmd.Wait()
}

func newExecCmd(r io.Reader, w io.Writer, cmd string) execTemplate {
	templ := execTemplate{
		cmd: exec.Command(cmd),
		r:   r,
		w:   w,
	}

	templ.cmd.Stdin = r
	templ.cmd.Stdout = w
	return templ
}

func main() {

	in := strings.NewReader("1 2\n")
	out := strings.Builder{}

	cmd := newExecCmd(in, &out, "./sum")

	if err := cmd.Run(); err != nil {
		fmt.Println("Run Error" + err.Error())
		return
	}

	fmt.Println(out.String())

	// Again Run
	// Reset out and cmd
	out.Reset()
	in.Reset("6 7\n")
	cmd = newExecCmd(in, &out, "./sum")

	if err := cmd.Start(); err != nil {
		fmt.Println("Start Error: " + err.Error())
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait Error: " + err.Error())
		return
	}

	fmt.Println(out.String())
}
