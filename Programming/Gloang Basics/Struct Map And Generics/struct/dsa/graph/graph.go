package main

import "fmt"

type Graph struct {
	child map[int][]int
}

func NewGraph() *Graph {
	return &Graph{
		child: make(map[int][]int),
	}
}

func (g *Graph) AddEdge(u, v int) {
	g.child[u] = append(g.child[u], v)
}

func (g *Graph) AddUndirectedEdge(u, v int) {
	g.AddEdge(u, v)
	g.AddEdge(v, u)
}

func (g *Graph) Print() {
	for node, neighbors := range g.child {
		fmt.Printf("%d -> %v\n", node, neighbors)
	}
}

func (g *Graph) DFS(start int, visited map[int]bool) {
	if visited[start] {
		return
	}
	fmt.Print(start, " ")
	visited[start] = true
	for _, neighbor := range g.child[start] {
		g.DFS(neighbor, visited)
	}
}

func (g *Graph) BFS(start int) {
	visited := make(map[int]bool)
	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		fmt.Print(node, " ")

		for _, neighbor := range g.child[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}

func main() {

}
