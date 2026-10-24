package main

import (
	"container/ring"
	"fmt"
)

// Todo -> Practical use of ring in system

func main() {

	r := ring.New(3) // create new ring of size n
	r.Next()         // move forward
	r.Prev()         // move backward
	_ = r.Value      // get/set value

	r2 := ring.New(3)

	// Initialize values
	for i := 0; i < 3; i++ {
		r.Value = i
		r = r.Next()

		r2.Value = i
		r2 = r2.Next()
	}

	// iterate and apply func
	r.Do(func(v interface{}) {
		fmt.Println(v)
	})

	r.Link(r2) // link r2 after r

	r2 = r.Unlink(2) // unlink n elements after r

}
