package main

/*

labib0x9@TestGoCode $ go run main.go
# command-line-arguments
./main.go:25:4: invalid operation: cannot receive from send-only channel chan<- any inStream (variable of type chan<- any)
./main.go:26:2: invalid operation: cannot send to receive-only channel <-chan any outStream (variable of type <-chan any)

*/

func main() {
	inStream := make(chan<- any)
	outStream := make(<-chan any)

	<-inStream
	outStream <- struct{}{}
}
