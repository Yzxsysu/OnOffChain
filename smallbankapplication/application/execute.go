package application

import (
	"bytes"
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	"log"
	"strconv"
)

const (
	GetBalance    uint8 = 1
	Amalgamate    uint8 = 2
	UpdateBalance uint8 = 3
	UpdateSaving  uint8 = 4
	SendPayment   uint8 = 5
	WriteCheck    uint8 = 6
	AddrLength    uint8 = 4
	DataLength    uint8 = 4
)

func Serialization(tx []byte) (uint32, []SmallBankTransaction) {
	// TxJson.Tx_type = 4
	Txs := bytes.Split(tx, []byte(">"))
	txNum := len(Txs) - 1
	ReceiveTx := make([]SmallBankTransaction, len(Txs)-1)
	for i := 0; i < txNum; i++ {
		tx := Txs[i]
		txElements := bytes.Split(tx, []byte(","))
		var TxJson SmallBankTransaction
		for _, txElement := range txElements {
			kv := bytes.Split(txElement, []byte("="))
			switch {
			case string(kv[0]) == "TxType":
				temp_type64, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				TxJson.TxType = uint8(temp_type64)
			case string(kv[0]) == "From":
				TxJson.From = make([]byte, len(kv[1]))
				copy(TxJson.From, kv[1])
			case string(kv[0]) == "To":
				TxJson.To = make([]byte, len(kv[1]))
				copy(TxJson.To, kv[1])
			case string(kv[0]) == "Balance":
				// temp_value := string(kv[1])
				TxJson.Balance = BytesToInt(kv[1])
			}
		}
		ReceiveTx[i] = TxJson
	}
	return 0, ReceiveTx
}

type SmallBankTransaction struct {
	TxType  uint8
	TxId    uint16
	From    []byte
	To      []byte
	Balance int
}

func (BCstate *BlockchainState) ExecuteSmallBankTransaction(s []SmallBankTransaction, threadNUm int) {
	Num := len(s)
	var TxType uint8
	var TxId uint16
	var From []byte
	var To []byte
	var Balance int

	txResult := make(chan TxResult, Num)

	for i := 0; i < Num; i++ {
		tx := s[i]
		TxType = tx.TxType
		TxId = tx.TxId
		From = tx.From
		To = tx.To
		Balance = tx.Balance

		switch TxType {
		case GetBalance:
			go BCstate.GetBalance(TxId, string(From), txResult)
		case Amalgamate:
			go BCstate.Amalgamate(TxId, string(From), string(To), txResult)
		case UpdateBalance:
			go BCstate.UpdateBalance(TxId, string(From), Balance, txResult)
		case UpdateSaving:
			go BCstate.UpdateSaving(TxId, string(From), Balance, txResult)
		case SendPayment:
			go BCstate.SendPayment(TxId, string(From), string(To), Balance, txResult)
		case WriteCheck:
			go BCstate.WriteCheck(TxId, string(From), Balance, txResult)
		default:
			fmt.Println("TxType doesn't match")
		}
	}
	l := len(txResult)
	pq := queue.NewPriorityQueue(l, true)
	visited := make([]bool, l)
	GenerateGraph(txResult, pq, visited)
}

func GenerateGraph(txResult <-chan TxResult, pq *queue.PriorityQueue, visited []bool) map[uint16][]Vertex {
	m := make(map[uint16][]Vertex)
	for t := range txResult {
		num := len(t.AccountName)
		for i := 0; i < num; i++ {
			if t.PreTxId[i] == 0 {
				// if the pre tx id is zero
				// weight is 0
				/*e := NewEdge()
				e.To = t.CurrentTxId
				err := pq.Put(e)
				if err != nil {
					log.Println(err)
				}*/
				break
			}
			v := NewVertex()
			e := NewEdge()
			v.CurrentTxId = t.CurrentTxId

			pre := t.PreTxId[i]
			e.From = pre
			e.To = t.CurrentTxId
			visited[pre] = true
			visited[v.CurrentTxId] = true

			v.AccountName = t.AccountName[i]

			v.CheckBool = t.CheckBool[i]
			if v.CheckBool {
				v.CheckVersion = t.CheckVersion[i]
				v.ConsistentCheckValue = t.ConsistentCheckValue[i]
				e.Weight += 12
			}
			v.SaveBool = t.SaveBool[i]
			if v.SaveBool {
				v.SaveVersion = t.SaveVersion[i]
				v.ConsistentSaveValue = t.ConsistentSaveValue[i]
				e.Weight += 12
			}
			m[pre] = append(m[pre], v)
			err := pq.Put(e)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return m
}

func CutGraph(m map[uint16][]Vertex, pq *queue.PriorityQueue, group int, visited []bool) ([][]GraphEdge, [][]uint16) {
	// init sub graph
	sub := make([][]GraphEdge, 0)

	l := len(visited)
	s := make([]uint16, 0)
	for i := 0; i < l; i++ {
		if visited[i] == false {
			s = append(s, uint16(i))
		}
	}
	SubV := SplitSlice(s, group)

	// init num of the group
	groupNUm := 1
	g := make([]GraphEdge, 0)

	// init the num of tx
	txNum := 1
	for !pq.Empty() {
		if groupNUm == group {
			qLen := pq.Len()
			for j := 0; j < qLen; j++ {
				edge := pq.Peek().(Edge)

				// from -> v1, v2
				length := len(m[edge.From])

				// v Vertex, maybe v1, v2
				v := m[edge.From]
				for i := 0; i < length; i++ {
					if v[i].CurrentTxId == edge.To {
						GE := NewGraphEdge()
						GE.From = edge.From
						GE.To = edge.To
						s := "name=" + v[i].AccountName
						if v[i].SaveVersion != 0 {
							s += ",SaveVersion=" + strconv.Itoa(int(v[i].SaveVersion))
							s += ",ConsistentSaveValue=" + strconv.Itoa(int(v[i].ConsistentSaveValue))
						}
						if v[i].CheckVersion != 0 {
							s += ",CheckVersion=" + strconv.Itoa(int(v[i].CheckVersion))
							s += ",ConsistentCheckValue=" + strconv.Itoa(int(v[i].ConsistentCheckValue))
						}
						GE.data = s
						g = append(g, GE)
						break
					}
				}
			}
			sub = append(sub, g)
			break
		}

		edge := pq.Peek().(Edge)
		// from -> v1, v2
		length := len(m[edge.From])
		// v Vertex, maybe v1, v2
		v := m[edge.From]
		// when txNum is zero
		if groupNUm != group && txNum >= l/group {
			// reset txNum
			txNum = 1

			// groupNum ++
			groupNUm++
			for i := 0; i < length; i++ {
				if v[i].CurrentTxId == edge.To {
					GE := NewGraphEdge()
					GE.From = edge.From
					GE.To = edge.To
					s := "name=" + v[i].AccountName
					if v[i].SaveVersion != 0 {
						s += ",SaveVersion=" + strconv.Itoa(int(v[i].SaveVersion))
						s += ",ConsistentSaveValue=" + strconv.Itoa(int(v[i].ConsistentSaveValue))
					}
					if v[i].CheckVersion != 0 {
						s += ",CheckVersion=" + strconv.Itoa(int(v[i].CheckVersion))
						s += ",ConsistentCheckValue=" + strconv.Itoa(int(v[i].ConsistentCheckValue))
					}
					GE.data = s
					g = append(g, GE)
					break
				}
			}
			sub = append(sub, g)
			g = make([]GraphEdge, 0)
		} else {
			for i := 0; i < length; i++ {
				if v[i].CurrentTxId == edge.To {
					GE := NewGraphEdge()
					GE.From = edge.From
					GE.To = edge.To
					txNum++
					g = append(g, GE)
					break
				}
			}
		}
	}

	return sub, SubV
}

func SplitSlice(slice []uint16, parts int) [][]uint16 {
	var result [][]uint16
	partLen := len(slice) / parts
	for i := 0; i < parts; i++ {
		start := i * partLen
		end := (i + 1) * partLen
		if i == parts-1 {
			end = len(slice)
		}
		result = append(result, slice[start:end])
	}
	return result
}
