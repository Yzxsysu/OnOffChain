package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Request struct {
	Message string `json:"message"`
}

func main() {
	// 监听请求
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// 读取请求内容
		buf := new(bytes.Buffer)
		_, err := io.Copy(buf, r.Body)
		if err != nil {
			log.Fatalf("read request body failed: %s\n", err)
		}
		body := buf.Bytes()

		// 将请求内容反序列化为Request结构体
		var req [][]uint16
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Fatalf("unmarshal request body failed: %s\n", err)
		}

		fmt.Printf("received message: %v\n", req)
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("start http server failed: %s\n", err)
	}
}
