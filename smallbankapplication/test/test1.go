package main

import (
	"fmt"
)

/*func main() {
	m := make(map[int]int)
	for i := 0; i < 100; i++ {
		WriteMap(m, i, i)
	}
	for i := 0; i < 100; i++ {
		go ReadMap(m, i)
	}
	// map支持并发读
	time.Sleep(time.Second * 2)
	var txResult TxResult
	txResult = NewTxResult()
	txResult.AccountName = append(txResult.AccountName, "A")
	txResult.ConsistentSaveValue = append(txResult.ConsistentCheckValue, 0)
	var SaveBool bool
	fmt.Println(SaveBool)
	fmt.Println(txResult)
}*/

func counter(out chan<- int) {
	for x := 0; x < 100; x++ {
		out <- x
	}
	//close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

/*
	func main() {
		naturals := make(chan int)
		go counter(naturals)
		for i := 0; i < 100; i++ {
			fmt.Println(<-naturals)
		}
		//squares := make(chan int)
		//go squarer(squares, naturals)
		//printer(squares)
	}
*/
func WriteMap(m map[int]int, i int, j int) {
	m[i] = j
}

func ReadMap(m map[int]int, i int) {
	fmt.Println(m[i])
}

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

func NewTxResult() TxResult {
	return TxResult{
		PreTxId:              make([]uint16, 0),
		CurrentTxId:          0,
		AccountName:          make([]string, 0),
		CheckVersion:         make([]uint16, 0),
		ConsistentCheckValue: make([]int, 0),
		SaveVersion:          make([]uint16, 0),
		ConsistentSaveValue:  make([]int, 0),
	}
}
