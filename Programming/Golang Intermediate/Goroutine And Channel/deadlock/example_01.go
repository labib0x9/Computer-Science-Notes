package main

import (
	"fmt"
	"sync"
	"time"
)

/*

labib0x9@TestGoCode $ ./out
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [sync.WaitGroup.Wait]:
sync.runtime_SemacquireWaitGroup(0xc000107e28?, 0x8?)
        /opt/homebrew/Cellar/go/1.27.2/libexec/src/runtime/sema.go:114 +0x38
sync.(*WaitGroup).Wait(0xc000012110)
        /opt/homebrew/Cellar/go/1.27.2/libexec/src/sync/waitgroup.go:206 +0xb8
main.main()
        /Users/labib0x9/Desktop/TestGoCode/main.go:37 +0x294

goroutine 4 [sync.Mutex.Lock]:
internal/sync.runtime_SemacquireMutex(0xc000068db8?, 0xd8?, 0xc000068de8?)
        /opt/homebrew/Cellar/go/1.27.2/libexec/src/runtime/sema.go:95 +0x28
internal/sync.(*Mutex).lockSlow(0xc000012130)
        /opt/homebrew/Cellar/go/1.27.2/libexec/src/internal/sync/mutex.go:149 +0x1e0
internal/sync.(*Mutex).Lock(0xc000012130)
        /opt/homebrew/Cellar/go/1.27.2/libexec/src/internal/sync/mutex.go:70 +0x4c
sync.(*Mutex).Lock(0x4a817c800?)
        /opt/homebrew/Cellar/go/1.27.2/libexec/src/sync/mutex.go:46 +0x20
main.main.func1({0x1029efd20, 0x1}, 0xc000012120, 0xc000012130)
        /Users/labib0x9/Desktop/TestGoCode/main.go:24 +0xd8
created by main.main in goroutine 1
        /Users/labib0x9/Desktop/TestGoCode/main.go:34 +0x1bc

goroutine 5 [sync.Mutex.Lock]:
internal/sync.runtime_SemacquireMutex(0xc000184db8?, 0xd8?, 0xc000184de8?)
        /opt/homebrew/Cellar/go/1.27.2/libexec/src/runtime/sema.go:95 +0x28
internal/sync.(*Mutex).lockSlow(0xc000012120)
        /opt/homebrew/Cellar/go/1.27.2/libexec/src/internal/sync/mutex.go:149 +0x1e0
internal/sync.(*Mutex).Lock(0xc000012120)
        /opt/homebrew/Cellar/go/1.27.2/libexec/src/internal/sync/mutex.go:70 +0x4c
sync.(*Mutex).Lock(0x4a817c800?)
        /opt/homebrew/Cellar/go/1.27.2/libexec/src/sync/mutex.go:46 +0x20
main.main.func1({0x1029f8fb8, 0x1}, 0xc000012130, 0xc000012120)
        /Users/labib0x9/Desktop/TestGoCode/main.go:24 +0xd8
created by main.main in goroutine 1
        /Users/labib0x9/Desktop/TestGoCode/main.go:35 +0x28c
labib0x9@TestGoCode $

*/

type value struct {
	mu    sync.Mutex
	value int
}

func main() {

	var wg sync.WaitGroup
	printSum := func(id string, v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()

		time.Sleep(20 * time.Second)

		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Println("ID=", id, " | a + b =", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)

	go printSum("a", &a, &b)
	go printSum("b", &b, &a)

	wg.Wait()
}
