package main

import (
	"fmt"
	"log"
)

func OVValidate(s *[]SmallBankTransaction, u *[][]uint16, group int, v chan map[string]AccountVersion) {
	//log.Println("Before OVValidate", group)
	if len(*u) == 0 {
		//log.Println("OVValidate", "v <- version", "len(*u) == 0", group)
		version := make(map[string]AccountVersion)
		v <- version
		return
	}
	lG := len((*u)[group])
	version := make(map[string]AccountVersion)
	var TxType uint8
	var From []byte
	var To []byte
	var Balance int
	for i := lG - 1; i >= 0; i-- {
		TxType = (*s)[(*u)[group][i]-1].T
		From = (*s)[(*u)[group][i]-1].F
		To = (*s)[(*u)[group][i]-1].O
		Balance = (*s)[(*u)[group][i]-1].B
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

	//log.Println("After OVValidate v <- version", group)
	v <- version
}

func OVGetBalance(A string, version map[string]AccountVersion) {
	AddComplexity(ByteLen, CycleNum)
	// don't need to modify the state of BlockchainState
}

func OVAmalgamate(A string, B string, version map[string]AccountVersion) {
	AddComplexity(ByteLen, CycleNum)
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

	SaveValue, err := syncSave.Load(A)
	//Save, err := mSave[A]
	if err != true {
		log.Println(err)
	}
	Save, ok := SaveValue.(int)
	if ok == false {
		Save = 0
		//log.Println(A, "syncMap panic:", "panic: interface conversion: interface {} is nil, not int")
	}
	SaveInt = Save

	CheckValue, err := syncCheck.Load(B)
	//Check, err := mCheck[B]
	if err != true {
		log.Println(err)
	}
	Check, ok := CheckValue.(int)
	if ok == false {
		Check = 0
		//log.Println(B, "syncMap panic:", "panic: interface conversion: interface {} is nil, not int")
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
	AddComplexity(ByteLen, CycleNum)
	var CheckInt int
	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}

	CheckValue, err := syncCheck.Load(A)
	//Check, err := mCheck[A]
	if err != true {
		log.Println(err)
	}
	Check, ok := CheckValue.(int)
	if ok == false {
		Check = 0
		//log.Println(A, "syncMap panic:", "panic: interface conversion: interface {} is nil, not int")
	}
	CheckInt = Check

	CheckInt += Balance

	//mCheck[A] = CheckInt

	AVersion := version[A]
	AVersion.Check = CheckInt
	version[A] = AVersion
}

func OVUpdateSaving(A string, Balance int, version map[string]AccountVersion) {
	AddComplexity(ByteLen, CycleNum)
	var SaveInt int
	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}

	SaveValue, err := syncSave.Load(A)
	//Save, err := mSave[A]
	if err != true {
		log.Println(err)
	}
	Save, ok := SaveValue.(int)
	if ok == false {
		Save = 0
		//log.Println(A, "syncMap panic:", "panic: interface conversion: interface {} is nil, not int")
	}
	SaveInt = Save

	SaveInt += Balance
	//mSave[A] = SaveInt

	AVersion := version[A]
	AVersion.Save = SaveInt
	version[A] = AVersion
}

func OVSendPayment(A string, B string, Balance int, version map[string]AccountVersion) {
	AddComplexity(ByteLen, CycleNum)
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

	CheckValueA, err := syncCheck.Load(A)
	//CheckA, err := mCheck[A]
	if err != true {
		log.Println(err)
	}
	CheckA, ok := CheckValueA.(int)
	if ok == false {
		CheckA = 0
		//log.Println(A, "syncMap panic:", "panic: interface conversion: interface {} is nil, not int")
	}
	CheckIntA = CheckA

	CheckValueB, err := syncCheck.Load(B)
	//CheckB, err := mCheck[B]
	if err != true {
		log.Println(err)
	}
	CheckB, ok := CheckValueB.(int)
	if ok == false {
		CheckB = 0
		//log.Println(B, "syncMap panic:", "panic: interface conversion: interface {} is nil, not int")
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
	AddComplexity(ByteLen, CycleNum)
	var SaveInt int
	var CheckInt int

	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}

	SaveValue, err := syncSave.Load(A)
	//Save, err := mSave[A]
	if err != true {
		log.Println(err)
	}
	Save, ok := SaveValue.(int)
	if ok == false {
		Save = 0
		log.Println(A, "syncMap panic:", "panic: interface conversion: interface {} is nil, not int")
	}
	SaveInt = Save

	CheckValue, err := syncCheck.Load(A)
	//Check, err := mCheck[A]
	if err != true {
		log.Println(err)
	}
	Check, ok := CheckValue.(int)
	if ok == false {
		Check = 0
		log.Println(A, "syncMap panic:", "panic: interface conversion: interface {} is nil, not int")
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
