package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	abciclient "github.com/tendermint/tendermint/abci/client"
	cfg "github.com/tendermint/tendermint/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/types"
	"log"
	smallbankapplication "onffchain/smallbankapplication/abci"
	"onffchain/smallbankapplication/application"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

var homeDir, isLeader, remotePorts string
var localPort, group, accountNum, coreNum uint

func init() {
	flag.StringVar(&homeDir, "home", "", "Path to the tendermint config directory (if empty, uses $HOME/.tendermint)")
	flag.StringVar(&isLeader, "leader", "false", "Is it a leader (default: false)")
	flag.UintVar(&accountNum, "accountNum", 1000, "The account num of the SmallBank")
	flag.UintVar(&group, "group", 0, "The group that the node belongs to")
	flag.UintVar(&coreNum, "coreNum", 8, "control the num of cpu's cores")
	flag.UintVar(&localPort, "inport", 10057, "beacon chain rpc port")
	flag.StringVar(&remotePorts, "outport", "20057,21057", "shards chain rpc port")
}

func main() {
	// Parse command-line arguments
	flag.Parse()
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
