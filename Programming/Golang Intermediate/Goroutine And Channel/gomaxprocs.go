package main

import (
	"fmt"
	"runtime"
)

/*

labib0x9@TestGoCode $ go run main.go
8
8
4
16

*/

func main() {

	// my pc core
	fmt.Println(runtime.NumCPU())

	// default = numCPU()
	// 8 threads competing for 8 cores,
	fmt.Println(runtime.GOMAXPROCS(0))

	// 4 threads running on 8 cores, 4 idle + 4 running
	runtime.GOMAXPROCS(4)
	fmt.Println(runtime.GOMAXPROCS(0))

	// 16 threads competing for 8 cores,
	runtime.GOMAXPROCS(16)
	fmt.Println(runtime.GOMAXPROCS(0))

}
