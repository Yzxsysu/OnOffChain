package application

import (
	dbm "github.com/Yzxsysu/onoffchain/src/tm-db"
	"log"
	"sync"
)

type BackendType string

var (
	ByteLen  int
	CycleNum int
)

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
	AccountLock map[string]*Lock
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

// NewBlockchainState can choose {goleveldb, cleveldb, memdb, boltdb, rocksdb, badgerdb}
// input Corresponding name
func NewBlockchainState(DBName BackendType, leader bool, dir string) (*BlockchainState, error, error) {
	var BaseCaseState BlockchainState
	var err error
	var err1 error
	// choose DB

	switch DBName {
	case GoLevelDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.GoLevelDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.GoLevelDBBackend, dir+"save")
	case CLevelDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.CLevelDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.CLevelDBBackend, dir+"save")
	case MemDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.MemDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.MemDBBackend, dir+"save")
	case BoltDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.BoltDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.BoltDBBackend, dir+"save")
	case RocksDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.RocksDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.RocksDBBackend, dir+"save")
	case BadgerDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.BadgerDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.BadgerDBBackend, dir+"save")
	}
	if err != nil || err1 != nil {
		log.Fatalf("Create db error: %v, %v", err, err1)
	}
	BaseCaseState.Leader = leader
	return &BlockchainState{
		CheckingStore: BaseCaseState.CheckingStore,
		SavingStore:   BaseCaseState.SavingStore,
		Height:        1,
		Leader:        leader,
		AccountLock:   make(map[string]*Lock),
	}, err, err1
}

// CreateAccount can create account with saving balance and checking balance
func (BCstate *BlockchainState) CreateAccount(AccountName string, SavingBalance int, CheckingBalance int) {
	// Create two separate accounts for two DB
	var err error
	err = BCstate.SavingStore.Set([]byte(AccountName), IntToBytes(SavingBalance))
	if err != nil {
		log.Println(err)
	}
	err = BCstate.CheckingStore.Set([]byte(AccountName), IntToBytes(CheckingBalance))
	if err != nil {
		log.Println(err)
	}
	// Create two separate locks for each account
	// Initialize the lock in BCstate
	BCstate.AccountLock[AccountName] = NewLock()
}

// GetBalance can get the balance of the account including check and save store
func (BCstate *BlockchainState) GetBalance(TxId uint16, A string, txResult chan<- TxResult) {
	// read only
	// lock the account
	BCstate.AccountLock[A].lock.RLock()
	AddComplexity(ByteLen, CycleNum)
	// use defer to unlock
	defer BCstate.AccountLock[A].lock.RUnlock()
	// construct account dataA
	var dataA AccountData
	var result TxResult
	result = NewTxResult()

	result.AccountName = append(result.AccountName, A)
	result.CurrentTxId = TxId

	value, _ := AccountDataMap.LoadOrStore(A, NewAccountData())
	// type assertion AccountData
	dataA = value.(AccountData)

	// GetBalance only read one account's check and save value
	result.SaveBool = append(result.SaveBool, true)
	result.CheckBool = append(result.CheckBool, true)

	result.PreTxId = append(result.PreTxId, dataA.WrittenBy) // 0 or pre witten I

	result.SaveVersion = append(result.SaveVersion, dataA.SaveVersion)    // 0 or next version
	result.CheckVersion = append(result.CheckVersion, dataA.CheckVersion) // 0 or next version

	// only read don't need to add 1
	dataA.SaveVersion = result.SaveVersion[0]
	dataA.CheckVersion = result.CheckVersion[0]

	Save, err := BCstate.SavingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	SaveInt := BytesToInt(Save)

	Check, err := BCstate.CheckingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	CheckInt := BytesToInt(Check)

	result.ConsistentSaveValue = append(result.ConsistentSaveValue, SaveInt)
	result.ConsistentCheckValue = append(result.ConsistentCheckValue, CheckInt)
	// add new AccountData
	dataA.WrittenBy = TxId

	// store the account name
	AccountDataMap.Store(A, dataA)
	txResult <- result
}

func (BCstate *BlockchainState) Amalgamate(TxId uint16, A string, B string, txResult chan<- TxResult) {
	BCstate.AccountLock[A].lock.Lock()
	BCstate.AccountLock[B].lock.Lock()
	AddComplexity(ByteLen, CycleNum)
	defer BCstate.AccountLock[A].lock.Unlock()
	defer BCstate.AccountLock[B].lock.Unlock()
	var dataA AccountData
	var dataB AccountData
	var result TxResult
	result = NewTxResult()

	// A first, B next
	result.AccountName = append(result.AccountName, A, B)
	result.CurrentTxId = TxId

	valueA, _ := AccountDataMap.LoadOrStore(A, NewAccountData())
	dataA = valueA.(AccountData)

	valueB, _ := AccountDataMap.LoadOrStore(B, NewAccountData())
	dataB = valueB.(AccountData)

	result.SaveBool = append(result.SaveBool, true, false)
	result.CheckBool = append(result.CheckBool, false, true)
	result.PreTxId = append(result.PreTxId, dataA.WrittenBy, dataB.WrittenBy)                 // 0 or pre witten I
	result.SaveVersion = append(result.SaveVersion, dataA.SaveVersion, dataB.SaveVersion)     // 0 or next version
	result.CheckVersion = append(result.CheckVersion, dataA.CheckVersion, dataB.CheckVersion) // 0 or next version

	dataA.SaveVersion = result.SaveVersion[0] + 1
	dataA.CheckVersion = result.CheckVersion[0]

	dataB.SaveVersion = result.SaveVersion[1]
	dataB.CheckVersion = result.CheckVersion[1] + 1

	Save, err := BCstate.SavingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	SaveInt := BytesToInt(Save)
	result.ConsistentSaveValue = append(result.ConsistentSaveValue, SaveInt, 0)

	Check, err := BCstate.CheckingStore.Get([]byte(B))
	if err != nil {
		log.Println(err)
	}
	CheckInt := BytesToInt(Check)
	result.ConsistentCheckValue = append(result.ConsistentCheckValue, 0, CheckInt)

	SaveInt = SaveInt + CheckInt
	err = BCstate.SavingStore.Set([]byte(A), IntToBytes(SaveInt))
	if err != nil {
		log.Println(err)
	}

	err = BCstate.CheckingStore.Set([]byte(B), IntToBytes(0))
	if err != nil {
		log.Println(err)
	}

	dataA.WrittenBy = TxId
	dataB.WrittenBy = TxId
	AccountDataMap.Store(A, dataA)
	AccountDataMap.Store(B, dataB)
	txResult <- result
}

func (BCstate *BlockchainState) UpdateBalance(TxId uint16, A string, Balance int, txResult chan<- TxResult) {
	BCstate.AccountLock[A].lock.Lock()
	AddComplexity(ByteLen, CycleNum)
	defer BCstate.AccountLock[A].lock.Unlock()

	var dataA AccountData
	var result TxResult
	result = NewTxResult()

	// only A and only change check store
	result.AccountName = append(result.AccountName, A)
	result.CurrentTxId = TxId

	valueA, _ := AccountDataMap.LoadOrStore(A, NewAccountData())
	dataA = valueA.(AccountData)

	result.SaveBool = append(result.SaveBool, false)
	result.CheckBool = append(result.CheckBool, true)
	result.PreTxId = append(result.PreTxId, dataA.WrittenBy)

	result.SaveVersion = append(result.SaveVersion, dataA.SaveVersion)
	result.CheckVersion = append(result.CheckVersion, dataA.CheckVersion)

	// SaveVersion can not be nil, otherwise it will encounter error
	dataA.SaveVersion = result.SaveVersion[0]
	dataA.CheckVersion = result.CheckVersion[0] + 1

	Check, err := BCstate.CheckingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	CheckInt := BytesToInt(Check)
	result.ConsistentCheckValue = append(result.ConsistentCheckValue, CheckInt)
	CheckInt += Balance

	err = BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckInt))
	if err != nil {
		log.Println(err)
	}
	dataA.WrittenBy = TxId
	AccountDataMap.Store(A, dataA)
	txResult <- result
}

func (BCstate *BlockchainState) UpdateSaving(TxId uint16, A string, Balance int, txResult chan<- TxResult) {
	BCstate.AccountLock[A].lock.Lock()
	AddComplexity(ByteLen, CycleNum)
	defer BCstate.AccountLock[A].lock.Unlock()

	var dataA AccountData
	var result TxResult
	result = NewTxResult()

	// only A and only change check store
	result.AccountName = append(result.AccountName, A)
	result.CurrentTxId = TxId

	valueA, _ := AccountDataMap.LoadOrStore(A, NewAccountData())
	dataA = valueA.(AccountData)

	result.SaveBool = append(result.SaveBool, true)
	result.CheckBool = append(result.CheckBool, false)
	result.PreTxId = append(result.PreTxId, dataA.WrittenBy)

	result.SaveVersion = append(result.SaveVersion, dataA.SaveVersion)
	result.CheckVersion = append(result.CheckVersion, dataA.CheckVersion)

	dataA.SaveVersion = result.SaveVersion[0] + 1
	dataA.CheckVersion = result.CheckVersion[0]

	Save, err := BCstate.SavingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	SaveInt := BytesToInt(Save)
	result.ConsistentSaveValue = append(result.ConsistentSaveValue, SaveInt)
	SaveInt += Balance

	err = BCstate.CheckingStore.Set([]byte(A), IntToBytes(SaveInt))
	if err != nil {
		log.Println(err)
	}
	dataA.WrittenBy = TxId
	AccountDataMap.Store(A, dataA)
	txResult <- result
}

func (BCstate *BlockchainState) SendPayment(TxId uint16, A string, B string, Balance int, txResult chan<- TxResult) {
	BCstate.AccountLock[A].lock.Lock()
	BCstate.AccountLock[B].lock.Lock()
	AddComplexity(ByteLen, CycleNum)
	defer BCstate.AccountLock[A].lock.Unlock()
	defer BCstate.AccountLock[B].lock.Unlock()

	var dataA AccountData
	var dataB AccountData
	var result TxResult
	result = NewTxResult()

	// A first, B next
	result.AccountName = append(result.AccountName, A, B)
	result.CurrentTxId = TxId

	valueA, _ := AccountDataMap.LoadOrStore(A, NewAccountData())
	dataA = valueA.(AccountData)

	valueB, _ := AccountDataMap.LoadOrStore(B, NewAccountData())
	dataB = valueB.(AccountData)

	result.SaveBool = append(result.SaveBool, false, false)
	result.CheckBool = append(result.CheckBool, true, true)
	result.PreTxId = append(result.PreTxId, dataA.WrittenBy, dataB.WrittenBy)
	result.SaveVersion = append(result.SaveVersion, dataA.SaveVersion, dataB.SaveVersion)     // 0 or next version
	result.CheckVersion = append(result.CheckVersion, dataA.CheckVersion, dataB.CheckVersion) // 0 or next version

	dataA.SaveVersion = result.SaveVersion[0]
	dataA.CheckVersion = result.CheckVersion[0] + 1

	dataB.SaveVersion = result.SaveVersion[1]
	dataB.CheckVersion = result.CheckVersion[1] + 1

	CheckA, err := BCstate.CheckingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}

	CheckIntA := BytesToInt(CheckA)

	CheckB, err := BCstate.CheckingStore.Get([]byte(B))
	if err != nil {
		log.Println(err)
	}
	CheckIntB := BytesToInt(CheckB)

	result.ConsistentCheckValue = append(result.ConsistentCheckValue, CheckIntA, CheckIntB)

	CheckIntA -= Balance
	CheckIntB += Balance

	// update check value
	err = BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckIntA))
	if err != nil {
		log.Println(err)
	}

	err = BCstate.CheckingStore.Set([]byte(B), IntToBytes(CheckIntB))
	if err != nil {
		log.Println(err)
	}

	dataA.WrittenBy = TxId
	dataB.WrittenBy = TxId

	AccountDataMap.Store(A, dataA)
	AccountDataMap.Store(B, dataB)
	txResult <- result
}

func (BCstate *BlockchainState) WriteCheck(TxId uint16, A string, Balance int, txResult chan<- TxResult) {
	BCstate.AccountLock[A].lock.Lock()
	AddComplexity(ByteLen, CycleNum)
	defer BCstate.AccountLock[A].lock.Unlock()

	var dataA AccountData
	var result TxResult
	result = NewTxResult()

	result.AccountName = append(result.AccountName, A)
	result.CurrentTxId = TxId

	valueA, _ := AccountDataMap.LoadOrStore(A, NewAccountData())
	dataA = valueA.(AccountData)

	// need to be added in the edge when creating
	result.SaveBool = append(result.SaveBool, true)
	result.CheckBool = append(result.CheckBool, true)
	result.PreTxId = append(result.PreTxId, dataA.WrittenBy)

	result.SaveVersion = append(result.SaveVersion, dataA.SaveVersion)
	result.CheckVersion = append(result.CheckVersion, dataA.CheckVersion)

	// save account is not be written
	// only read so the version doesn't need to add 1
	dataA.SaveVersion = result.SaveVersion[0]
	dataA.CheckVersion = result.CheckVersion[0] + 1

	Save, err := BCstate.SavingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	SaveInt := BytesToInt(Save)
	result.ConsistentSaveValue = append(result.ConsistentSaveValue, SaveInt)

	Check, err := BCstate.CheckingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	CheckInt := BytesToInt(Check)
	result.ConsistentCheckValue = append(result.ConsistentCheckValue, CheckInt)

	if SaveInt+CheckInt < Balance {
		CheckInt = CheckInt - Balance - 1
	} else {
		CheckInt = CheckInt - Balance
	}
	err = BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckInt))
	if err != nil {
		log.Println(err)
	}
	dataA.WrittenBy = TxId
	AccountDataMap.Store(A, dataA)
	txResult <- result
}
