package main

import (
	"fmt"
	"runtime"
	"runtime/trace"
)

// Incomplete

func main() {

	trace.WithRegion()
	trace.IsEnabled()
	trace.NewTask()
	trace.Start()
	trace.StartRegion()
	trace.Stop()

	fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(0))
}
