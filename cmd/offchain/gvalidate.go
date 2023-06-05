package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func OValidate(s *[]SmallBankTransaction, GE *[][]GraphEdge, group int, v chan map[string]AccountVersion) {
	order, m := Dfs(*GE, group)
	RWm := NewRWMap(len(m))
	RWm.m = m

	log.Println("OValidate:group", order, group)

	lG := len(order)
	version := make(map[string]AccountVersion)
	var TxType uint8
	var TxId uint16
	var From []byte
	var To []byte
	var Balance int
	for i := lG - 1; i >= 0; i-- {
		TxType = (*s)[order[i]-1].T
		TxId = (*s)[order[i]-1].I
		From = (*s)[order[i]-1].F
		To = (*s)[order[i]-1].O
		Balance = (*s)[order[i]-1].B
		switch TxType {
		case GetBalance:
			RWOGetBalance(TxId, string(From), RWm, version)
			//OGetBalance(TxId, string(From), m, version)
		case Amalgamate:
			RWOAmalgamate(TxId, string(From), string(To), RWm, version)
			//OAmalgamate(TxId, string(From), string(To), m, version)
		case UpdateBalance:
			RWOUpdateBalance(TxId, string(From), Balance, RWm, version)
			//OUpdateBalance(TxId, string(From), Balance, m, version)
		case UpdateSaving:
			RWOUpdateSaving(TxId, string(From), Balance, RWm, version)
			//OUpdateSaving(TxId, string(From), Balance, m, version)
		case SendPayment:
			RWOSendPayment(TxId, string(From), string(To), Balance, RWm, version)
			//OSendPayment(TxId, string(From), string(To), Balance, m, version)
		case WriteCheck:
			RWOWriteCheck(TxId, string(From), Balance, RWm, version)
			//OWriteCheck(TxId, string(From), Balance, m, version)
		default:
			fmt.Println("T doesn't match")
		}
	}

	log.Println("v <- version", version, group)
	v <- version
}

func OGetBalance(TxId uint16, A string, m map[uint16]string, version map[string]AccountVersion) {
	// don't need to modify the state of BlockchainState
}

func OAmalgamate(TxId uint16, A string, B string, m map[uint16]string, version map[string]AccountVersion) {
	var SaveInt int
	var CheckInt int

	var name string
	var SaveVersion string
	var ConsistentSaveValue string
	var CheckVersion string
	var ConsistentCheckValue string

	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}
	_, ok = version[B]
	if !ok {
		version[B] = NewAccountVersion()
	}

	v, ok := m[TxId]
	if ok {
		// name="" < name=""
		s1 := strings.Split(v, "<")
		l := len(s1)
		if l == 0 {
			log.Println("D is nil when validating graph edge")
		}
		// "name="1",SaveVersion=10,ConsistentSaveValue=2"
		for _, elements := range s1 {
			element := strings.Split(elements, ",")
			for _, e := range element {
				kv := strings.Split(e, "=")
				switch kv[0] {
				case "name":
					name = kv[1]
				case "SaveVersion":
					SaveVersion = kv[1]
				case "ConsistentSaveValue":
					ConsistentSaveValue = kv[1]
				case "CheckVersion":
					CheckVersion = kv[1]
				case "ConsistentCheckValue":
					ConsistentCheckValue = kv[1]
				}
			}
			// modify the account version value
			if name == A {
				AVersion := version[A]
				temp, _ := strconv.ParseInt(SaveVersion, 10, 64)
				AVersion.SaveVersion = uint16(temp)
				version[A] = AVersion
			}
			if name == B {
				AVersion := version[B]
				temp, _ := strconv.ParseInt(CheckVersion, 10, 64)
				AVersion.CheckVersion = uint16(temp)
				version[B] = AVersion
			}
		}
	}

	if ConsistentSaveValue != "" {
		temp, _ := strconv.ParseInt(ConsistentSaveValue, 10, 64)
		SaveInt = int(temp)
	} else {
		SaveValue, err := syncSave.Load(A)
		Save := SaveValue.(int)
		//Save, err := mSave[A]
		if err != true {
			log.Println(err)
		}
		SaveInt = Save
	}

	if ConsistentCheckValue != "" {
		temp, _ := strconv.ParseInt(ConsistentCheckValue, 10, 64)
		CheckInt = int(temp)
	} else {
		CheckValue, err := syncCheck.Load(B)
		Check := CheckValue.(int)
		//Check, err := mCheck[B]
		if err != true {
			log.Println(err)
		}
		CheckInt = Check
	}

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

func OUpdateBalance(TxId uint16, A string, Balance int, m map[uint16]string, version map[string]AccountVersion) {
	var name string
	var CheckVersion string
	var ConsistentCheckValue string

	var CheckInt int
	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}
	v, ok := m[TxId]
	if ok {
		// name="" < name=""
		s1 := strings.Split(v, "<")
		l := len(s1)
		if l == 0 {
			log.Println("D is nil when validating graph edge")
		}
		// "name="1",CheckVersion=10,ConsistentCheckValue=2"
		for _, elements := range s1 {
			element := strings.Split(elements, ",")
			for _, e := range element {
				kv := strings.Split(e, "=")
				switch kv[0] {
				case "name":
					name = kv[1]
				case "CheckVersion":
					CheckVersion = kv[1]
				case "ConsistentCheckValue":
					ConsistentCheckValue = kv[1]
				}
			}
			// modify the account version value
			if name == A {
				AVersion := version[A]
				temp, _ := strconv.ParseInt(CheckVersion, 10, 64)
				AVersion.CheckVersion = uint16(temp)
				version[A] = AVersion
			}
		}
	}

	if ConsistentCheckValue != "" {
		temp, _ := strconv.ParseInt(ConsistentCheckValue, 10, 64)
		CheckInt = int(temp)
	} else {
		CheckValue, err := syncCheck.Load(A)
		Check := CheckValue.(int)
		//Check, err := mCheck[A]
		if err != true {
			log.Println(err)
		}
		CheckInt = Check
	}

	CheckInt += Balance

	//mCheck[A] = CheckInt

	AVersion := version[A]
	AVersion.Check = CheckInt
	version[A] = AVersion
}

// OUpdateSaving fatal error: concurrent map read and map write
func OUpdateSaving(TxId uint16, A string, Balance int, m map[uint16]string, version map[string]AccountVersion) {
	var name string
	var SaveVersion string
	var ConsistentSaveValue string

	var SaveInt int
	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}
	v, ok := m[TxId]
	if ok {
		// name="" < name=""
		s1 := strings.Split(v, "<")
		l := len(s1)
		if l == 0 {
			log.Println("D is nil when validating graph edge")
		}
		// "name="1",SaveVersion=10,ConsistentSaveValue=2"
		for _, elements := range s1 {
			element := strings.Split(elements, ",")
			for _, e := range element {
				kv := strings.Split(e, "=")
				switch kv[0] {
				case "name":
					name = kv[1]
				case "SaveVersion":
					SaveVersion = kv[1]
				case "ConsistentSaveValue":
					ConsistentSaveValue = kv[1]
				}
			}
			// modify the account version value
			if name == A {
				AVersion := version[A]
				temp, _ := strconv.ParseInt(SaveVersion, 10, 64)
				AVersion.SaveVersion = uint16(temp)
				version[A] = AVersion
			}
		}
	}
	if ConsistentSaveValue != "" {
		temp, _ := strconv.ParseInt(ConsistentSaveValue, 10, 64)
		SaveInt = int(temp)
	} else {
		SaveValue, err := syncSave.Load(A)
		Save := SaveValue.(int)
		//Save, err := mSave[A]
		if err != true {
			log.Println(err)
		}
		SaveInt = Save
	}

	SaveInt += Balance
	//mSave[A] = SaveInt

	AVersion := version[A]
	AVersion.Save = SaveInt
	version[A] = AVersion
}

func OSendPayment(TxId uint16, A string, B string, Balance int, m map[uint16]string, version map[string]AccountVersion) {
	var name string
	var CheckVersion string
	var ConsistentCheckValue string

	var CheckAS string
	var CheckBS string

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

	v, ok := m[TxId]
	if ok {
		// name="" < name=""
		s1 := strings.Split(v, "<")
		l := len(s1)
		if l == 0 {
			log.Println("D is nil when validating graph edge")
		}
		// "name="1",CheckVersion=10,ConsistentCheckValue=2"
		for _, elements := range s1 {
			element := strings.Split(elements, ",")
			for _, e := range element {
				kv := strings.Split(e, "=")
				switch kv[0] {
				case "name":
					name = kv[1]
				case "CheckVersion":
					CheckVersion = kv[1]
				case "ConsistentCheckValue":
					ConsistentCheckValue = kv[1]
				}
			}
			// modify the account version value
			if name == A {
				AVersion := version[A]
				temp, _ := strconv.ParseInt(CheckVersion, 10, 64)
				AVersion.CheckVersion = uint16(temp)
				version[A] = AVersion
				CheckAS = ConsistentCheckValue
			}
			if name == B {
				AVersion := version[B]
				temp, _ := strconv.ParseInt(CheckVersion, 10, 64)
				AVersion.CheckVersion = uint16(temp)
				version[B] = AVersion
				CheckBS = ConsistentCheckValue
			}
		}
	}

	if CheckAS != "" {
		temp, _ := strconv.ParseInt(CheckAS, 10, 64)
		CheckIntA = int(temp)
	} else {
		CheckValue, err := syncCheck.Load(A)
		CheckA := CheckValue.(int)
		//CheckA, err := mCheck[A]
		if err != true {
			log.Println(err)
		}
		CheckIntA = CheckA
	}

	if CheckBS != "" {
		temp, _ := strconv.ParseInt(CheckBS, 10, 64)
		CheckIntB = int(temp)
	} else {
		CheckValue, err := syncCheck.Load(B)
		CheckB := CheckValue.(int)
		//CheckB, err := mCheck[B]
		if err != true {
			log.Println(err)
		}
		CheckIntB = CheckB
	}

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

func OWriteCheck(TxId uint16, A string, Balance int, m map[uint16]string, version map[string]AccountVersion) {
	var SaveInt int
	var CheckInt int

	var name string
	var SaveVersion string
	var ConsistentSaveValue string
	var CheckVersion string
	var ConsistentCheckValue string

	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}

	v, ok := m[TxId]
	if ok {
		// name="" < name=""
		s1 := strings.Split(v, "<")
		l := len(s1)
		if l == 0 {
			log.Println("D is nil when validating graph edge")
		}
		// "name="1",SaveVersion=10,ConsistentSaveValue=2"
		for _, elements := range s1 {
			element := strings.Split(elements, ",")
			for _, e := range element {
				kv := strings.Split(e, "=")
				switch kv[0] {
				case "name":
					name = kv[1]
				case "SaveVersion":
					SaveVersion = kv[1]
				case "ConsistentSaveValue":
					ConsistentSaveValue = kv[1]
				case "CheckVersion":
					CheckVersion = kv[1]
				case "ConsistentCheckValue":
					ConsistentCheckValue = kv[1]
				}
			}
			// modify the account version value
			if name == A {
				AVersion := version[A]
				if SaveVersion != "" {
					temp, _ := strconv.ParseInt(SaveVersion, 10, 64)
					AVersion.SaveVersion = uint16(temp)
				}
				if CheckVersion != "" {
					temp, _ := strconv.ParseInt(CheckVersion, 10, 64)
					AVersion.CheckVersion = uint16(temp)
				}
				version[A] = AVersion
			}
		}
	}
	if ConsistentSaveValue != "" {
		temp, _ := strconv.ParseInt(ConsistentSaveValue, 10, 64)
		SaveInt = int(temp)
	} else {
		SaveValue, err := syncSave.Load(A)
		Save := SaveValue.(int)
		//Save, err := mSave[A]
		if err != true {
			log.Println(err)
		}
		SaveInt = Save
	}

	if ConsistentCheckValue != "" {
		temp, _ := strconv.ParseInt(ConsistentCheckValue, 10, 64)
		CheckInt = int(temp)
	} else {
		CheckValue, err := syncCheck.Load(A)
		Check := CheckValue.(int)
		//Check, err := mCheck[A]
		if err != true {
			log.Println(err)
		}
		CheckInt = Check
	}

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
