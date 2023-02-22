package application

import (
	"fmt"
	"log"
)

func (BCstate *BlockchainState) VValidate(s *[]SmallBankTransaction, v *[]uint16, ch chan bool) {
	log.Println("VValidate:", *v)
	l := len(*v)
	if l == 0 {
		return
	}
	var TxType uint8
	//var I uint16
	var From []byte
	var To []byte
	var Balance int

	for i := 0; i < l; i++ {
		tx := (*s)[(*v)[i]-1]
		TxType = tx.T
		//I = tx.I
		From = tx.F
		To = tx.O
		Balance = tx.B
		switch TxType {
		case GetBalance:
			go BCstate.VGetBalance(string(From))
		case Amalgamate:
			go BCstate.VAmalgamate(string(From), string(To))
		case UpdateBalance:
			go BCstate.VUpdateBalance(string(From), Balance)
		case UpdateSaving:
			go BCstate.VUpdateSaving(string(From), Balance)
		case SendPayment:
			go BCstate.VSendPayment(string(From), string(To), Balance)
		case WriteCheck:
			go BCstate.VWriteCheck(string(From), Balance)
		default:
			fmt.Println("T doesn't match")
		}
	}
	ch <- true
}

func (BCstate *BlockchainState) VGetBalance(A string) {
	// don't need to modify the state of BlockchainState
}

func (BCstate *BlockchainState) VAmalgamate(A string, B string) {
	Save, err := BCstate.SavingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	SaveInt := BytesToInt(Save)

	Check, err := BCstate.CheckingStore.Get([]byte(B))
	if err != nil {
		log.Println(err)
	}
	CheckInt := BytesToInt(Check)

	SaveInt = SaveInt + CheckInt
	err = BCstate.SavingStore.Set([]byte(A), IntToBytes(SaveInt))
	if err != nil {
		log.Println(err)
	}

	err = BCstate.CheckingStore.Set([]byte(B), IntToBytes(0))
	if err != nil {
		log.Println(err)
	}
}

func (BCstate *BlockchainState) VUpdateBalance(A string, Balance int) {
	Check, err := BCstate.CheckingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	CheckInt := BytesToInt(Check)
	CheckInt += Balance

	err = BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckInt))
	if err != nil {
		log.Println(err)
	}
}

func (BCstate *BlockchainState) VUpdateSaving(A string, Balance int) {
	Save, err := BCstate.SavingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	SaveInt := BytesToInt(Save)
	SaveInt += Balance

	err = BCstate.SavingStore.Set([]byte(A), IntToBytes(SaveInt))
	if err != nil {
		log.Println(err)
	}
}

func (BCstate *BlockchainState) VSendPayment(A string, B string, Balance int) {
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
}

func (BCstate *BlockchainState) VWriteCheck(A string, Balance int) {
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

	if SaveInt+CheckInt < Balance {
		CheckInt = CheckInt - Balance - 1
	} else {
		CheckInt = CheckInt - Balance
	}
	err = BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckInt))
	if err != nil {
		log.Println(err)
	}
}
