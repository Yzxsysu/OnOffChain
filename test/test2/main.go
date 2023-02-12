package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var dataQueue [][]byte

// 发送端代码
func sendData(url string, data []byte) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// 接收端代码
func receiveData(w http.ResponseWriter, r *http.Request) {

	// 将接收到的数据放入队列
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
	if len(dataQueue) >= 10 {
		time.Sleep(time.Second)
	}
	// 将数据加入队列
	dataQueue = append(dataQueue, body)
	_, err = fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	if err != nil {
		return
	}
}

// 启动接收端
func startReceiver() {
	http.HandleFunc("/api", receiveData)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// 处理队列中的数据
func processData() {
	//time.Sleep(time.Second)
	for {
		// 从队列中取出数据
		if len(dataQueue) > 0 {
			data := dataQueue[0]
			//dataQueue = dataQueue[1:]
			dataQueue = dataQueue[:copy(dataQueue, dataQueue[1:])]
			fmt.Println(string(data))
		}
	}

}

func main() {
	dataQueue = make([][]byte, 3)
	// 启动接收端
	go startReceiver()

	// 启动数据处理线程
	go processData()
	i := 0
	// 发送数据
	for {
		str := fmt.Sprintf(`%v, {"message": "Hello, World!"}`, i)
		i++
		data := []byte(str)
		err := sendData("http://localhost:8080/api", data)
		if err != nil {
			log.Println("Error sending data:", err)
		}
		//time.Sleep(time.Second)
	}
}
