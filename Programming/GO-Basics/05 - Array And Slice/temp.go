// // // PrintSlice(s1)
// // // 	PrintSlice(s2)

// // // 	PrintSlice(s3)
// // // 	PrintSlice(s4)

// // // 	// s[i : j] -> i to j - 1
// // // 	// 0 <= i <= j <= cap(s)
// // // 	// slice beyond cap(s) causes panic
// // // 	// slice beyond len(s) extends the slice
// // // 	// length = j - i
// // // 	arr := []int{1, 2, 4, 5, 6, 7, 8, 9, 10}
// // // 	s5 := arr[2:4]
// // // 	PrintSlice(s5)

// // // 	// cap(s) = 8
// // // 	// len(s) = 5
// // // 	s := []int{0, 1, 2, 3}
// // // 	s = append(s, 4)

// // // 	// Within cap, beyond len
// // // 	s6 := s[2:6]
// // // 	PrintSlice(s6)

// // // 	// Beyond cap, panic
// // // 	// s7 := s[2 : 10]
// // // 	// PrintSlice(s7) // Panic

// // // 	// len(s) = 2
// // // 	// cap(s) = 4
// // // 	// len(s8) = 1
// // // 	// cap(s8) = 3 = cap(s) - i
// // // 	// s8 = s[i : j]
// // // 	s = make([]int, 2, 4)
// // // 	s8 := s[1:]
// // // 	fmt.Println(len(s8), cap(s8))

// // // 	// len(s) = 5
// // // 	// cap(s) = 10
// // // 	// len(s8) = 3
// // // 	// cap(s8) = 8 = cap(s) - i
// // // 	// s9 = s[i : j]
// // // 	s = make([]int, 5, 10)
// // // 	s9 := s[2:5]
// // // 	fmt.Println(len(s9), cap(s9))

// // // 	// clear
// // // 	clear(s)

// // // 	//
// // // 	// s[0:len(s1)]

// // // 	/** --- **/

// // // 	// Interesting case
// // // 	newS := []int{1, 2, 3}
// // // 	newY := newS // newS and newY point to the same backing array
// // // 	newY[0] = 100

// // // 	// Why ?
// // // 	// Because they both share the same backed-array,
// // // 	// so, changing one can effect other,
// // // 	// use copy() function to avoid this behaviour

// // // 	PrintSlice(newS) // Output: [100 2 3]
// // // 	PrintSlice(newY) // Output: [100 2 3]

// // // 	newY = append(newY, 4) // Now newY points to a new backing array
// // // 	newY[0] = 200          // Modifying newY does not affect newS anymore
// // // 	PrintSlice(newS)       // Output: [100 2 3]
// // // 	PrintSlice(newY)       // Output: [200 2 3 4]


// // /****/


// // package main

// // import "fmt"

// // func PrintSlice(s []int) {
// // 	fmt.Println(s)
// // }

// // func main() {
// // 	s1 := []int{10, 20, 10}
// // 	s2 := make([]int, 0, 10)

// // 	// Didn't append
// // 	Append(s1)
// // 	Append(s2)

// // 	PrintSlice(s1)
// // 	PrintSlice(s2)

// // 	// Append
// // 	AppendPointer(&s1)
// // 	AppendPointer(&s2)

// // 	PrintSlice(s1)
// // 	PrintSlice(s2)

// // 	// printPointerSlice(&s1)
// // }

// // // to-do Print address
// // func Append(s []int) {
// // 	s = append(s, 40)
// // }

// // // Passing pointer
// // func AppendPointer(s *[]int) {
// // 	*s = append(*s, 30)
// // }

// // // What is happening here ??
// // func printPointerSlice(s *[]int) {

// // 	// a is a copy of orignal s[i]
// // 	for _, a := range *s {
// // 		fmt.Printf("%d ", a)
// // 	}
// // 	fmt.Println()

// // 	// original value of s[i]
// // 	// dereferencing using (*s)
// // 	for i := range *s {
// // 		fmt.Printf("%d ", (*s)[i])
// // 	}
// // 	fmt.Println()

// // 	// append another slice
	
// // 	// pointer slice
// // 	// var s []*int
// // }


// package main

// import "fmt"

// // len(s) == cap(s),
// // In this situation append in s
// // increases the capacity by 2 * len(s)
// func main() {
// 	s := []int{1}

// 	fmt.Println(len(s), cap(s))
// 	s = append(s, 3)
// 	fmt.Println(len(s), cap(s))
// 	s = append(s, 4)
// 	fmt.Println(len(s), cap(s))
// 	s = append(s, 5)
// 	fmt.Println(len(s), cap(s))
// 	s = append(s, 6)
// 	fmt.Println(len(s), cap(s))
// }

// // s[:] referencing the slice of s which is extendable to s's capacity
// // so changing on s[:] in range 0 to len(s) - 1, affects s (change the value)
// // after append if slices length becomes  more than the capacity of s
// // then s doesn't affected
// // len(s) = 2
// // cap(s) = 4
// func Slice() {
// 	s := make([]int, 2, 4)
// 	Print2(s[:])
// 	fmt.Println("main", s)
// }

// // Affects
// // because the memory of s[:] = c is the actual s
// // s[:] = c length doesn't exceed capacity of s
// // changes s[0] = 300
// // len(c) = 3
// // cap(s) = 4
// func Print(c []int) {
// 	fmt.Println(c, len(c), cap(c))
// 	fmt.Println(c)
// 	c = append(c, 10)
// 	fmt.Println(c)

// 	c[0] = 3000
// 	fmt.Println(c)
// }

// // Affects
// // because the memory of s[:] = c is the actual s
// // s[:] = c length doesn't exceed capacity of s
// // changes s[0] = 500
// // len(c) = 4
// // cap(s) = 4
// func Print1(c []int) {
// 	c[0] = 3000
// 	c = append(c, 20)
// 	c = append(c, 30)
// 	c[0] = 5000
// 	fmt.Println(c)
// }

// // Doesn't Affects
// // because the memory of s[:] = c is now changes as
// // s[:] = c length exceed capacity of s
// // len(c) = 6
// // cap(s) = 4
// func Print2(c []int) {
// 	c = append(c, 20)
// 	c = append(c, 30)
// 	c = append(c, 40)
// 	c = append(c, 50)
// 	c[0] = 5000
// 	fmt.Println(c)
// }
