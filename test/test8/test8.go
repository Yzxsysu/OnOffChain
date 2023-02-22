package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"onffchain/smallbankapplication/application"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
			}
			re := make(map[string]string)
			err = json.Unmarshal(body, &re)
			if err != nil {
				fmt.Println(err)
			}
			log.Println(re)
			_, err = w.Write([]byte("Received a POST request"))
			if err != nil {
				log.Println(err)
			}
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})
	go http.HandleFunc("/S", WSHandlerSV)
	go Subscribe("127.0.0.1:8090")
	//log.Fatal(http.ListenAndServe(":8090", nil))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func Subscribe(ip string) {
	var err error
	err = http.ListenAndServe(ip, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}

func WSHandlerSV(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Receive msg
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		SubV := make([][]uint16, 0)
		err = json.Unmarshal(body, &SubV)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Sub:", SubV)
		application.MsgSV <- SubV
		a := <-application.MsgSV
		fmt.Println(a)
		fmt.Println("send to channel")
		_, err = w.Write([]byte("ok"))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
