package application

import (
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
	//AddrLength    uint8 = 4
	//DataLength    uint8 = 4
)

/*func (BCstate *BlockchainState) ExecuteSmallBankTransaction(s []SmallBankTransaction, threadNUm int) {
	length := len(s)
	var TxType uint8
	var TxId uint16
	var From []byte
	var To []byte
	var Balance int

	txResult := make(chan TxResult, length)

	for i := 0; i < length; i++ {
		tx := s[i]
		TxType = tx.T
		TxId = tx.I
		From = tx.F
		To = tx.O
		Balance = tx.B

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
			fmt.Println("T doesn't match")
		}
	}

	pq := queue.NewPriorityQueue(length, true)
	visited := make([]bool, length)
	GenerateGraph(txResult, pq, visited, length)
}*/

func GenerateGraph(txResult <-chan TxResult, pq *queue.PriorityQueue, visited []bool, length int) map[uint16][]Vertex {
	m := make(map[uint16][]Vertex)
	for l := 0; l < length; l++ {
		t := <-txResult
		num := len(t.AccountName)
		// tx1 -> tx3, tx2 -> tx3
		for i := 0; i < num; i++ {
			if t.PreTxId[i] == 0 {
				// if the pre tx id is zero
				// weight is 0
				/*e := NewEdge()
				e.O = t.CurrentTxId
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
	//log.Println(m)
	return m
	//for t := range txResult {
	//	num := len(t.AccountName)
	//	for i := 0; i < num; i++ {
	//		if t.PreTxId[i] == 0 {
	//			// if the pre tx id is zero
	//			// weight is 0
	//			/*e := NewEdge()
	//			e.O = t.CurrentTxId
	//			err := pq.Put(e)
	//			if err != nil {
	//				log.Println(err)
	//			}*/
	//			break
	//		}
	//		v := NewVertex()
	//		e := NewEdge()
	//		v.CurrentTxId = t.CurrentTxId
	//
	//		pre := t.PreTxId[i]
	//		e.From = pre
	//		e.To = t.CurrentTxId
	//		visited[pre] = true
	//		visited[v.CurrentTxId] = true
	//
	//		v.AccountName = t.AccountName[i]
	//
	//		v.CheckBool = t.CheckBool[i]
	//		if v.CheckBool {
	//			v.CheckVersion = t.CheckVersion[i]
	//			v.ConsistentCheckValue = t.ConsistentCheckValue[i]
	//			e.Weight += 12
	//		}
	//		v.SaveBool = t.SaveBool[i]
	//		if v.SaveBool {
	//			v.SaveVersion = t.SaveVersion[i]
	//			v.ConsistentSaveValue = t.ConsistentSaveValue[i]
	//			e.Weight += 12
	//		}
	//		m[pre] = append(m[pre], v)
	//		err := pq.Put(e)
	//		if err != nil {
	//			log.Println(err)
	//		}
	//	}
	//}
}

func CutGraph(m map[uint16][]Vertex, pq *queue.PriorityQueue, group int, visited []bool) ([][]GraphEdge, [][]uint16) {
	// init sub graph
	sub := make([][]GraphEdge, 0)
	var SubV [][]uint16
	// 1001
	l := len(visited)
	s := make([]uint16, 0)
	// 0 - 1000 (1001个) -> 1 - 1000 | i < l -> 1000 < 1001暂停
	// 1 to l - 1才是需要遍历的
	go func() {
		for i := 1; i < l; i++ {
			if visited[i] == false {
				s = append(s, uint16(i))
			}
		}
		SubV = SplitSlice(s, group)
	}()

	// init num of the group
	groupNUm := 1
	g := make([]GraphEdge, 0)

	// init the num of tx
	txNum := 1
	for !pq.Empty() {
		if groupNUm == group {
			qLen := pq.Len()
			for j := 0; j < qLen; j++ {
				edges, err := pq.Get(1)
				if err != nil {
					log.Println(err)
				}
				edge := edges[0].(Edge)
				//edge := pq.Peek().(Edge)

				// from -> v1, v2
				length := len(m[edge.From])

				// v Vertex, maybe v1, v2
				v := m[edge.From]
				/*for i := 0; i < length; i++ {
					if v[i].CurrentTxId == edge.To {
						GE := NewGraphEdge()
						GE.F = edge.From
						GE.T = edge.To
						edgeData := "name=" + v[i].AccountName
						if v[i].SaveVersion != 0 {
							edgeData += ",SaveVersion=" + strconv.Itoa(int(v[i].SaveVersion))
							edgeData += ",ConsistentSaveValue=" + strconv.Itoa(int(v[i].ConsistentSaveValue))
						}
						if v[i].CheckVersion != 0 {
							edgeData += ",CheckVersion=" + strconv.Itoa(int(v[i].CheckVersion))
							edgeData += ",ConsistentCheckValue=" + strconv.Itoa(int(v[i].ConsistentCheckValue))
						}
						GE.D = edgeData
						g = append(g, GE)
						break
					}
				}*/
				for i := 0; i < length; i++ {
					if v[i].CurrentTxId == edge.To {
						GE := NewGraphEdge()
						GE.F = edge.From
						GE.T = edge.To
						txNum++
						g = append(g, GE)
						break
					}
				}
			}
			sub = append(sub, g)
			break
		}

		edges, err := pq.Get(1)
		if err != nil {
			log.Println(err)
		}
		edge := edges[0].(Edge)
		//edge := pq.Peek().(Edge)
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
					GE.F = edge.From
					GE.T = edge.To
					edgeData := "name=" + v[i].AccountName
					if v[i].SaveVersion != 0 {
						edgeData += ",SaveVersion=" + strconv.Itoa(int(v[i].SaveVersion))
						edgeData += ",ConsistentSaveValue=" + strconv.Itoa(int(v[i].ConsistentSaveValue))
					}
					if v[i].CheckVersion != 0 {
						edgeData += ",CheckVersion=" + strconv.Itoa(int(v[i].CheckVersion))
						edgeData += ",ConsistentCheckValue=" + strconv.Itoa(int(v[i].ConsistentCheckValue))
					}
					GE.D = edgeData
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
					GE.F = edge.From
					GE.T = edge.To
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
