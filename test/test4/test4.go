package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func sendData(url string, data []byte) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(buf)
	return nil
}

func main() {
	data := []byte(`{"message": "Hello, World!"}`)
	err := sendData("http://localhost:8080/api", data)
	if err != nil {
		log.Println("Error sending data:", err)
	}
}
