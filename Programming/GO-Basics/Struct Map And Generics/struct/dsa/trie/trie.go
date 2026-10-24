package main

type Node struct {
	children map[rune]*Node
	isEnd    bool
}

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{
		root: &Node{children: make(map[rune]*Node)},
	}
}

func (t *Trie) Insert(word string) {
	curr := t.root
	for _, ch := range word {
		if _, exists := curr.children[ch]; !exists {
			curr.children[ch] = &Node{children: make(map[rune]*Node)}
		}
		curr = curr.children[ch]
	}
	curr.isEnd = true
}

func (t *Trie) Search(word string) bool {
	curr := t.root
	for _, ch := range word {
		if _, exists := curr.children[ch]; !exists {
			return false
		}
		curr = curr.children[ch]
	}
	return curr.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	curr := t.root
	for _, ch := range prefix {
		if _, exists := curr.children[ch]; !exists {
			return false
		}
		curr = curr.children[ch]
	}
	return true
}

func main() {

}
