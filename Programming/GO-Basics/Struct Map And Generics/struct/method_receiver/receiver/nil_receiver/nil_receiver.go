package main

type IntTree struct {
	val int
	left, right *IntTree
}

func (it *IntTree) Insert(val int) *IntTree {
	if it == nil {
		return &IntTree{val: val}
	}
	if it.val < val {
		it.right = it.right.Insert(val)
	} else {
		it.left = it.left.Insert(val)
	}
	return it
}

func (it *IntTree) Search(val int) bool {
	switch {
	case it == nil:
		return false
	case it.val > val:
		return it.left.Search(val)
	case it.val < val:
		return it.right.Search(val)
	default:
		return true
	}
}

func main() {

}