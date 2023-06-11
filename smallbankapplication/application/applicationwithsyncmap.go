package application

import (
	"log"
	"sync"
)

// GetBalance can get the balance of the account including check and save store
func (BCstate *BlockchainState) GetBalanceWithSyncMap(TxId uint16, A string, txResult chan<- TxResult, sMap *sync.Map) {
	// read only
	// lock the account
	BCstate.AccountLock[A].lock.Lock()
	AddComplexity(ByteLen, CycleNum)
	// use defer to unlock
	defer BCstate.AccountLock[A].lock.Unlock()
	// construct account dataA
	var dataA AccountData
	var result TxResult
	result = NewTxResult()

	result.AccountName = append(result.AccountName, A)
	result.CurrentTxId = TxId

	value, _ := sMap.LoadOrStore(A, NewAccountData())
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
	sMap.Store(A, dataA)
	txResult <- result
}

func (BCstate *BlockchainState) AmalgamateWithSyncMap(TxId uint16, A string, B string, txResult chan<- TxResult, sMap *sync.Map) {
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

	valueA, _ := sMap.LoadOrStore(A, NewAccountData())
	dataA = valueA.(AccountData)

	valueB, _ := sMap.LoadOrStore(B, NewAccountData())
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
	sMap.Store(A, dataA)
	sMap.Store(B, dataB)
	txResult <- result
}

func (BCstate *BlockchainState) UpdateBalanceWithSyncMap(TxId uint16, A string, Balance int, txResult chan<- TxResult, sMap *sync.Map) {
	BCstate.AccountLock[A].lock.Lock()
	AddComplexity(ByteLen, CycleNum)
	defer BCstate.AccountLock[A].lock.Unlock()

	var dataA AccountData
	var result TxResult
	result = NewTxResult()

	// only A and only change check store
	result.AccountName = append(result.AccountName, A)
	result.CurrentTxId = TxId

	valueA, _ := sMap.LoadOrStore(A, NewAccountData())
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
	sMap.Store(A, dataA)
	txResult <- result
}

func (BCstate *BlockchainState) UpdateSavingWithSyncMap(TxId uint16, A string, Balance int, txResult chan<- TxResult, sMap *sync.Map) {
	BCstate.AccountLock[A].lock.Lock()
	AddComplexity(ByteLen, CycleNum)
	defer BCstate.AccountLock[A].lock.Unlock()

	var dataA AccountData
	var result TxResult
	result = NewTxResult()

	// only A and only change check store
	result.AccountName = append(result.AccountName, A)
	result.CurrentTxId = TxId

	valueA, _ := sMap.LoadOrStore(A, NewAccountData())
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
	sMap.Store(A, dataA)
	txResult <- result
}

func (BCstate *BlockchainState) SendPaymentWithSyncMap(TxId uint16, A string, B string, Balance int, txResult chan<- TxResult, sMap *sync.Map) {
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

	valueA, _ := sMap.LoadOrStore(A, NewAccountData())
	dataA = valueA.(AccountData)

	valueB, _ := sMap.LoadOrStore(B, NewAccountData())
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

	sMap.Store(A, dataA)
	sMap.Store(B, dataB)
	txResult <- result
}

func (BCstate *BlockchainState) WriteCheckWithSyncMap(TxId uint16, A string, Balance int, txResult chan<- TxResult, sMap *sync.Map) {
	BCstate.AccountLock[A].lock.Lock()
	AddComplexity(ByteLen, CycleNum)
	defer BCstate.AccountLock[A].lock.Unlock()

	var dataA AccountData
	var result TxResult
	result = NewTxResult()

	result.AccountName = append(result.AccountName, A)
	result.CurrentTxId = TxId

	valueA, _ := sMap.LoadOrStore(A, NewAccountData())
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
	sMap.Store(A, dataA)
	txResult <- result
}
