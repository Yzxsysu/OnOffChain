package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"log"
	"os"
)

var homeDir, isLeader, remotePorts string
var localPort uint

func init() {
	flag.StringVar(&homeDir, "home", "", "Path to the tendermint config directory (if empty, uses $HOME/.tendermint)")
	flag.StringVar(&isLeader, "leader", "false", "Is it a leader (default: false)")
	flag.UintVar(&localPort, "inport", 10057, "beacon chain rpc port")
	flag.StringVar(&remotePorts, "outport", "20057,21057", "shards chain rpc port")
}

func main() {
	// Parse command-line arguments
	flag.Parse()
	if homeDir == "" {
		homeDir = os.ExpandEnv("$HOME/.tendermint")
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
	/*gf, err := types.GenesisDocFromFile(config.GenesisFile())
	if err != nil {
		log.Fatalf("Loading genesis document: %v", err)
	}*/

	/*dbPath := filepath.Join(homeDir, "badger")
	db, err := application.NewBlockchainState()
	if err != nil {
		log.Fatalf("Opening database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Closing database: %v", err)
		}
	}()*/

}
