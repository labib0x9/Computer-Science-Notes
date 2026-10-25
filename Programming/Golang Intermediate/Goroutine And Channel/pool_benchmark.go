package main

// main_test.go

/*
labib0x9@TestGoCode $ go test -benchtime=10s -bench=BenchmarkNetworkRequest
goos: darwin
goarch: arm64
pkg: test
cpu: Apple M1
BenchmarkNetworkRequest-8             10        1002302904 ns/op		// -> 1e9
PASS
ok      test    11.209s

*/

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

// Without pooling

func connectToServer() any { // must be thread-safe
	time.Sleep(time.Second)
	return struct{}{}
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			panic(err)
		}
		defer server.Close()

		wg.Done() // wait group exits ?

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("accept error: %v", err)
				continue
			}
			connectToServer()
			fmt.Fprint(conn, "")
			conn.Close()
		}
	}()
	return &wg
}

// func init() {
// 	daemonStarted := startNetworkDaemon()
// 	daemonStarted.Wait()
// }

func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatalf("Dial error= %v\n", err)
		}

		if _, err := io.ReadAll(conn); err != nil {
			b.Fatalf("Read error= %v\n", err)
		}

		conn.Close()
	}
}

// With Pooling
/*

labib0x9@TestGoCode $ go test -benchtime=10s -bench=BenchmarkNetworkRequest
goos: darwin
goarch: arm64
pkg: test
cpu: Apple M1
BenchmarkNetworkRequest-8           3837           4232341 ns/op
PASS
ok      test    44.628s

*/

func warmServiceConnCache() *sync.Pool {
	// p has 10 pool.
	p := &sync.Pool{
		New: connectToServer, // connection
	}
	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}

func startNetworkDaemonPool() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		connPool := warmServiceConnCache()
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			panic(err)
		}
		defer server.Close()

		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("accept error: %v", err)
				continue
			}
			svcConn := connPool.Get()
			fmt.Fprint(conn, "")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()
	return &wg
}

func init() {
	daemonStarted := startNetworkDaemonPool()
	daemonStarted.Wait()
}
