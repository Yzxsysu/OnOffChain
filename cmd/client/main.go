package main

import (
	"fmt"
	"net/http"
	"onffchain/smallbankapplication/application"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

// tx format: 127.0.0.1:20057/broadcast_tx_commit?tx="T=3,I=1,F=1,O=3,B=156>T=1,I=2,F=2,O=1,B=190"
func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		_ = fmt.Errorf("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		_ = fmt.Errorf("could not start CPU profile: ", err)
	}
	//runtime.GC()
	defer pprof.StopCPUProfile()

	f1, err := os.Create("mem.prof")
	if err != nil {
		_ = fmt.Errorf("could not create memory profile: ", err)
	}
	//runtime.GC()
	if err := pprof.WriteHeapProfile(f1); err != nil {
		fmt.Errorf("could not write memory profile: ", err)
	}
	f1.Close()

	f2, err := os.Create("goroutine.prof")
	if err != nil {
		_ = fmt.Errorf("could not create goroutine profile: ", err)
	}
	//runtime.GC()
	if gProf := pprof.Lookup("goroutine"); gProf == nil {
		fmt.Errorf("could not write goroutine: ")
	} else {
		gProf.WriteTo(f2, 0)
	}
	f2.Close()
	i := 0
	for {
		i++
		if i == 10 {
			break
		}
		time.Sleep(time.Millisecond * 100)
		var err error
		txs := application.GenerateTx(1000, 1000, 1)

		/*result, _ := json.Marshal(txs)

		sm := make([]application.SmallBankTransaction, 0)
		err := json.Unmarshal(result, &sm)
		fmt.Println(len(result))
		fmt.Println(string(result))
		fmt.Println(sm)*/
		str := ""
		l := len(txs)
		for i, tx := range txs {
			str += "T=" + strconv.Itoa(int(tx.T))
			str += "," + "I=" + strconv.Itoa(int(tx.I))
			str += "," + "F=" + string(tx.F)
			str += "," + "O=" + string(tx.O)
			str += "," + "B=" + strconv.Itoa(tx.B)
			if i != l-1 {
				str += ">"
			}
		}
		//requestBody := []byte(str)
		/*resp, err := http.Post("http://127.0.0.1:20057/broadcast_tx_commit", "application/json", bytes.NewBuffer(result))
		fmt.Println(resp)
		fmt.Println(len(result))
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(result)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)*/
		go func(str string) {
			fmt.Println(len(str))
			request1 := "127.0.0.1:20057/broadcast_tx_commit?tx=\"" + str + "\""
			_, err = http.Get("http://" + request1)
			if err != nil {
				fmt.Println(err)
			}
		}(str)
		/*fmt.Println(len(str))
		request1 := "127.0.0.1:20057/broadcast_tx_commit?tx=\"" + str + "\""
		//request1 := "127.0.0.1:20057/broadcast_tx_commit?tx=\"" + str + "\""
		_, err = http.Get("http://" + request1)
		if err != nil {
			fmt.Println(err)
		}*/
	}
}
