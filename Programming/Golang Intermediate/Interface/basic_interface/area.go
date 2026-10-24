package main

import "fmt"

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return 3.1416 * c.radius
}

type Rectangular struct {
	width, height float64
}

func (r Rectangular) Area() float64 {
	return r.height * r.width
}

type Shape interface {
	Area() float64
}

func calculateShape(s Shape) float64 {
	return s.Area()
}

func main() {

	rec := Rectangular{ width: 10, height: 30}
	cir := Circle{radius: 10}

	fmt.Println(calculateShape(rec))
	fmt.Println(calculateShape(cir))
}