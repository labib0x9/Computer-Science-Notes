package main

import (
	"container/list"
	"fmt"
)

func main() {

	l := list.New()             // creates a new list
	l.PushFront(1)              // insert at front
	Element := l.PushBack(2)               // insert at back
	l.InsertBefore(2, Element) // insert before a node
	l.InsertAfter(2, Element)  // insert after a node
	l.Remove(Element)          // remove a node

	l.Front() // returns first node
	l.Back()  // returns last node
	l.Len()   // returns length of list

	e := l.Front()
	_ = e.Value // access value
	e.Next()    // next node
	e.Prev()    // previous node

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
