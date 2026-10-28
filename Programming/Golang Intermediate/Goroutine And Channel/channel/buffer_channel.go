package main

/*

Why is output like this? did not expect it

labib0x9@TestGoCode $ go run main.go
Seeding: 0
Seeding: 1
Seeding: 2
Seeding: 3
Seeding: 4
Seeding: 5

Recieved 0
Recieved 1
Recieved 2
Recieved 3
Recieved 4
Recieved 5

Seeding: 6
Seeding: 7
Seeding: 8
Seeding: 9

Producer done!

Recieved 6
Recieved 7
Recieved 8
Recieved 9
*/

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4) // buffer channel
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer done!")

		for i := 0; i < 10; i++ {
			fmt.Fprintf(&stdoutBuff, "Seeding: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "Recieved %v\n", integer)
	}
}
