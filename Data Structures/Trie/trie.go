// https://codeforces.com/contest/1902/problem/E

package main

import (
	"bufio"
	"fmt"
	"os"
	_ "sort"
	"strconv"
	"strings"
)

var (
	reader = bufio.NewReaderSize(os.Stdin, 1e6)
	writer = bufio.NewWriterSize(os.Stdout, 1e6)
)

func ReadInt() int {
	s, _ := reader.ReadString('\n')
	x, _ := strconv.Atoi(strings.TrimSpace(s))
	// x, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return x
}

func ReadPairs() (int64, int64, int64, int64, int64) {
	line, _ := reader.ReadString('\n')
	parts := strings.Fields(line)
	u, _ := strconv.ParseInt(parts[0], 10, 64)
	v, _ := strconv.ParseInt(parts[1], 10, 64)
	x, _ := strconv.ParseInt(parts[2], 10, 64)
	y, _ := strconv.ParseInt(parts[3], 10, 64)
	z, _ := strconv.ParseInt(parts[4], 10, 64)

	return u, v, x, y, z
}

func ReadArray() []string {
	s, _ := reader.ReadString('\n')
	fields := strings.Fields(s)
	arr := make([]string, len(fields))
	for i, _ := range fields {
		// arr[i], _ = strconv.ParseInt(v, 10, 64)
		// // arr[i], _ = strconv.ParseFloat(v, 64)
		arr[i] = strings.TrimSpace(s)
	}
	return arr
}

func ReadString() string {
	s, _ := reader.ReadString('\n')
	return strings.TrimSpace(s)
}

func Write(a ...interface{}) {
	fmt.Fprintln(writer, a...)
}

type node struct {
	child [26]*node
	isEnd bool
	freq  int
}

func NewNode() *node {
	return &node{
		isEnd: false,
	}
}

type Trie struct {
	root *node
}

func (t *Trie) Insert(s string) {
	cur_node := t.root
	cur_node.freq++
	for _, c := range s {
		ch := c - 'a'
		if cur_node.child[ch] == nil {
			cur_node.child[ch] = NewNode()
		}
		cur_node = cur_node.child[ch]
		cur_node.freq++
	}
	cur_node.isEnd = true
}

func (t *Trie) LCP(s string) (cnt int) {
	cur_node := t.root
	for _, c := range s {
		ch := c - 'a'
		if cur_node.child[ch] == nil {
			return
		}
		cur_node = cur_node.child[ch]
		cnt += cur_node.freq
	}
	return
}

func dfs(cur *node, s string) {
	if cur == nil {
		return
	}
	for i := 0; i < 26; i++ {
		if cur.child[i] == nil {
			continue
		}
		dfs(cur.child[i], s+string(byte('a'+i)))
	}
}

func NewTrie() Trie {
	return Trie{
		root: NewNode(),
	}
}

func solve() {

	n := ReadInt()
	arr := make([]string, n, n)

	for i := 0; i < n; i++ {
		arr[i] = ReadString()
	}

	t := NewTrie()
	ans := int64(0)

	for i := range arr {
		ans += int64(len(arr[i]))
		t.Insert(arr[i])
	}

	ans = ans * int64(2*n)

	for i := range arr {
		sLen := len(arr[i])
		newS := make([]byte, sLen, sLen)
		for j := sLen - 1; j >= 0; j-- {
			newS[sLen-j-1] = arr[i][j]
		}

		ans = ans - int64(2*t.LCP(string(newS)))
	}

	Write(ans)
}

func main() {
	defer writer.Flush()

	testCase := !true
	if testCase {
		t := ReadInt()
		for i := 0; i < t; i++ {
			solve()
		}
	} else {
		solve()
	}
}
