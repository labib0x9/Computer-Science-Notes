package main

import (
	"fmt"
	"sync"
)

/*
labib0x9@TestGoCode $ go build -gcflags="all=-N -l" -o out main.go
labib0x9@TestGoCode $ ./out
K= bar V= 2
K= foo V= 1
labib0x9@TestGoCode $ ./out
K= bar V= 2
labib0x9@TestGoCode $ ./out
K= bar V= 2
labib0x9@TestGoCode $ ./out
K= bar V= 2
labib0x9@TestGoCode $ ./out
K= bar V= 2
K= foo V= 1
labib0x9@TestGoCode $ ./out
K= bar V= 2
K= foo V= 1
labib0x9@TestGoCode $ ./out
K= bar V= 2
*/

/*
labib0x9@TestGoCode $ go build -race -o out main.go
labib0x9@TestGoCode $ ./out
K= bar V= 2
==================
WARNING: DATA RACE
Write at 0x00c0001800f0 by goroutine 7:
  runtime.mapaccess2_faststr()
      /opt/homebrew/Cellar/go/1.27.2/libexec/src/internal/runtime/maps/runtime_faststr.go:120 +0x2bc
  main.main.func1()
      /Users/labib0x9/Desktop/TestGoCode/main.go:30 +0x40

Previous write at 0x00c0001800f0 by main goroutine:
  runtime.mapaccess2_faststr()
      /opt/homebrew/Cellar/go/1.27.2/libexec/src/internal/runtime/maps/runtime_faststr.go:120 +0x2bc
  main.main()
      /Users/labib0x9/Desktop/TestGoCode/main.go:33 +0xa8

Goroutine 7 (running) created at:
  main.main()
      /Users/labib0x9/Desktop/TestGoCode/main.go:29 +0x8c
==================
Found 1 data race(s)
*/

/*
use mutex lock and unlock, DATA Race will be gone, but race condition will still present.

use wait group and mutex, then DATA Race and Race Condition are both gone...
*/

// BOTH DATA RACE AND RACE CONDITION
/*
func main() {
	mp := make(map[string]int)

	go func() {
		mp["foo"] = 1
	}()

	mp["bar"] = 2

	for k, v := range mp {
		fmt.Println("K=", k, "V=", v)
	}
}
*/

/* NO DATA RACE, BUT RACE CONDITION
labib0x9@TestGoCode $ go build -race -o out main.go
labib0x9@TestGoCode $ ./out
It should be 10000, GOT= 9720
*/

/*
func main() {

	var mu sync.Mutex
	var value int

	for range 10000 {
		go func() {
			mu.Lock()
			value++
			mu.Unlock()
		}()
	}

	mu.Lock()
	fmt.Println("It should be 10000, GOT=", value)
	mu.Unlock()
}
*/

/*
labib0x9@TestGoCode $ go build -race -o out main.go
labib0x9@TestGoCode $ ./out
==================
WARNING: DATA RACE
Read at 0x00c00009e048 by goroutine 8:
  main.main.func1()
      /Users/labib0x9/Desktop/TestGoCode/main.go:110 +0x68

Previous write at 0x00c00009e048 by goroutine 7:
  main.main.func1()
      /Users/labib0x9/Desktop/TestGoCode/main.go:110 +0x78

Goroutine 8 (running) created at:
  main.main()
      /Users/labib0x9/Desktop/TestGoCode/main.go:108 +0x6c

Goroutine 7 (finished) created at:
  main.main()
      /Users/labib0x9/Desktop/TestGoCode/main.go:108 +0x6c
==================
==================
WARNING: DATA RACE
Write at 0x00c00009e048 by goroutine 9:
  main.main.func1()
      /Users/labib0x9/Desktop/TestGoCode/main.go:110 +0x78

Previous write at 0x00c00009e048 by goroutine 21:
  main.main.func1()
      /Users/labib0x9/Desktop/TestGoCode/main.go:110 +0x78

Goroutine 9 (running) created at:
  main.main()
      /Users/labib0x9/Desktop/TestGoCode/main.go:108 +0x6c

Goroutine 21 (running) created at:
  main.main()
      /Users/labib0x9/Desktop/TestGoCode/main.go:108 +0x6c
==================
It should be 10000, GOT= 8300
Found 2 data race(s)
*/

// RACE Condition and Data Race is present. Even though every 10000 goroutine completes.
/*
func main() {

	var wg sync.WaitGroup
	var value int

	for range 10000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value++
		}()
	}

	wg.Wait()
	fmt.Println("It should be 10000, GOT=", value)
}
*/

// No Race Condition and data race
func main() {

	var mu sync.Mutex
	var wg sync.WaitGroup
	var value int

	for range 10000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			value++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("It should be 10000, GOT=", value)
}
