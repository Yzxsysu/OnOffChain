package main

import (
	"log"
	"strconv"
	"strings"
	"sync"
)

type RWMap struct {
	m map[uint16]string
	sync.RWMutex
}

func NewRWMap(n int) *RWMap {
	return &RWMap{
		m: make(map[uint16]string, n),
	}
}

func (m *RWMap) Get(k uint16) (string, bool) {
	m.RLock()
	defer m.RUnlock()
	v, ok := m.m[k]
	return v, ok
}

func (m *RWMap) Set(k uint16, v string) {
	m.Lock()
	defer m.Unlock()
	m.m[k] = v
}

func (m *RWMap) Delete(k uint16) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, k)
}

func (m *RWMap) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.m)
}

func (m *RWMap) Each(f func(k uint16, v string) bool) {
	// 遍历期间有读锁
	m.RLock()
	defer m.RUnlock()

	for k, v := range m.m {
		if !f(k, v) {
			return
		}
	}
}

var (
	ByteLen  int
	CycleNum int
)

func RWOGetBalance(TxId uint16, A string, m *RWMap, version map[string]AccountVersion) {
	AddComplexity(ByteLen, CycleNum)
	// don't need to modify the state of BlockchainState
}

func RWOAmalgamate(TxId uint16, A string, B string, m *RWMap, version map[string]AccountVersion) {
	AddComplexity(ByteLen, CycleNum)
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

	v, ok := m.Get(TxId)

	if ok {
		// name="" < name=""
		s1 := strings.Split(v, "<")
		l := len(s1)
		if l == 0 {
			//log.Println("D is nil when validating graph edge")
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
	}

	if ConsistentCheckValue != "" {
		temp, _ := strconv.ParseInt(ConsistentCheckValue, 10, 64)
		CheckInt = int(temp)
	} else {
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

func RWOUpdateBalance(TxId uint16, A string, Balance int, m *RWMap, version map[string]AccountVersion) {
	AddComplexity(ByteLen, CycleNum)
	var name string
	var CheckVersion string
	var ConsistentCheckValue string

	var CheckInt int
	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}
	v, ok := m.Get(TxId)
	if ok {
		// name="" < name=""
		s1 := strings.Split(v, "<")
		l := len(s1)
		if l == 0 {
			//log.Println("D is nil when validating graph edge")
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
	}

	CheckInt += Balance

	//mCheck[A] = CheckInt

	AVersion := version[A]
	AVersion.Check = CheckInt
	version[A] = AVersion
}

// OUpdateSaving fatal error: concurrent map read and map write
func RWOUpdateSaving(TxId uint16, A string, Balance int, m *RWMap, version map[string]AccountVersion) {
	AddComplexity(ByteLen, CycleNum)
	var name string
	var SaveVersion string
	var ConsistentSaveValue string

	var SaveInt int
	// init account version
	_, ok := version[A]
	if !ok {
		version[A] = NewAccountVersion()
	}

	v, ok := m.Get(TxId)

	if ok {
		// name="" < name=""
		s1 := strings.Split(v, "<")
		l := len(s1)
		if l == 0 {
			//log.Println("D is nil when validating graph edge")
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
	}

	SaveInt += Balance
	//mSave[A] = SaveInt

	AVersion := version[A]
	AVersion.Save = SaveInt
	version[A] = AVersion
}

func RWOSendPayment(TxId uint16, A string, B string, Balance int, m *RWMap, version map[string]AccountVersion) {
	AddComplexity(ByteLen, CycleNum)
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

	v, ok := m.Get(TxId)
	if ok {
		// name="" < name=""
		s1 := strings.Split(v, "<")
		l := len(s1)
		if l == 0 {
			//log.Println("D is nil when validating graph edge")
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
		//CheckA, err := mCheck[A]
		if err != true {
			log.Println(err)
		}
		CheckA, ok := CheckValue.(int)
		if ok == false {
			CheckA = 0
			//log.Println(A, "syncMap panic:", "panic: interface conversion: interface {} is nil, not int")
		}
		CheckIntA = CheckA
	}

	if CheckBS != "" {
		temp, _ := strconv.ParseInt(CheckBS, 10, 64)
		CheckIntB = int(temp)
	} else {
		CheckValue, err := syncCheck.Load(B)
		//CheckB, err := mCheck[B]
		if err != true {
			log.Println(err)
		}
		CheckB, ok := CheckValue.(int)
		if ok == false {
			CheckB = 0
			log.Println(B, "syncMap panic:", "panic: interface conversion: interface {} is nil, not int")
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

func RWOWriteCheck(TxId uint16, A string, Balance int, m *RWMap, version map[string]AccountVersion) {
	AddComplexity(ByteLen, CycleNum)
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

	v, ok := m.Get(TxId)
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
	}

	if ConsistentCheckValue != "" {
		temp, _ := strconv.ParseInt(ConsistentCheckValue, 10, 64)
		CheckInt = int(temp)
	} else {
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
