package main

import (
	"fmt"
	"log"
	"net/http"
	"onffchain/smallbankapplication/application"
	"os"
	"strconv"
	"time"
)

// tx format: 127.0.0.1:20057/broadcast_tx_commit?tx="T=3,I=1,F=1,O=3,B=156>T=1,I=2,F=2,O=1,B=190"
func main() {
	// 先删除
	err := os.Remove("client_tx.log") // 指定文件路径及名称
	if err != nil {
		// 如果发生错误，则打印错误信息
		panic(err)
	}
	file, err := os.OpenFile("client_tx.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(0)
	log.SetOutput(file)
	defer file.Close()
	txs := application.GenerateTx(1000, 100000, 0.1)
	//fmt.Println(txs)
	for {
		var err error
		str := ""
		//txs := application.GenerateTx(1000, 1000, 1)

		l := len(txs)
		fmt.Println("l =", l)
		for i, tx := range txs {
			currentTime := time.Now()
			t := currentTime.Nanosecond()
			str += "T=" + strconv.Itoa(int(tx.T))
			str += "," + "I=" + strconv.Itoa(int(tx.I))
			str += "," + "F=" + string(tx.F)
			str += "," + "O=" + string(tx.O)
			str += "," + "B=" + strconv.Itoa(tx.B)
			str += "," + "t=" + strconv.Itoa(t)
			if i != l-1 {
				str += ">"
			}
		}
		go func(str string) {
			//fmt.Println(len(str))
			request1 := "172.172.0.3:26657/broadcast_tx_commit?tx=\"" + str + "\""
			_, err = http.Get("http://" + request1)
			log.Println(time.Now().Format("2006-01-02T15:04:05+08:00"))
			if err != nil {
				fmt.Println(err)
			}
		}(str)
		time.Sleep(time.Millisecond * 230)
	}
}
