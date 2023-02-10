package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	url := "http://localhost:8000/api"

	// 设置请求内容
	//requestBody := []byte("{\"message\":\"Hello world\"}")
	u := make([][]uint16, 0)
	a := make([]uint16, 10)
	u = append(u, a)

	for i := 0; i < 100; i++ {
		u[0][0] = uint16(i)
		b, err := json.Marshal(u)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 发送请求
		_, err = http.Post(url, "application/json", bytes.NewBuffer(b))
		if err != nil {
			log.Fatalf("http post failed: %s\n", err)
		}

		// 关闭连接
		/*defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)*/
	}
}
