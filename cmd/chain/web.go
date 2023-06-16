package main

import (
	"encoding/json"
	"fmt"
	"github.com/Yzxsysu/onoffchain/smallbankapplication/application"
	"io"
	"log"
	"net/http"
)

func Subscribe(ip string) {
	var err error
	err = http.ListenAndServe(ip, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
	return
}

func WSHandlerS(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		w.Header().Set("Connection", "keep-alive")
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		Sub := make([][]application.GraphEdge, 0)
		err = json.Unmarshal(body, &Sub)
		if err != nil {
			fmt.Println(err)
		}
		//log.Println("Sub:", Sub)
		application.MsgS <- Sub
		//log.Println("send to channel")
		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
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
		//log.Println("SubV:", SubV)
		application.MsgSV <- SubV
		//log.Println("send to channel")
		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func WSHandler1(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		w.Header().Set("Connection", "keep-alive")
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		V1 := make(map[string]application.AccountVersion, 0)
		err = json.Unmarshal(body, &V1)
		if err != nil {
			fmt.Println(err)
		}
		//log.Println("V1:", V1)
		application.MsgV1 <- V1
		//log.Println("send to channel")
		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func WSHandler2(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		w.Header().Set("Connection", "keep-alive")
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		V2 := make(map[string]application.AccountVersion, 0)
		err = json.Unmarshal(body, &V2)
		if err != nil {
			fmt.Println(err)
		}
		//log.Println("V2:", V2)
		application.MsgV2 <- V2
		//log.Println("send to channel")
		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func WSHandler3(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		w.Header().Set("Connection", "keep-alive")
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		V3 := make(map[string]application.AccountVersion, 0)
		err = json.Unmarshal(body, &V3)
		if err != nil {
			fmt.Println(err)
		}
		//log.Println("V3:", V3)
		application.MsgV3 <- V3
		//log.Println("send to channel")
		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func WSHandler4(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		w.Header().Set("Connection", "keep-alive")
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		V4 := make(map[string]application.AccountVersion, 0)
		err = json.Unmarshal(body, &V4)
		if err != nil {
			fmt.Println(err)
		}
		//log.Println("V4:", V4)
		application.MsgV4 <- V4
		//log.Println("send to channel")
		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func WSHandler5(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		w.Header().Set("Connection", "keep-alive")
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		V5 := make(map[string]application.AccountVersion, 0)
		err = json.Unmarshal(body, &V5)
		if err != nil {
			fmt.Println(err)
		}
		//log.Println("V5:", V5)
		application.MsgV5 <- V5
		//log.Println("send to channel")
		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func WSHandler6(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		w.Header().Set("Connection", "keep-alive")
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		V6 := make(map[string]application.AccountVersion, 0)
		err = json.Unmarshal(body, &V6)
		if err != nil {
			fmt.Println(err)
		}
		//log.Println("V6:", V6)
		application.MsgV6 <- V6
		//log.Println("send to channel")
		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
