package main

import (
	"fmt"
	"log"
)

func OVValidate(s *[]SmallBankTransaction, u [][]uint16, group int, v chan map[string]AccountVersion) {

	lG := len(u[group])
	version := make(map[string]AccountVersion)
	var TxType uint8
	var From []byte
	var To []byte
	var Balance int
	for i := lG - 1; i >= 0; i-- {
		TxType = (*s)[u[group][i]].T
		From = (*s)[u[group][i]].F
		To = (*s)[u[group][i]].O
		Balance = (*s)[u[group][i]].B
		switch TxType {
		case GetBalance:
			OVGetBalance(string(From), version)
		case Amalgamate:
			OVAmalgamate(string(From), string(To), version)
		case UpdateBalance:
			OVUpdateBalance(string(From), Balance, version)
		case UpdateSaving:
			OVUpdateSaving(string(From), Balance, version)
		case SendPayment:
			OVSendPayment(string(From), string(To), Balance, version)
		case WriteCheck:
			OVWriteCheck(string(From), Balance, version)
		default:
			fmt.Println("T doesn't match")
		}
	}
	v <- version
}

func OVGetBalance(A string, version map[string]AccountVersion) {
	// don't need to modify the state of BlockchainState
}

func OVAmalgamate(A string, B string, version map[string]AccountVersion) {
	var SaveInt int
	var CheckInt int
	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}
	_, ok = version[B]
	if !ok {
		version[B] = NewAccountVersion()
	}

	Save, err := mSave[A]
	if err != true {
		log.Println(err)
	}
	SaveInt = Save

	Check, err := mCheck[B]
	if err != true {
		log.Println(err)
	}
	CheckInt = Check

	SaveInt = SaveInt + CheckInt
	//mSave[A] = SaveInt

	AVersion := version[A]
	AVersion.Save = SaveInt
	version[A] = AVersion

	//mCheck[B] = 0

	AVersion = version[B]
	AVersion.Check = 0
	version[B] = AVersion
}

func OVUpdateBalance(A string, Balance int, version map[string]AccountVersion) {
	var CheckInt int
	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}

	Check, err := mCheck[A]
	if err != true {
		log.Println(err)
	}
	CheckInt = Check

	CheckInt += Balance

	//mCheck[A] = CheckInt

	AVersion := version[A]
	AVersion.Check = CheckInt
	version[A] = AVersion
}

func OVUpdateSaving(A string, Balance int, version map[string]AccountVersion) {
	var SaveInt int
	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}

	Save, err := mSave[A]
	if err != true {
		log.Println(err)
	}
	SaveInt = Save

	SaveInt += Balance
	//mSave[A] = SaveInt

	AVersion := version[A]
	AVersion.Save = SaveInt
	version[A] = AVersion
}

func OVSendPayment(A string, B string, Balance int, version map[string]AccountVersion) {
	var CheckIntA int
	var CheckIntB int

	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}
	_, ok = version[B]
	if !ok {
		version[B] = NewAccountVersion()
	}

	CheckA, err := mCheck[A]
	if err != true {
		log.Println(err)
	}
	CheckIntA = CheckA

	CheckB, err := mCheck[B]
	if err != true {
		log.Println(err)
	}
	CheckIntB = CheckB

	CheckIntA -= Balance
	CheckIntB += Balance
	// update check value
	//mCheck[A] = CheckIntA

	//mCheck[B] = CheckIntB

	AVersion := version[A]
	AVersion.Check = CheckIntA
	version[A] = AVersion

	AVersion = version[B]
	AVersion.Check = CheckIntB
	version[B] = AVersion
}

func OVWriteCheck(A string, Balance int, version map[string]AccountVersion) {
	var SaveInt int
	var CheckInt int

	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}

	Save, err := mSave[A]
	if err != true {
		log.Println(err)
	}
	SaveInt = Save

	Check, err := mCheck[A]
	if err != true {
		log.Println(err)
	}
	CheckInt = Check

	if SaveInt+CheckInt < Balance {
		CheckInt = CheckInt - Balance - 1
	} else {
		CheckInt = CheckInt - Balance
	}
	//mCheck[A] = CheckInt

	AVersion := version[A]
	AVersion.Check = CheckInt
	version[A] = AVersion
}
