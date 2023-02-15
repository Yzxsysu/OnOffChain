package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	smallbankapplication "onffchain/smallbankapplication/abci"
	"onffchain/smallbankapplication/application"
	"os"
	"os/signal"
	"syscall"
)

var leaderIp string
var Txs chan []SmallBankTransaction
var MsgS chan [][]GraphEdge
var MsgSV chan [][]uint16
var accountNum uint

var mV1 chan map[string]AccountVersion
var mV2 chan map[string]AccountVersion
var mV3 chan map[string]AccountVersion
var mV4 chan map[string]AccountVersion
var mV5 chan map[string]AccountVersion
var mV6 chan map[string]AccountVersion

var mSave map[string]int
var mCheck map[string]int

func init() {
	flag.UintVar(&accountNum, "accountNum", 1000, "The account num of the SmallBank")
	flag.StringVar(&leaderIp, "leaderIp", "0.0.0.0:26657", "Let replica subscribe the websocket")
}

// 监听proposal的Sub和SubV即可
func main() {
	flag.Parse()

	mV1 = make(chan map[string]AccountVersion)
	mV2 = make(chan map[string]AccountVersion)
	mV3 = make(chan map[string]AccountVersion)
	mV4 = make(chan map[string]AccountVersion)
	mV5 = make(chan map[string]AccountVersion)
	mV6 = make(chan map[string]AccountVersion)
	// Create account the same state
	CreateAccountNum(int(accountNum))

	MsgS = make(chan [][]GraphEdge)
	MsgSV = make(chan [][]uint16)
	Txs = make(chan []SmallBankTransaction)
	// Subscribe to the leader websocket using goroutine
	// define the url
	smallbankapplication.US = url.URL{Scheme: "ws", Host: leaderIp, Path: "/S"}
	smallbankapplication.USV = url.URL{Scheme: "ws", Host: leaderIp, Path: "/SV"}
	LeaderIp := url.URL{Scheme: "ws", Host: leaderIp, Path: "/Tx"}
	go http.HandleFunc(smallbankapplication.US.String(), WSHandlerS)
	go http.HandleFunc(smallbankapplication.USV.String(), WSHandlerSV)
	go http.HandleFunc(LeaderIp.String(), WSHandlerTx)
	//go Execute()
	go Validate()
	go Subscribe(leaderIp)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

// Execute simulate the execute
func Validate() {
	for {
		if len(Txs) != 0 && len(MsgS) != 0 && len(MsgSV) != 0 {
			s := <-Txs
			// init mChan
			mS := <-MsgS
			mSV := <-MsgSV
			go OValidate(&s, mS, 0, mV1)
			go OValidate(&s, mS, 1, mV2)
			go OValidate(&s, mS, 2, mV3)
			go OVValidate(&s, mSV, 0, mV4)
			go OVValidate(&s, mSV, 1, mV5)
			go OVValidate(&s, mSV, 2, mV6)
			go Merge()
		}
	}
}

func Merge() {
	m1 := <-mV1
	m2 := <-mV2
	m3 := <-mV3
	m4 := <-mV4
	m5 := <-mV5
	m6 := <-mV6
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
	jm1, err := json.Marshal(m1)
	if err != nil {
		return
	}
	mV1URL := url.URL{Scheme: "ws", Host: "localhost:8090", Path: "/mV1"}
	go SendData(&jm1, mV1URL)

	jm2, err := json.Marshal(m2)
	if err != nil {
		fmt.Println(err)
	}
	mV2URL := url.URL{Scheme: "ws", Host: "localhost:8090", Path: "/mV2"}
	go SendData(&jm2, mV2URL)

	jm3, err := json.Marshal(m3)
	if err != nil {
		fmt.Println(err)
	}
	mV3URL := url.URL{Scheme: "ws", Host: "localhost:8090", Path: "/mV3"}
	go SendData(&jm3, mV3URL)

	jm4, err := json.Marshal(m4)
	if err != nil {
		fmt.Println(err)
	}
	mV4URL := url.URL{Scheme: "ws", Host: "localhost:8090", Path: "/mV4"}
	go SendData(&jm4, mV4URL)

	jm5, err := json.Marshal(m5)
	if err != nil {
		fmt.Println(err)
	}
	mV5URL := url.URL{Scheme: "ws", Host: "localhost:8090", Path: "/mV5"}
	go SendData(&jm5, mV5URL)

	jm6, err := json.Marshal(m6)
	if err != nil {
		fmt.Println(err)
	}
	mV6URL := url.URL{Scheme: "ws", Host: "localhost:8090", Path: "/mV6"}
	go SendData(&jm6, mV6URL)
}

func SendData(msg *[]byte, u url.URL) {
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Dial error:", err)
		return
	}
	// close conn
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("conn close error:", err)
		}
	}(conn)

	if err := conn.WriteMessage(websocket.BinaryMessage, *msg); err != nil {
		//if err := conn.WriteMessage(1, []byte("今天。。。"));err != nil {
		log.Println("Writeing error...", err)
		return
	}
	return
}
func Subscribe(ip string) {
	var err error
	err = http.ListenAndServe(ip+"/S", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
	err = http.ListenAndServe(ip+"/SV", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
	err = http.ListenAndServe(ip+"/Tx", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}

var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024 * 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024 * 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WSHandlerTx(w http.ResponseWriter, r *http.Request) {
	//   完成握手 升级为 WebSocket长连接，使用conn发送和接收消息。
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	//调用连接的WriteMessage和ReadMessage方法以一片字节发送和接收消息。实现如何回显消息：
	Sub := make([]SmallBankTransaction, 0)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		err = json.Unmarshal(msg, &Sub)
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		Txs <- Sub
		log.Printf("Read from leader msg:%s \n", msg)
	}
}

func WSHandlerS(w http.ResponseWriter, r *http.Request) {
	//   完成握手 升级为 WebSocket长连接，使用conn发送和接收消息。
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	//调用连接的WriteMessage和ReadMessage方法以一片字节发送和接收消息。实现如何回显消息：
	Sub := make([][]application.GraphEdge, 0)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		err = json.Unmarshal(msg, &Sub)
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		application.MsgS <- Sub
		log.Printf("Read from leader msg:%s \n", msg)
	}
}

func WSHandlerSV(w http.ResponseWriter, r *http.Request) {
	//   完成握手 升级为 WebSocket长连接，使用conn发送和接收消息。
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	//调用连接的WriteMessage和ReadMessage方法以一片字节发送和接收消息。实现如何回显消息：
	SubV := make([][]uint16, 0)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		err = json.Unmarshal(msg, &SubV)
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		application.MsgSV <- SubV
		log.Printf("Read from leader msg:%s \n", msg)
	}
}
