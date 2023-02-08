package main

import (
	"flag"
	"fmt"
	"time"
)

type Graph struct {
	v       uint16
	adj     [][]uint16
	visited []bool
	order   []uint16
}

func NewGraph(v uint16) *Graph {
	g := &Graph{
		v:       v,
		adj:     make([][]uint16, v),
		visited: make([]bool, v),
	}
	return g
}

func (g *Graph) AddEdge(s uint16, t uint16) {
	g.adj[s] = append(g.adj[s], t)
}

func (g *Graph) TopoSortByDFS() []uint16 {
	for i := 0; uint16(i) < g.v; i++ {
		if !g.visited[i] {
			g.DFS(uint16(i))
		}
	}
	return g.order
}

func (g *Graph) DFS(vertex uint16) {
	g.visited[vertex] = true
	for _, v := range g.adj[vertex] {
		if !g.visited[v] {
			g.DFS(v)
		}
	}
	g.order = append(g.order, vertex)
}

func read(ch chan bool) {
	ch <- true
}

func write(ch chan bool) {
	time.Sleep(0)
	ch <- true
}

func main() {
	g := NewGraph(6)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)
	order := g.TopoSortByDFS()
	fmt.Println(order)
	ch := make(chan bool, 2)
	go write(ch)
	go read(ch)
	fmt.Println(<-ch, <-ch)
	// 定义一个字符串标志name，默认值为"world"，使用说明为"the target to greet"
	name := flag.String("name", "world", "the target to greet")

	// 解析命令行参数
	flag.Parse()

	// 使用解析后的命令行参数
	fmt.Println("Hello,", *name)
}
