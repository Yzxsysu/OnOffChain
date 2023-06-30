package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	smallbankapplication "github.com/Yzxsysu/onoffchain/smallbankapplication/abci"
	"github.com/Yzxsysu/onoffchain/smallbankapplication/application"
	"github.com/spf13/viper"
	abciclient "github.com/tendermint/tendermint/abci/client"
	cfg "github.com/tendermint/tendermint/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/types"
)

var homeDir, isLeader string
var groupNum, accountNum, coreNum, group uint
var leaderIp string
var OffChainIp string
var OffChainPort string
var SetNum string

var webIp, webPort string
var subscribeIp string

func init() {
	flag.StringVar(&homeDir, "home", "", "Path to the tendermint config directory (if empty, uses $HOME/.tendermint)")
	flag.StringVar(&isLeader, "leader", "false", "Is it a leader (default: false)")
	flag.StringVar(&leaderIp, "leaderIp", "172.17.0.3:26657", "Let replica subscribe the websocket")
	flag.UintVar(&accountNum, "accountNum", 1000, "The account num of the SmallBank")
	flag.UintVar(&groupNum, "groupNum", 3, "the total number of group")
	flag.StringVar(&OffChainIp, "OffChainIp", "172.17.0.100", "Let replica subscribe the websocket")
	flag.StringVar(&OffChainPort, "OffChainPort", "8090", "Let replica subscribe the websocket")
	flag.UintVar(&group, "group", 1, "The group that the node belongs to")
	flag.UintVar(&coreNum, "coreNum", 8, "control the num of cpu's cores")
	// send message port sends the leader's graphs to validators, begin with 1
	flag.StringVar(&webIp, "webIp", "172.17.0.3,172.17.0.4,172.17.0.5", "send message websocket ip")
	flag.StringVar(&webPort, "webPort", "10157,10257,10357", "send message websocket ports")
	flag.StringVar(&SetNum, "SetNum", "2f", "Group Num")
	flag.StringVar(&subscribeIp, "subscribeIp", "", "the replica that listens and serves ip")
}

func main() {
	flag.Parse()
	// Set the core num
	runtime.GOMAXPROCS(int(coreNum))
	application.ByteLen = 1024
	application.CycleNum = 10
	// application.ByteLen = 0
	// application.CycleNum = 0
	application.SetNum = SetNum
	application.Group = int(group)
	smallbankapplication.Ips = strings.Split(webIp, ",")
	smallbankapplication.Ports = strings.Split(webPort, ",")
	smallbankapplication.OffChainIp = OffChainIp
	smallbankapplication.OffChainPort = OffChainPort
	http.HandleFunc("/S", WSHandlerS)
	http.HandleFunc("/SV", WSHandlerSV)
	http.HandleFunc("/mV1", WSHandler1)
	http.HandleFunc("/mV2", WSHandler2)
	http.HandleFunc("/mV3", WSHandler3)
	http.HandleFunc("/mV4", WSHandler4)
	http.HandleFunc("/mV5", WSHandler5)
	http.HandleFunc("/mV6", WSHandler6)
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
			log.Fatalf("Closing SavingStore database: %v", err)
		}
		if err := db.CheckingStore.Close(); err != nil {
			log.Fatalf("Closing Checkingstore database: %v", err)
		}
	}()
	if isLeader == "true" {
		db.Leader = true
	} else if isLeader == "false" {
		db.Leader = false
		if subscribeIp != "" {
			go Subscribe(subscribeIp)
		}
	}

	// create account
	for i := 1; i <= int(accountNum); i++ {
		db.CreateAccount(strconv.Itoa(i), 1000, 1000)
	}

	app := smallbankapplication.NewSmallBankApplication(db)
	acc := abciclient.NewLocalCreator(app)
	logger := tmlog.MustNewDefaultLogger(tmlog.LogFormatPlain, tmlog.LogLevelInfo, false)
	node, err := nm.New(config, logger, acc, gf)

	if err != nil {
		log.Fatalf("Creating node: %v", err)
	}

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
