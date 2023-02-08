package application

import (
	"github.com/Workiva/go-datastructures/queue"
	"net"
	"strconv"
	"sync"
)

// TxResult records the
type TxResult struct {
	PreTxId              []uint16
	CurrentTxId          uint16
	AccountName          []string
	CheckBool            []bool
	CheckVersion         []uint16
	ConsistentCheckValue []int
	SaveBool             []bool
	SaveVersion          []uint16
	ConsistentSaveValue  []int
}

func NewTxResult() TxResult {
	return TxResult{
		PreTxId:              make([]uint16, 0),
		CurrentTxId:          0,
		AccountName:          make([]string, 0),
		CheckBool:            make([]bool, 0),
		CheckVersion:         make([]uint16, 0),
		ConsistentCheckValue: make([]int, 0),
		SaveBool:             make([]bool, 0),
		SaveVersion:          make([]uint16, 0),
		ConsistentSaveValue:  make([]int, 0),
	}
}

var AccountDataMap sync.Map

type AccountData struct {
	WrittenBy    uint16
	CheckVersion uint16
	SaveVersion  uint16
}

func NewAccountData() AccountData {
	return AccountData{
		WrittenBy:    0,
		CheckVersion: 0,
		SaveVersion:  0,
	}
}

type IpAddress struct {
	Ip   net.IP
	Port uint16
}

func IntToBytes(n int) []byte {
	return []byte(strconv.Itoa(n))
}

func BytesToInt(b []byte) int {
	s := string(b)
	i, _ := strconv.Atoi(s)
	return i
}

type Edge struct {
	From   uint16
	To     uint16
	Weight uint8
}

func NewEdge() Edge {
	return Edge{}
}

func (p Edge) Compare(other queue.Item) int {
	otherPerson := other.(Edge)
	if p.Weight < otherPerson.Weight {
		return 1
	} else if p.Weight == otherPerson.Weight {
		return 0
	} else {
		return -1
	}
}

type Vertex struct {
	CurrentTxId          uint16
	AccountName          string
	CheckBool            bool
	CheckVersion         uint16
	ConsistentCheckValue int
	SaveBool             bool
	SaveVersion          uint16
	ConsistentSaveValue  int
}

func NewVertex() Vertex {
	return Vertex{}
}

type GraphEdge struct {
	From uint16
	To   uint16
	data string
}

func NewGraphEdge() GraphEdge {
	return GraphEdge{}
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

/*func main() {
	var m sync.Map
	m.Store("a", AccountData{1, 1, 1})
	val, _ := m.Load("a")
	WrittenBy, _ := val.(AccountData)

	fmt.Println(WrittenBy.WrittenBy)
	fmt.Println(val)

}
*/
