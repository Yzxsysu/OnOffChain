package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	abciclient "github.com/tendermint/tendermint/abci/client"
	cfg "github.com/tendermint/tendermint/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/types"
	"log"
	"net/http"
	"net/url"
	smallbankapplication "onffchain/smallbankapplication/abci"
	"onffchain/smallbankapplication/application"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

var homeDir, isLeader, remotePorts string
var localPort, groupNum, accountNum, coreNum, group uint
var leaderIp string
var OffChainIp string
var SetNum string

func init() {
	flag.StringVar(&homeDir, "home", "", "Path to the tendermint config directory (if empty, uses $HOME/.tendermint)")
	flag.StringVar(&isLeader, "leader", "false", "Is it a leader (default: false)")
	flag.StringVar(&leaderIp, "leaderIp", "0.0.0.0:26657", "Let replica subscribe the websocket")
	flag.UintVar(&accountNum, "accountNum", 1000, "The account num of the SmallBank")
	flag.UintVar(&groupNum, "groupNum", 3, "the total number of group")
	flag.StringVar(&OffChainIp, "OffChainIp", "0.0.0.0:8090", "Let replica subscribe the websocket")
	flag.UintVar(&group, "group", 1, "The group that the node belongs to")
	flag.UintVar(&coreNum, "coreNum", 8, "control the num of cpu's cores")
	flag.UintVar(&localPort, "inport", 10057, "beacon chain rpc port")
	flag.StringVar(&remotePorts, "outport", "20057,21057", "shards chain rpc port")
	flag.StringVar(&SetNum, "SetNum", "2f", "Group Num")
	smallbankapplication.US = url.URL{Scheme: "ws", Host: leaderIp, Path: "/S"}
	smallbankapplication.USV = url.URL{Scheme: "ws", Host: leaderIp, Path: "/SV"}
	smallbankapplication.OffChianURL = url.URL{Scheme: "ws", Host: leaderIp, Path: "/Tx"}
	application.MsgS = make(chan [][]application.GraphEdge)
	application.MsgSV = make(chan [][]uint16)
	application.MsgV1 = make(chan map[string]application.AccountVersion)
	application.MsgV2 = make(chan map[string]application.AccountVersion)
	application.MsgV3 = make(chan map[string]application.AccountVersion)
	application.MsgV4 = make(chan map[string]application.AccountVersion)
	application.MsgV5 = make(chan map[string]application.AccountVersion)
	application.MsgV6 = make(chan map[string]application.AccountVersion)
	application.Version = make(chan map[string]application.AccountVersion)
}

func main() {
	// Parse command-line arguments
	flag.Parse()
	application.SetNum = SetNum
	application.Group = int(group)
	// Allocate MsgS and MsgSV memory
	// Subscribe to the leader websocket using goroutine
	// define the url
	UmV1 := url.URL{Scheme: "ws", Host: OffChainIp, Path: "/mV1"}
	UmV2 := url.URL{Scheme: "ws", Host: OffChainIp, Path: "/mV2"}
	UmV3 := url.URL{Scheme: "ws", Host: OffChainIp, Path: "/mV3"}
	UmV4 := url.URL{Scheme: "ws", Host: OffChainIp, Path: "/mV4"}
	UmV5 := url.URL{Scheme: "ws", Host: OffChainIp, Path: "/mV5"}
	UmV6 := url.URL{Scheme: "ws", Host: OffChainIp, Path: "/mV6"}
	go http.HandleFunc(smallbankapplication.US.String(), WSHandlerS)
	go http.HandleFunc(smallbankapplication.USV.String(), WSHandlerSV)

	if homeDir == "" {
		homeDir = os.ExpandEnv("/home/.tendermint")
	}
	// Set default path and arguments
	config := cfg.DefaultConfig()

	// Set root(the location of tendermint node)
	config.SetRoot(homeDir)

	// Set viper working location
	viper.SetConfigFile(fmt.Sprintf("%s/%s", homeDir, "config/config.toml"))

	// Read in config into viper
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Reading config: %v", err)
	}
	// Unmarshal unmarshals the config into a Struct. Make sure that the tags
	// on the fields of the structure are properly set.
	// It can change the default setting of the config -> config := cfg.DefaultConfig()
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Decoding config: %v", err)
	}
	// After unmarshaling the config into the config struct, need to validate
	// whether the config struct is a proper form
	if err := config.ValidateBasic(); err != nil {
		log.Fatalf("Invalid configuration data: %v", err)
	}
	// GenesisFile is a string type in config struct
	gf, err := types.GenesisDocFromFile(config.GenesisFile())
	if err != nil {
		log.Fatalf("Loading genesis document: %v", err)
	}

	dbPath := filepath.Join(homeDir, "badger")
	db, err, err1 := application.NewBlockchainState("badgerdb", true, dbPath)

	if err != nil || err1 != nil {
		log.Fatalf("Opening database: %v, %v", err, err1)
	}
	defer func() {
		if err := db.SavingStore.Close(); err != nil {
			log.Fatalf("Closing database: %v", err)
		}
		if err := db.CheckingStore.Close(); err != nil {
			log.Fatalf("Closing database: %v", err)
		}
	}()
	if isLeader == "true" {
		db.Leader = true
	} else if isLeader == "false" {
		go ListenOffChain(OffChainIp)
		if SetNum == "2f" {
			if group == 1 {
				go http.HandleFunc(UmV3.String(), WSHandler3)
				go http.HandleFunc(UmV6.String(), WSHandler6)
			}
			if group == 2 {
				go http.HandleFunc(UmV1.String(), WSHandler1)
				go http.HandleFunc(UmV4.String(), WSHandler4)
			}
			if group == 3 {
				go http.HandleFunc(UmV2.String(), WSHandler2)
				go http.HandleFunc(UmV5.String(), WSHandler5)
			}
		} else if SetNum == "f" {
			if group == 1 {
				go http.HandleFunc(UmV2.String(), WSHandler2)
				go http.HandleFunc(UmV3.String(), WSHandler3)
				go http.HandleFunc(UmV5.String(), WSHandler5)
				go http.HandleFunc(UmV6.String(), WSHandler6)
			}
			if group == 2 {
				go http.HandleFunc(UmV1.String(), WSHandler1)
				go http.HandleFunc(UmV3.String(), WSHandler3)
				go http.HandleFunc(UmV4.String(), WSHandler4)
				go http.HandleFunc(UmV6.String(), WSHandler6)
			}
			if group == 3 {
				go http.HandleFunc(UmV1.String(), WSHandler1)
				go http.HandleFunc(UmV2.String(), WSHandler2)
				go http.HandleFunc(UmV4.String(), WSHandler4)
				go http.HandleFunc(UmV5.String(), WSHandler5)
			}
		}

		db.Leader = false
	}

	// create account
	for i := 0; i < int(accountNum); i++ {
		db.CreateAccount(strconv.Itoa(i), 1000, 1000)
	}

	app := smallbankapplication.NewSmallBankApplication(db)
	acc := abciclient.NewLocalCreator(app)

	logger := tmlog.MustNewDefaultLogger(tmlog.LogFormatPlain, tmlog.LogLevelInfo, false)
	node, err := nm.New(config, logger, acc, gf)

	if err != nil {
		log.Fatalf("Creating node: %v", err)
	}
	go Subscribe(leaderIp)
	err = node.Start()
	if err != nil {
		log.Fatalf("Starting node: %v", err)
	}
	defer func() {
		err = node.Stop()
		if err != nil {
			log.Fatalf("Stoping node: %v", err)
		}
		node.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
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

func ListenOffChain(i string) {
	err := http.ListenAndServe(i+"/mV1", nil)
	if err != nil {
		fmt.Println(err)
	}
	err = http.ListenAndServe(i+"/mV2", nil)
	if err != nil {
		fmt.Println(err)
	}
	err = http.ListenAndServe(i+"/mV3", nil)
	if err != nil {
		fmt.Println(err)
	}
	err = http.ListenAndServe(i+"/mV4", nil)
	if err != nil {
		fmt.Println(err)
	}
	err = http.ListenAndServe(i+"/mV5", nil)
	if err != nil {
		fmt.Println(err)
	}
	err = http.ListenAndServe(i+"/mV6", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func WSHandler1(w http.ResponseWriter, r *http.Request) {
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
	V1 := make(map[string]application.AccountVersion, 0)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		err = json.Unmarshal(msg, &V1)
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		application.MsgV1 <- V1
		log.Printf("Read from leader msg:%s \n", msg)
	}
}

func WSHandler2(w http.ResponseWriter, r *http.Request) {
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
	V2 := make(map[string]application.AccountVersion, 0)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		err = json.Unmarshal(msg, &V2)
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		application.MsgV2 <- V2
		log.Printf("Read from leader msg:%s \n", msg)
	}
}

func WSHandler3(w http.ResponseWriter, r *http.Request) {
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
	V3 := make(map[string]application.AccountVersion, 0)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		err = json.Unmarshal(msg, &V3)
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		application.MsgV3 <- V3
		log.Printf("Read from leader msg:%s \n", msg)
	}
}

func WSHandler4(w http.ResponseWriter, r *http.Request) {
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
	V4 := make(map[string]application.AccountVersion, 0)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		err = json.Unmarshal(msg, &V4)
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		application.MsgV4 <- V4
		log.Printf("Read from leader msg:%s \n", msg)
	}
}

func WSHandler5(w http.ResponseWriter, r *http.Request) {
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
	V5 := make(map[string]application.AccountVersion, 0)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		err = json.Unmarshal(msg, &V5)
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		application.MsgV5 <- V5
		log.Printf("Read from leader msg:%s \n", msg)
	}
}

func WSHandler6(w http.ResponseWriter, r *http.Request) {
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
	V6 := make(map[string]application.AccountVersion, 0)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		err = json.Unmarshal(msg, &V6)
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		application.MsgV6 <- V6
		log.Printf("Read from leader msg:%s \n", msg)
	}
}
