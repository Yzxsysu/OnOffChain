package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func HttpPost(url string, data interface{}, headers map[string]string) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func SendData(msg *[]byte, ip string, port string, path string) {
	u := url.URL{Scheme: "http", Host: ip + ":" + port, Path: path}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(*msg))
	if err != nil {
		log.Println("SendData err:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("resp err:", resp)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("resp.Body close err:", err)
		}
	}(resp.Body)

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("respData ReadAll err", err)
	}
	log.Println(string(respData))
}

func main() {
	url := "http://127.0.0.1:8080/Tx"

	/*postBody, err := json.Marshal(map[string]string{
		"title":  "foo",
		"body":   "bar",
		"userId": "1",
	})
	if err != nil {
		log.Println(err)
	}*/
	p := map[string]string{
		"title":  "a",
		"body":   "b",
		"userId": "v",
	}
	pjson, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}

	arr := make([][]uint16, 5)
	for i := range arr {
		arr[i] = make([]uint16, 3)
	}
	arrjson, err := json.Marshal(arr)
	if err != nil {
		log.Println(err)
	}

	go SendData(&pjson, "127.0.0.1", "8090", "/")
	go SendData(&arrjson, "127.0.0.1", "8090", "/SV")
	time.Sleep(time.Second)
	post, err := HttpPost(url, p, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(post))

	/*responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	log.Printf(string(body))*/
}
