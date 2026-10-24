package main

import (
	"sync"
	"testing"
)

/*
labib0x9@TestGoCode $ go test -bench=BenchmarkContextSwitch -cpu=1 main_test.go
goos: darwin
goarch: arm64
cpu: Apple M1
BenchmarkContextSwitch  10384272               107.4 ns/op
PASS
ok      command-line-arguments  1.708s
*/

// Why one CPU is used here ??

func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin // blocks until the channel is closed
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}

	receiver := func() {
		defer wg.Done()
		<-begin // blocks until the channel is closed
		for i := 0; i < b.N; i++ {
			<-c
		}
	}

	wg.Add(2)
	go sender()
	go receiver()

	b.StartTimer()
	close(begin)
	wg.Wait()
}
