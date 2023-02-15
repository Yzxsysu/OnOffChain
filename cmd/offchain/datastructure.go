package main

const (
	GetBalance    uint8 = 1
	Amalgamate    uint8 = 2
	UpdateBalance uint8 = 3
	UpdateSaving  uint8 = 4
	SendPayment   uint8 = 5
	WriteCheck    uint8 = 6
	//AddrLength    uint8 = 4
	//DataLength    uint8 = 4
)

type SmallBankTransaction struct {
	T uint8
	I uint16
	F []byte
	O []byte
	B int
}

type AccountVersion struct {
	Check        int
	CheckVersion uint16
	Save         int
	SaveVersion  uint16
}

func NewAccountVersion() AccountVersion {
	return AccountVersion{}
}

type GraphEdge struct {
	F uint16
	T uint16
	D string
}

type SortedGraph struct {
	v       uint16
	adj     [][]uint16
	visited []bool
	order   []uint16
}

func NewSortedGraph(v uint16) *SortedGraph {
	g := &SortedGraph{
		v:       v,
		adj:     make([][]uint16, v),
		visited: make([]bool, v),
	}
	return g
}

func (g *SortedGraph) AddEdge(s uint16, t uint16) {
	g.adj[s] = append(g.adj[s], t)
}

func (g *SortedGraph) TopoSortByDFS() []uint16 {
	for i := 0; uint16(i) < g.v; i++ {
		if !g.visited[i] {
			g.DFS(uint16(i))
		}
	}
	return g.order
}

func (g *SortedGraph) DFS(vertex uint16) {
	g.visited[vertex] = true
	for _, v := range g.adj[vertex] {
		if !g.visited[v] {
			g.DFS(v)
		}
	}
	g.order = append(g.order, vertex)
}
