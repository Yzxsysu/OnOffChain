package application

import (
	"log"
)

var bufferLength = 1000
var MsgS = make(chan [][]GraphEdge, bufferLength)
var MsgSV = make(chan [][]uint16, bufferLength)
var MsgV1 = make(chan map[string]AccountVersion, bufferLength)
var MsgV2 = make(chan map[string]AccountVersion, bufferLength)
var MsgV3 = make(chan map[string]AccountVersion, bufferLength)
var MsgV4 = make(chan map[string]AccountVersion, bufferLength)
var MsgV5 = make(chan map[string]AccountVersion, bufferLength)
var MsgV6 = make(chan map[string]AccountVersion, bufferLength)
var Version = make(chan map[string]AccountVersion, bufferLength)

func (BCstate *BlockchainState) Validate(s []SmallBankTransaction, GE *[][]GraphEdge, u *[][]uint16, group int) {
	//log.Println("Validate:", group)
	if len(*u) == 0 {
		*u = [][]uint16{{0}, {0}, {0}}
	}
	v := (*u)[group]
	ch := make(chan bool, 2)
	go BCstate.VValidate(s, &v, ch)
	go BCstate.GValidate(s, GE, group, Version, ch)
	// go OffChainExecute
	<-ch
	<-ch
	//log.Println("Validate end", Group)
}

var SetNum string
var Group int

func (BCstate *BlockchainState) DValidate(s []SmallBankTransaction) {
	if SetNum == "2f" {
		if Group == 1 {
			//log.Println("DValidate:", Group)
			GE := <-MsgS
			//log.Println("GE:", GE)
			u := <-MsgSV
			//log.Println("u:", u)
			go BCstate.Validate(s, &GE, &u, 1-1)
			go BCstate.Validate(s, &GE, &u, 2-1)
			// 可能是这里被阻塞了
			// offchain的问题 并且和链上的goroutine的channel相关
			go BCstate.MergeSV(<-MsgV6)
			ver1 := <-Version
			ver2 := <-Version
			BCstate.MergeS2(ver1, ver2, <-MsgV3)
		}
		if Group == 2 {
			GE := <-MsgS
			u := <-MsgSV
			go BCstate.Validate(s, &GE, &u, 2-1)
			go BCstate.Validate(s, &GE, &u, 3-1)
			go BCstate.MergeSV(<-MsgV4)
			ver1 := <-Version
			ver2 := <-Version
			BCstate.MergeS2(ver1, ver2, <-MsgV1)
		}
		if Group == 3 {
			GE := <-MsgS
			u := <-MsgSV
			go BCstate.Validate(s, &GE, &u, 1-1)
			go BCstate.Validate(s, &GE, &u, 3-1)
			go BCstate.MergeSV(<-MsgV5)
			ver1 := <-Version
			ver2 := <-Version
			BCstate.MergeS2(ver1, ver2, <-MsgV2)
		}
	}
	if SetNum == "f" {
		if Group == 1 {
			GE := <-MsgS
			u := <-MsgSV
			go BCstate.Validate(s, &GE, &u, 1-1)
			go BCstate.MergeSV(<-MsgV5)
			go BCstate.MergeSV(<-MsgV6)
			v := <-Version
			BCstate.MergeS(v, <-MsgV2)
			BCstate.MergeS(v, <-MsgV3)
		}
		if Group == 2 {
			GE := <-MsgS
			u := <-MsgSV
			go BCstate.Validate(s, &GE, &u, 2-1)
			go BCstate.MergeSV(<-MsgV4)
			go BCstate.MergeSV(<-MsgV6)
			v := <-Version
			BCstate.MergeS(v, <-MsgV1)
			BCstate.MergeS(v, <-MsgV3)
		}
		if Group == 3 {
			GE := <-MsgS
			u := <-MsgSV
			go BCstate.Validate(s, &GE, &u, 3-1)
			go BCstate.MergeSV(<-MsgV4)
			go BCstate.MergeSV(<-MsgV5)
			v := <-Version
			BCstate.MergeS(v, <-MsgV1)
			BCstate.MergeS(v, <-MsgV2)
		}
	}
}

type AccountVersion struct {
	Check        int
	CheckVersion uint16
	Save         int
	SaveVersion  uint16
}

func NewAccountVersion() AccountVersion {
	return AccountVersion{}
}

func (BCstate *BlockchainState) MergeS(OriginV map[string]AccountVersion, NewV map[string]AccountVersion) {
	//log.Println("MergeS:", Group)
	for key, value := range NewV {
		_, err := OriginV[key]
		if err != true {
			OriginV[key] = value
		}

		if OriginV[key].SaveVersion <= value.SaveVersion {
			err := BCstate.SavingStore.Set([]byte(key), IntToBytes(value.Save))
			if err != nil {
				log.Println(err)
			}
			tempV := OriginV[key]
			tempV.SaveVersion = value.SaveVersion
			tempV.Save = value.Save
			OriginV[key] = tempV
		}
		if OriginV[key].CheckVersion <= value.CheckVersion {
			err := BCstate.CheckingStore.Set([]byte(key), IntToBytes(value.Check))
			if err != nil {
				log.Println(err)
			}
			tempV := OriginV[key]
			tempV.CheckVersion = value.CheckVersion
			tempV.Check = value.Check
			OriginV[key] = tempV
		}
	}
	//log.Println("MergeS end")
}

func (BCstate *BlockchainState) MergeS2(OriginV map[string]AccountVersion, NewV1 map[string]AccountVersion, NewV2 map[string]AccountVersion) {
	//log.Println("MergeS2:", Group)
	for key, value := range NewV1 {
		_, err := OriginV[key]
		if err != true {
			OriginV[key] = value
		}

		if OriginV[key].SaveVersion <= value.SaveVersion {
			err := BCstate.SavingStore.Set([]byte(key), IntToBytes(value.Save))
			if err != nil {
				log.Println(err)
			}
			tempV := OriginV[key]
			tempV.SaveVersion = value.SaveVersion
			tempV.Save = value.Save
			OriginV[key] = tempV
		}
		if OriginV[key].CheckVersion <= value.CheckVersion {
			err := BCstate.CheckingStore.Set([]byte(key), IntToBytes(value.Check))
			if err != nil {
				log.Println(err)
			}
			tempV := OriginV[key]
			tempV.CheckVersion = value.CheckVersion
			tempV.Check = value.Check
			OriginV[key] = tempV
		}
	}
	for key, value := range NewV2 {
		_, err := OriginV[key]
		if err != true {
			OriginV[key] = value
		}

		if OriginV[key].SaveVersion <= value.SaveVersion {
			err := BCstate.SavingStore.Set([]byte(key), IntToBytes(value.Save))
			if err != nil {
				log.Println(err)
			}
			tempV := OriginV[key]
			tempV.SaveVersion = value.SaveVersion
			tempV.Save = value.Save
			OriginV[key] = tempV
		}
		if OriginV[key].CheckVersion <= value.CheckVersion {
			err := BCstate.CheckingStore.Set([]byte(key), IntToBytes(value.Check))
			if err != nil {
				log.Println(err)
			}
			tempV := OriginV[key]
			tempV.CheckVersion = value.CheckVersion
			tempV.Check = value.Check
			OriginV[key] = tempV
		}
	}
	//log.Println("MergeS2 end")
}

func (BCstate *BlockchainState) MergeSV(NewV map[string]AccountVersion) {
	//log.Println("MergeSV", Group)
	for key, value := range NewV {
		err := BCstate.SavingStore.Set([]byte(key), IntToBytes(value.Save))
		if err != nil {
			log.Println(err)
		}
		err = BCstate.CheckingStore.Set([]byte(key), IntToBytes(value.Check))
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("MergeSV end")
}
