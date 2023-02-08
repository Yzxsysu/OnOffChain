package app

import (
	"fmt"
	"sync"
)

// TxResult records the
type TxResult struct {
	PreTxId              []uint16
	CurrentTxId          uint16
	AccountName          []string
	CheckVersion         []uint16
	ConsistentCheckValue []int
	SaveVersion          []uint16
	ConsistentSaveValue  []int
}

var SaveMap sync.Map
var CheckMap sync.Map

type AccountData struct {
	WrittenBy uint16
	Version   uint16
}

func main() {
	var m sync.Map
	m.Store("a", AccountData{1, 1})
	val, _ := m.Load("a")
	WrittenBy, _ := val.(AccountData)

	fmt.Println(WrittenBy.WrittenBy)
	fmt.Println(val)
	//fmt.Println(m.Load("a"))
}
