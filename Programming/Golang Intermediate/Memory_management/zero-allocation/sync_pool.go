package main

import (
	"fmt"
	"sync"
)

type Buffer []byte

func (b *Buffer) Append(s string) {
	*b = append(*b, s...)
}

func (b *Buffer) Truncate() {
	*b = (*b)[0:]
}

var bufPool = sync.Pool{
	New: func() any {
		fmt.Println("Allocating ...")
		return new(Buffer)
	},
}

func main() {

	// buf := bufPool.Get() -> Append() method doesn't work, without type assertion buf's type is any, we need to specify the type is Buffer.

	buf := bufPool.Get().(*Buffer) // pool is empty, calls New

	buf.Append("ABCD")
	fmt.Println(buf)

	buf.Truncate()
	bufPool.Put(buf) // Put buf back to pool

	x := bufPool.Get().(*Buffer) // retrive from pool, pool is now empty
	fmt.Println(x)

	y := bufPool.Get().(*Buffer) // pool is empty, calls New
	fmt.Println(y)
}
