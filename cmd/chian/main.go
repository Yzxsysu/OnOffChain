package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"onffchain/smallbankapplication/application"
	"strconv"
)

// tx format: 127.0.0.1:20057/broadcast_tx_commit?tx="TxType=3,TxId=1,From=1,To=3,Balance=156>TxType=1,TxId=2,From=2,To=1,Balance=190"
func main() {
	txs := application.GenerateTx(1000, 1000, 1)
	str := ""
	l := len(txs)
	for i, tx := range txs {
		str += "TxType=" + strconv.Itoa(int(tx.TxType))
		str += "," + "TxId=" + strconv.Itoa(int(tx.TxId))
		str += "," + "From=" + string(tx.From)
		str += "," + "To=" + string(tx.To)
		str += "," + "Balance=" + strconv.Itoa(tx.Balance)
		if i != l-1 {
			str += ">"
		}
	}
	requestBody := []byte(str)
	resp, err := http.Post("http://127.0.0.1:20057/broadcast_tx_commit", "application/x-www-form-urlencoded", bytes.NewBuffer(requestBody))
	//fmt.Println(str)
	fmt.Println(len(requestBody))
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	/*fmt.Println(len(str))
	request1 := "127.0.0.1:20057/broadcast_tx_commit?tx=\"" + str + "\""
	//request1 := "127.0.0.1:20057/broadcast_tx_commit?tx=\"" + str + "\""
	_, err = http.Get("http://" + request1)
	if err != nil {
		fmt.Println(err)
	}*/

}
