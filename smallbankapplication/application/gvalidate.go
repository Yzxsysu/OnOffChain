package application

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func Dfs(GE [][]GraphEdge, groupNum int) ([]uint16, map[uint16]string) {
	//log.Println("Dfs:", groupNum)
	if len(GE) == 0 {
		return make([]uint16, 0), make(map[uint16]string)
	}
	sub := GE[groupNum]
	l := len(sub)
	m := make(map[uint16]string)
	s := NewSortedGraph(uint16(l))
	for i := 0; i < l; i++ {
		from := sub[i].F
		to := sub[i].T
		s.AddEdge(from, to)
		_, ok := m[to]
		if ok {
			m[to] += ">" + sub[i].D
		} else {
			m[to] = sub[i].D
		}
	}
	order := s.TopoSortByDFS()
	return order, m
}

func (BCstate *BlockchainState) GValidate(s []SmallBankTransaction, GE *[][]GraphEdge, group int, v chan map[string]AccountVersion, ch chan bool) {
	order, m := Dfs(*GE, group)
	lG := len(order)
	//log.Println("GValidate:", group)
	if lG == 0 {
		//.Println("lG == 0")
		version := make(map[string]AccountVersion)
		v <- version
		ch <- true
		return
	}
	version := make(map[string]AccountVersion)
	var TxType uint8
	var TxId uint16
	var From []byte
	var To []byte
	var Balance int
	for i := lG - 1; i >= 0; i-- {
		TxType = (s)[order[i]-1].T
		TxId = (s)[order[i]-1].I
		From = (s)[order[i]-1].F
		To = (s)[order[i]-1].O
		Balance = (s)[order[i]-1].B
		switch TxType {
		case GetBalance:
			BCstate.GGetBalance(TxId, string(From), m, version)
		case Amalgamate:
			BCstate.GAmalgamate(TxId, string(From), string(To), m, version)
		case UpdateBalance:
			BCstate.GUpdateBalance(TxId, string(From), Balance, m, version)
		case UpdateSaving:
			BCstate.GUpdateSaving(TxId, string(From), Balance, m, version)
		case SendPayment:
			BCstate.GSendPayment(TxId, string(From), string(To), Balance, m, version)
		case WriteCheck:
			BCstate.GWriteCheck(TxId, string(From), Balance, m, version)
		default:
			fmt.Println("T doesn't match")
		}
	}
	v <- version
	ch <- true
}

func (BCstate *BlockchainState) GGetBalance(TxId uint16, A string, m map[uint16]string, version map[string]AccountVersion) {
	// don't need to modify the state of BlockchainState
	AddComplexity(ByteLen, CycleNum)
}

func (BCstate *BlockchainState) GAmalgamate(TxId uint16, A string, B string, m map[uint16]string, version map[string]AccountVersion) {
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

	v, ok := m[TxId]
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
		Save, err := BCstate.SavingStore.Get([]byte(A))
		if err != nil {
			log.Println(err)
		}
		SaveInt = BytesToInt(Save)
	}

	if ConsistentCheckValue != "" {
		temp, _ := strconv.ParseInt(ConsistentCheckValue, 10, 64)
		CheckInt = int(temp)
	} else {
		Check, err := BCstate.CheckingStore.Get([]byte(B))
		if err != nil {
			log.Println(err)
		}
		CheckInt = BytesToInt(Check)
	}

	SaveInt = SaveInt + CheckInt
	err := BCstate.SavingStore.Set([]byte(A), IntToBytes(SaveInt))
	if err != nil {
		log.Println(err)
	}
	AVersion := version[A]
	AVersion.Save = SaveInt
	version[A] = AVersion

	err = BCstate.CheckingStore.Set([]byte(B), IntToBytes(0))
	if err != nil {
		log.Println(err)
	}
	AVersion = version[B]
	AVersion.Check = 0
	version[B] = AVersion
}

func (BCstate *BlockchainState) GUpdateBalance(TxId uint16, A string, Balance int, m map[uint16]string, version map[string]AccountVersion) {
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
	v, ok := m[TxId]
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
		Check, err := BCstate.CheckingStore.Get([]byte(A))
		if err != nil {
			log.Println(err)
		}
		CheckInt = BytesToInt(Check)
	}

	CheckInt += Balance

	err := BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckInt))
	if err != nil {
		log.Println(err)
	}

	AVersion := version[A]
	AVersion.Check = CheckInt
	version[A] = AVersion
}

func (BCstate *BlockchainState) GUpdateSaving(TxId uint16, A string, Balance int, m map[uint16]string, version map[string]AccountVersion) {
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
	v, ok := m[TxId]
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
		Save, err := BCstate.SavingStore.Get([]byte(A))
		if err != nil {
			log.Println(err)
		}
		SaveInt = BytesToInt(Save)
	}

	SaveInt += Balance
	err := BCstate.CheckingStore.Set([]byte(A), IntToBytes(SaveInt))
	if err != nil {
		log.Println(err)
	}
	AVersion := version[A]
	AVersion.Save = SaveInt
	version[A] = AVersion
}

func (BCstate *BlockchainState) GSendPayment(TxId uint16, A string, B string, Balance int, m map[uint16]string, version map[string]AccountVersion) {
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

	v, ok := m[TxId]
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
		CheckA, err := BCstate.CheckingStore.Get([]byte(A))
		if err != nil {
			log.Println(err)
		}
		CheckIntA = BytesToInt(CheckA)
	}

	if CheckBS != "" {
		temp, _ := strconv.ParseInt(CheckBS, 10, 64)
		CheckIntB = int(temp)
	} else {
		CheckB, err := BCstate.CheckingStore.Get([]byte(B))
		if err != nil {
			log.Println(err)
		}
		CheckIntB = BytesToInt(CheckB)
	}

	CheckIntA -= Balance
	CheckIntB += Balance
	// update check value
	err := BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckIntA))
	if err != nil {
		log.Println(err)
	}
	err = BCstate.CheckingStore.Set([]byte(B), IntToBytes(CheckIntB))
	if err != nil {
		log.Println(err)
	}
	AVersion := version[A]
	AVersion.Check = CheckIntA
	version[A] = AVersion

	AVersion = version[B]
	AVersion.Check = CheckIntB
	version[B] = AVersion
}

func (BCstate *BlockchainState) GWriteCheck(TxId uint16, A string, Balance int, m map[uint16]string, version map[string]AccountVersion) {
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

	v, ok := m[TxId]
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
		Save, err := BCstate.SavingStore.Get([]byte(A))
		if err != nil {
			log.Println(err)
		}
		SaveInt = BytesToInt(Save)
	}

	if ConsistentCheckValue != "" {
		temp, _ := strconv.ParseInt(ConsistentCheckValue, 10, 64)
		CheckInt = int(temp)
	} else {
		Check, err := BCstate.CheckingStore.Get([]byte(A))
		if err != nil {
			log.Println(err)
		}
		CheckInt = BytesToInt(Check)
	}

	if SaveInt+CheckInt < Balance {
		CheckInt = CheckInt - Balance - 1
	} else {
		CheckInt = CheckInt - Balance
	}
	err := BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckInt))
	if err != nil {
		log.Println(err)
	}
	AVersion := version[A]
	AVersion.Check = CheckInt
	version[A] = AVersion
}
