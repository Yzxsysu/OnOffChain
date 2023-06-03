package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var accountNum uint

var bufferLength = 10000
var Txs = make(chan []SmallBankTransaction, bufferLength)
var MsgS = make(chan [][]GraphEdge, bufferLength)
var MsgSV = make(chan [][]uint16, bufferLength)
var mV1 = make(chan map[string]AccountVersion, bufferLength)
var mV2 = make(chan map[string]AccountVersion, bufferLength)
var mV3 = make(chan map[string]AccountVersion, bufferLength)
var mV4 = make(chan map[string]AccountVersion, bufferLength)
var mV5 = make(chan map[string]AccountVersion, bufferLength)
var mV6 = make(chan map[string]AccountVersion, bufferLength)

var mSave = make(map[string]int)
var mCheck = make(map[string]int)
var webIp, webPort string
var groupIp1, groupIp2, groupIp3 []string
var groupPort1, groupPort2, groupPort3 []string
var SetNum string
var offChainIp, offChainPort string

func init() {
	flag.UintVar(&accountNum, "accountNum", 1000, "The account num of the SmallBank")
	flag.StringVar(&webIp, "webIp", "127.0.0.1,127.0.0.1,127.0.0.1", "send message websocket ip")
	flag.StringVar(&webPort, "webPort", "10157,10257,10357", "send message port")
	flag.StringVar(&offChainIp, "offChainIp", "127.0.0.1", "the ip of offchain node")
	flag.StringVar(&offChainPort, "offChainPort", "8090", "the port of offchain node")
	flag.StringVar(&SetNum, "SetNum", "2f", "Group Num")
}

// 监听proposal的Sub和SubV即可
func main() {
	flag.Parse()
	groupIp1, groupIp2, groupIp3 = SplitToThree(webIp)
	groupPort1, groupPort2, groupPort3 = SplitToThree(webPort)

	CreateAccountNum(int(accountNum))

	http.HandleFunc("/S", WSHandlerS)
	http.HandleFunc("/SV", WSHandlerSV)
	http.HandleFunc("/Tx", WSHandlerTx)
	//go http.ListenAndServe()
	go Validate()
	go Subscribe(offChainIp + ":" + offChainPort)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func Validate() {
	for {
		s := <-Txs
		log.Println("s := <- Txs")
		// init mChan
		mS := <-MsgS
		log.Println("mS := <-MsgS")
		mSV := <-MsgSV
		log.Println("mSV := <-MsgSV")
		go OValidate(&s, &mS, 0, mV1)
		go OValidate(&s, &mS, 1, mV2)
		go OValidate(&s, &mS, 2, mV3)
		go OVValidate(&s, &mSV, 0, mV4)
		go OVValidate(&s, &mSV, 1, mV5)
		go OVValidate(&s, &mSV, 2, mV6)
		go Merge()
	}
}

func Merge() {
	m1 := <-mV1
	m2 := <-mV2
	m3 := <-mV3
	m4 := <-mV4
	m5 := <-mV5
	m6 := <-mV6
	if SetNum == "2f" {
		go func() {
			for i, port := range groupPort3 {
				go SendData(m2, groupIp3[i], port, "/mV2")
				go SendData(m5, groupIp3[i], port, "/mV5")
			}
		}()
		go func() {
			for i, port := range groupPort2 {
				go SendData(m1, groupIp2[i], port, "/mV1")
				go SendData(m4, groupIp2[i], port, "/mV4")
			}
		}()
		go func() {
			for i, port := range groupPort1 {
				go SendData(m3, groupIp1[i], port, "/mV3")
				go SendData(m6, groupIp1[i], port, "/mV6")
			}
		}()
	} else if SetNum == "f" {
		go func() {
			for i, port := range groupPort3 {
				go SendData(m1, groupIp3[i], port, "/mV1")
				go SendData(m2, groupIp3[i], port, "/mV2")
				go SendData(m4, groupIp3[i], port, "/mV4")
				go SendData(m5, groupIp3[i], port, "/mV5")
			}
		}()
		go func() {
			for i, port := range groupPort2 {
				go SendData(m1, groupIp2[i], port, "/mV1")
				go SendData(m3, groupIp2[i], port, "/mV3")
				go SendData(m4, groupIp2[i], port, "/mV4")
				go SendData(m6, groupIp2[i], port, "/mV6")
			}
		}()
		go func() {
			for i, port := range groupPort1 {
				go SendData(m2, groupIp1[i], port, "/mV2")
				go SendData(m3, groupIp1[i], port, "/mV3")
				go SendData(m5, groupIp1[i], port, "/mV5")
				go SendData(m6, groupIp1[i], port, "/mV6")
			}
		}()
	}
	tempCheckVersion := make(map[string]uint16)
	tempSaveVersion := make(map[string]uint16)
	for key, value := range m1 {
		_, err := tempSaveVersion[key]
		if err != true {
			tempSaveVersion[key] = value.SaveVersion
		}
		_, err = tempCheckVersion[key]
		if err != true {
			tempSaveVersion[key] = value.CheckVersion
		}
		if tempSaveVersion[key] < value.SaveVersion {
			tempSaveVersion[key] = value.SaveVersion
			mSave[key] = value.Save
		}
		if tempCheckVersion[key] < value.CheckVersion {
			tempCheckVersion[key] = value.CheckVersion
			mCheck[key] = value.Check
		}
	}
	for key, value := range m2 {
		_, err := tempSaveVersion[key]
		if err != true {
			tempSaveVersion[key] = value.SaveVersion
		}
		_, err = tempCheckVersion[key]
		if err != true {
			tempSaveVersion[key] = value.CheckVersion
		}
		if tempSaveVersion[key] < value.SaveVersion {
			tempSaveVersion[key] = value.SaveVersion
			mSave[key] = value.Save
		}
		if tempCheckVersion[key] < value.CheckVersion {
			tempCheckVersion[key] = value.CheckVersion
			mCheck[key] = value.Check
		}
	}
	for key, value := range m3 {
		_, err := tempSaveVersion[key]
		if err != true {
			tempSaveVersion[key] = value.SaveVersion
		}
		_, err = tempCheckVersion[key]
		if err != true {
			tempSaveVersion[key] = value.CheckVersion
		}
		if tempSaveVersion[key] < value.SaveVersion {
			tempSaveVersion[key] = value.SaveVersion
			mSave[key] = value.Save
		}
		if tempCheckVersion[key] < value.CheckVersion {
			tempCheckVersion[key] = value.CheckVersion
			mCheck[key] = value.Check
		}
	}
	for key, value := range m4 {
		mSave[key] = value.Save
		mCheck[key] = value.Check
	}
	for key, value := range m5 {
		mSave[key] = value.Save
		mCheck[key] = value.Check
	}
	for key, value := range m6 {
		mSave[key] = value.Save
		mCheck[key] = value.Check
	}
}

func SplitToThree(ports string) ([]string, []string, []string) {
	p := strings.Split(ports, ",")
	l := len(p)
	return p[:l/3], p[l/3 : 2*l/3], p[2*l/3:]
}
