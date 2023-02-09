package main

import (
	"fmt"
	dbm "onffchain/src/tm-db"
	"os"
	"path/filepath"
)

func main() {
	// Open a BadgerDB database
	homeDir := os.ExpandEnv("/home/.tendermint")
	dbPath := filepath.Join(homeDir, "badger")
	db, err := dbm.NewDB("mydb", dbm.BoltDBBackend, dbPath)
	if err != nil {
		panic(err)
	}
	fmt.Println(1)
	defer db.Close()

	// Use the database as needed
	// ...
}
