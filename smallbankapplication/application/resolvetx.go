package application

import (
	"bytes"
	"log"
	"strconv"
)

func (BCstate *BlockchainState) ResolveAndExecuteTx(request []byte, threadNUm int) {
	TxType := []byte("TxType")
	TxId := []byte("TxId")
	From := []byte("From")
	To := []byte("To")
	Balance := []byte("Balance")

	// TxType=3,TxId=1,From=1,To=3,Balance=156>TxType=1,TxId=2,From=2,To=1,Balance=190"
	txs := bytes.Split(request, []byte(">"))
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
			switch kv[0] {
			case TxType:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.TxType = uint8(temp)
			case TxId:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.TxId = uint16(temp)
			case From:
				tx.From = make([]byte, len(kv[1]))
				copy(tx.From, kv[1])
			case To:
				tx.To = make([]byte, len(kv[1]))
				copy(tx.To, kv[1])
			case Balance:
				// temp_value := string(kv[1])
				tx.Balance = BytesToInt(kv[1])
			}
		}
		ReceiveTx[i] = tx
	}
	BCstate.ExecuteSmallBankTransaction(ReceiveTx, threadNUm)
}
