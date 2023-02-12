package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func receiveData(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		log.Fatalf("read request body failed: %s\n", err)
	}
	body := buf.Bytes()
	if err != nil {
		log.Println("Error reading the request body:", err)
		return
	}
	fmt.Println(string(body))
	// 处理接收到的信息
	// ...

	// 回复信息
	resp := []byte(`{"message": "Hello, Client!"}`)

	_, err = w.Write(resp)
	if err != nil {
		fmt.Println(err)
	}

}

func startReceiver() {
	http.HandleFunc("/api", receiveData)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	// 启动服务端监听
	startReceiver()
}
