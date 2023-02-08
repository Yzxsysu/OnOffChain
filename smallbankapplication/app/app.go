package app

import (
	dbm "github.com/tendermint/tm-db"
	"log"
	"strconv"
	"sync"
)

type BackendType string

// These are valid backend types.
const (
	// GoLevelDBBackend represents goleveldb (github.com/syndtr/goleveldb - most
	// popular implementation)
	//   - pure go
	//   - stable
	GoLevelDBBackend BackendType = "goleveldb"
	// CLevelDBBackend represents cleveldb (uses levigo wrapper)
	//   - fast
	//   - requires gcc
	//   - use cleveldb build tag (go build -tags cleveldb)
	CLevelDBBackend BackendType = "cleveldb"
	// MemDBBackend represents in-memory key value store, which is mostly used
	// for testing.
	MemDBBackend BackendType = "memdb"
	// BoltDBBackend represents bolt (uses etcd's fork of bolt -
	// github.com/etcd-io/bbolt)
	//   - EXPERIMENTAL
	//   - may be faster is some use-cases (random reads - indexer)
	//   - use boltdb build tag (go build -tags boltdb)
	BoltDBBackend BackendType = "boltdb"
	// RocksDBBackend represents rocksdb (uses github.com/tecbot/gorocksdb)
	//   - EXPERIMENTAL
	//   - requires gcc
	//   - use rocksdb build tag (go build -tags rocksdb)
	RocksDBBackend BackendType = "rocksdb"

	BadgerDBBackend BackendType = "badgerdb"
)

type BlockchainState struct {
	// SmallBank contains two types of accounts
	CheckingStore dbm.DB
	SavingStore   dbm.DB
	Height        uint32
	Leader        bool
	AppHash       []byte
	// account -> pointer of a lock
	// the lock can lock separate instance of the struct
	CheckLock map[string]*Lock
	SaveLock  map[string]*Lock
}

// Lock create locks for specific check account
type Lock struct {
	lock sync.RWMutex
}

// NewLock constructs lock for each Lock instance
func NewLock() *Lock {
	return &Lock{
		lock: sync.RWMutex{},
	}
}

// NewBlockchainState can choose {goleveldb, cleveldb, memdb, boltdb , rocksdb, badgerdb}
// input Corresponding name
func NewBlockchainState(DBName string, leader bool, dir string) *BlockchainState {
	var BaseCaseState BlockchainState
	var err error
	// choose DB
	switch DBName {
	case string(GoLevelDBBackend):
		BaseCaseState.CheckingStore, err = dbm.NewDB(DBName, dbm.GoLevelDBBackend, dir)
		BaseCaseState.SavingStore, err = dbm.NewDB(DBName, dbm.GoLevelDBBackend, dir)
	case string(CLevelDBBackend):
		BaseCaseState.CheckingStore, err = dbm.NewDB(DBName, dbm.CLevelDBBackend, dir)
		BaseCaseState.SavingStore, err = dbm.NewDB(DBName, dbm.CLevelDBBackend, dir)
	case string(MemDBBackend):
		BaseCaseState.CheckingStore, err = dbm.NewDB(DBName, dbm.MemDBBackend, dir)
		BaseCaseState.SavingStore, err = dbm.NewDB(DBName, dbm.MemDBBackend, dir)
	case string(BoltDBBackend):
		BaseCaseState.CheckingStore, err = dbm.NewDB(DBName, dbm.BoltDBBackend, dir)
		BaseCaseState.SavingStore, err = dbm.NewDB(DBName, dbm.BoltDBBackend, dir)
	case string(RocksDBBackend):
		BaseCaseState.CheckingStore, err = dbm.NewDB(DBName, dbm.RocksDBBackend, dir)
		BaseCaseState.SavingStore, err = dbm.NewDB(DBName, dbm.RocksDBBackend, dir)
	case string(BadgerDBBackend):
		BaseCaseState.CheckingStore, err = dbm.NewDB(DBName, dbm.BadgerDBBackend, dir)
		BaseCaseState.SavingStore, err = dbm.NewDB(DBName, dbm.BadgerDBBackend, dir)
	}
	if err != nil {
		log.Fatalf("Create database error: %v", err)
	}
	BaseCaseState.Leader = leader
	return &BlockchainState{
		CheckingStore: BaseCaseState.CheckingStore,
		SavingStore:   BaseCaseState.SavingStore,
		Height:        1,
		Leader:        leader,
		CheckLock:     make(map[string]*Lock),
		SaveLock:      make(map[string]*Lock),
	}
}

func (BCstate *BlockchainState) DeliverGraph() {

}

// CreateAccount can create account with saving balance and checking balance
func (BCstate *BlockchainState) CreateAccount(AccountName string, SavingBalance int, CheckingBalance int) {
	// Create two separate accounts for two DB
	var err error
	err = BCstate.SavingStore.Set([]byte(AccountName), []byte(strconv.Itoa(SavingBalance)))
	if err != nil {
		panic(err)
	}
	err = BCstate.CheckingStore.Set([]byte(AccountName), []byte(strconv.Itoa(CheckingBalance)))
	if err != nil {
		panic(err)
	}

	// Create two separate locks for each account
	// Initialize the lock in BCstate
	BCstate.CheckLock[AccountName] = NewLock()
	BCstate.SaveLock[AccountName] = NewLock()
}

// GetBalance can get the balance of the account including check and save store
func (BCstate *BlockchainState) GetBalance(TxId uint16, AccountName string, txResult chan<- TxResult) int {
	// read only
	// lock the account
	BCstate.SaveLock[AccountName].lock.RLock()
	BCstate.CheckLock[AccountName].lock.RLock()
	// use defer to unlock
	defer BCstate.SaveLock[AccountName].lock.RUnlock()
	defer BCstate.CheckLock[AccountName].lock.RUnlock()

	// construct save account savedata
	var savedata AccountData
	var result TxResult
	result.AccountName = append(result.AccountName, AccountName)
	result.CurrentTxId = TxId
	value, ok := SaveMap.Load(AccountName)
	if ok {
		// type assertion AccountData
		savedata = value.(AccountData)
		result.PreTxId = append(result.PreTxId, savedata.WrittenBy)
		result.SaveVersion = append(result.SaveVersion, savedata.Version)
	} else {
		// if fail will return nil
		// saveaccountdata := value
	}
	Save, err := BCstate.SavingStore.Get([]byte(AccountName))
	if err != nil {
		log.Println(err)
	}
	SaveInt, err := strconv.Atoi(string(Save))
	// add consistent save value
	result.ConsistentSaveValue = append(result.ConsistentSaveValue, SaveInt)
	if err != nil {
		log.Println(err)
	}
	// construct save account savedata done

	// construct check account savedata
	var checkdata AccountData
	value, ok = CheckMap.Load(AccountName)
	if ok {
		checkdata = value.(AccountData)
		if checkdata.WrittenBy == result.PreTxId[0] {
			result.CheckVersion = append(result.CheckVersion, checkdata.Version)
		} else {
			result.PreTxId = append(result.PreTxId, checkdata.WrittenBy)

		}
	}
	Check, err := BCstate.CheckingStore.Get([]byte(AccountName))
	if err != nil {
		log.Println(err)
	}
	CheckInt, err := strconv.Atoi(string(Check))
	if err != nil {
		log.Println(err)
	}

	return SaveInt + CheckInt
}

func (BCstate *BlockchainState) Amalgamate(A string, B string, Graph *sync.Map, Txid uint32) {

}
