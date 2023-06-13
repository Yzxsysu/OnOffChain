package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func Subscribe(ip string) {
	var err error
	err = http.ListenAndServe(ip, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}

func SendData(msg interface{}, ip string, port string, path string) {
	//log.Println("Before SendData", path)
	u := url.URL{Scheme: "http", Host: ip + ":" + port, Path: path}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Println("SendData func json err:", err)
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("SendData err:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println("resp err:", err)
	}
	//log.Println("After SendData", path)
}

//func SendData(msg interface{}, ip string, port string, path string) {
//	u := url.URL{Scheme: "http", Host: ip + ":" + port, Path: path}
//	log.Println(u.String())
//	log.Println(msg, ip, port, path)
//	jsonData, err := json.Marshal(msg)
//	if err != nil {
//		log.Println("SendData func json err:", err)
//	}
//	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
//	req.Header.Set("Connection", "keep-alive")
//	if err != nil {
//		log.Println("SendData err:", err)
//	}
//	req.Header.Set("Content-Type", "application/json")
//
//	client := http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Println("resp err:", resp)
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Println("resp.Body close err:", err)
//		}
//	}(resp.Body)
//
//	respData, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Println("respData ReadAll err", err)
//	}
//	log.Println(string(respData))
//}

func WSHandlerTx(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		w.Header().Set("Connection", "keep-alive")
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		Tx := make([]SmallBankTransaction, 0)
		err = json.Unmarshal(body, &Tx)
		if err != nil {
			fmt.Println(err)
		}
		Txs <- Tx
		_, err = w.Write([]byte("off ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	//log.Println("Tx send to channel")
}

func WSHandlerS(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		w.Header().Set("Connection", "keep-alive")
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		Sub := make([][]GraphEdge, 0)
		err = json.Unmarshal(body, &Sub)
		if err != nil {
			fmt.Println(err)
		}
		MsgS <- Sub
		_, err = w.Write([]byte("off Sub ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	//log.Println("Sub send to channel")
}

func WSHandlerSV(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		w.Header().Set("Connection", "keep-alive")
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		SubV := make([][]uint16, 0)
		err = json.Unmarshal(body, &SubV)
		if err != nil {
			fmt.Println(err)
		}
		MsgSV <- SubV
		_, err = w.Write([]byte("off SubV ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	//log.Println("SubV send to channel")
}
