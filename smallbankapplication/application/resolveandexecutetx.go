package application

import (
	"bytes"
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	"log"
	"strconv"
	"sync"
)

const (
	TxTypeString  string = "T"
	TxIdString    string = "I"
	FromString    string = "F"
	ToString      string = "O"
	BalanceString string = "B"
)

func ResolveTx(request []byte) []SmallBankTransaction {
	log.Println("Before ReceiveTx:")
	txs := bytes.Split(request, []byte(">"))
	l := len(txs)
	if l == 0 {
		log.Println("the tx is nil")
	}
	ReceiveTx := make([]SmallBankTransaction, l)
	/*err := json.Unmarshal(*request, &ReceiveTx)
	if err != nil {
		log.Println(err)
	}*/
	for i, elements := range txs {
		var tx SmallBankTransaction
		element := bytes.Split(elements, []byte(","))
		for _, e := range element {
			kv := bytes.Split(e, []byte("="))
			switch string(kv[0]) {
			case TxTypeString:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.T = uint8(temp)
			case TxIdString:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.I = uint16(temp)
			case FromString:
				tx.F = make([]byte, len(kv[1]))
				copy(tx.F, kv[1])
			case ToString:
				tx.O = make([]byte, len(kv[1]))
				copy(tx.O, kv[1])
			case BalanceString:
				// temp_value := string(kv[1])
				tx.B = BytesToInt(kv[1])
			}
		}
		ReceiveTx[i] = tx
	}
	log.Println("After ReceiveTx:")
	return ReceiveTx
}

func (BCstate *BlockchainState) ResolveAndExecuteTx(request *[]byte) ([][]GraphEdge, [][]uint16, []SmallBankTransaction) {
	// T=3,I=1,F=1,O=3,B=156>T=1,I=2,F=2,O=1,B=190"
	txs := bytes.Split(*request, []byte(">"))
	l := len(txs)
	if l == 0 {
		log.Println("the tx is nil")
	}
	ReceiveTx := make([]SmallBankTransaction, l)
	for i, elements := range txs {
		var tx SmallBankTransaction
		element := bytes.Split(elements, []byte(","))
		for _, e := range element {
			kv := bytes.Split(e, []byte("="))
			switch string(kv[0]) {
			case TxTypeString:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.T = uint8(temp)
			case TxIdString:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.I = uint16(temp)
			case FromString:
				tx.F = make([]byte, len(kv[1]))
				copy(tx.F, kv[1])
			case ToString:
				tx.O = make([]byte, len(kv[1]))
				copy(tx.O, kv[1])
			case BalanceString:
				// temp_value := string(kv[1])
				tx.B = BytesToInt(kv[1])
			}
		}
		ReceiveTx[i] = tx
	}

	var TxType uint8
	var TxId uint16
	var From []byte
	var To []byte
	var Balance int

	l = len(ReceiveTx)
	txResult := make(chan TxResult, l)
	AccountDataMap = sync.Map{}
	for i := 0; i < l; i++ {
		tx := ReceiveTx[i]
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

	pq := queue.NewPriorityQueue(l, true)
	visited := make([]bool, l+1)
	Sub, SubV := CutGraph(GenerateGraph(txResult, pq, visited, l), pq, 3, visited)
	close(txResult)
	// need to send to on chain and other off chain
	return Sub, SubV, ReceiveTx
}

func (BCstate *BlockchainState) ResolveAndExecuteTxWithSyncMap(request *[]byte) ([][]GraphEdge, [][]uint16, []SmallBankTransaction) {
	// T=3,I=1,F=1,O=3,B=156>T=1,I=2,F=2,O=1,B=190"
	txs := bytes.Split(*request, []byte(">"))
	l := len(txs)
	if l == 0 {
		log.Println("the tx is nil")
	}

	log.Println("Before ReceiveTx")
	ReceiveTx := make([]SmallBankTransaction, l)
	for i, elements := range txs {
		var tx SmallBankTransaction
		element := bytes.Split(elements, []byte(","))
		for _, e := range element {
			kv := bytes.Split(e, []byte("="))
			switch string(kv[0]) {
			case TxTypeString:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.T = uint8(temp)
			case TxIdString:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.I = uint16(temp)
			case FromString:
				tx.F = make([]byte, len(kv[1]))
				copy(tx.F, kv[1])
			case ToString:
				tx.O = make([]byte, len(kv[1]))
				copy(tx.O, kv[1])
			case BalanceString:
				// temp_value := string(kv[1])
				tx.B = BytesToInt(kv[1])
			}
		}
		ReceiveTx[i] = tx
	}
	log.Println("After ReceiveTx")
	var TxType uint8
	var TxId uint16
	var From []byte
	var To []byte
	var Balance int

	log.Println("Before txResult")
	txResult := make(chan TxResult, l)
	sMap := sync.Map{}
	for i := 0; i < l; i++ {
		tx := ReceiveTx[i]
		TxType = tx.T
		TxId = tx.I
		From = tx.F
		To = tx.O
		Balance = tx.B

		switch TxType {
		case GetBalance:
			go BCstate.GetBalanceWithSyncMap(TxId, string(From), txResult, &sMap)
		case Amalgamate:
			go BCstate.AmalgamateWithSyncMap(TxId, string(From), string(To), txResult, &sMap)
		case UpdateBalance:
			go BCstate.UpdateBalanceWithSyncMap(TxId, string(From), Balance, txResult, &sMap)
		case UpdateSaving:
			go BCstate.UpdateSavingWithSyncMap(TxId, string(From), Balance, txResult, &sMap)
		case SendPayment:
			go BCstate.SendPaymentWithSyncMap(TxId, string(From), string(To), Balance, txResult, &sMap)
		case WriteCheck:
			go BCstate.WriteCheckWithSyncMap(TxId, string(From), Balance, txResult, &sMap)
		default:
			fmt.Println("T doesn't match")
		}
	}
	log.Println("After txResult")

	pq := queue.NewPriorityQueue(l, true)
	visited := make([]bool, l+1)
	log.Println("Before CutGraph")
	Sub, SubV := CutGraph(GenerateGraph(txResult, pq, visited, l), pq, 3, visited)
	log.Println("After CutGraph")
	close(txResult)
	// need to send to on chain and other off chain
	return Sub, SubV, ReceiveTx
}
