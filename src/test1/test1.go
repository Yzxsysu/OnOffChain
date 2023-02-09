package main

import (
	"bytes"
	"io"
	"net/http"
)

func main() {
	requestBody := []byte("tx=fromid=0,toid=1,type=1,from=ABCD,to=DCBA,value=20,data=DATA,nonce=1")
	resp, err := http.Post("http://127.0.0.1:20057/broadcast_tx_commit", "application/x-www-form-urlencoded", bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
}
