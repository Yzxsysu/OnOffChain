package application

import (
	"fmt"
)

var MsgS chan [][]GraphEdge
var MsgSV chan [][]uint16
var MsgV1 chan map[string]AccountVersion
var MsgV2 chan map[string]AccountVersion
var MsgV3 chan map[string]AccountVersion
var MsgV4 chan map[string]AccountVersion
var MsgV5 chan map[string]AccountVersion
var MsgV6 chan map[string]AccountVersion

var Version chan map[string]AccountVersion

func (BCstate *BlockchainState) Validate(s *[]SmallBankTransaction, group int) {
	GE := <-MsgS
	u := <-MsgSV
	v := u[group]
	ch := make(chan bool, 2)
	go BCstate.VValidate(s, v, ch)
	go BCstate.GValidate(s, GE, group, Version, ch)
	// go OffChainExecute
	<-ch
	<-ch
}

var SetNum string
var Group int

func (BCstate *BlockchainState) DValidate(s *[]SmallBankTransaction) {
	if SetNum == "f" {
		if Group == 1 {
			go BCstate.Validate(s, 1)
			v2 := <-MsgV2
			v3 := <-MsgV3
			v5 := <-MsgV5
			v6 := <-MsgV6
			v := <-Version
			BCstate.MergeS(v, v2)
			BCstate.MergeS(v, v3)
			go BCstate.MergeSV(v5)
			go BCstate.MergeSV(v6)
		}
		if Group == 2 {
			go BCstate.Validate(s, 2)
			v1 := <-MsgV1
			v3 := <-MsgV3
			v4 := <-MsgV4
			v6 := <-MsgV6
			v := <-Version
			BCstate.MergeS(v, v1)
			BCstate.MergeS(v, v3)
			go BCstate.MergeSV(v4)
			go BCstate.MergeSV(v6)
		}
		if Group == 3 {
			go BCstate.Validate(s, 3)
			v1 := <-MsgV1
			v2 := <-MsgV2
			v4 := <-MsgV4
			v5 := <-MsgV5
			v := <-Version
			BCstate.MergeS(v, v1)
			BCstate.MergeS(v, v2)
			go BCstate.MergeSV(v4)
			go BCstate.MergeSV(v5)
		}
	}
	if SetNum == "2f" {
		if Group == 1 {
			go BCstate.Validate(s, 1)
			go BCstate.Validate(s, 2)
			v3 := <-MsgV3
			v6 := <-MsgV6
			v := <-Version
			BCstate.MergeS(v, v3)
			go BCstate.MergeSV(v6)
		}
		if Group == 2 {
			go BCstate.Validate(s, 2)
			go BCstate.Validate(s, 3)
			v1 := <-MsgV1
			v4 := <-MsgV4
			v := <-Version
			BCstate.MergeS(v, v1)
			go BCstate.MergeSV(v4)
		}
		if Group == 3 {
			go BCstate.Validate(s, 1)
			go BCstate.Validate(s, 3)
			v2 := <-MsgV2
			v5 := <-MsgV5
			v := <-Version
			BCstate.MergeS(v, v2)
			go BCstate.MergeSV(v5)
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
	for key, value := range NewV {
		_, err := OriginV[key]
		if err != true {
			OriginV[key] = value
		}

		if OriginV[key].SaveVersion < value.SaveVersion {
			err := BCstate.SavingStore.Set([]byte(key), IntToBytes(value.Save))
			if err != nil {
				fmt.Println(err)
			}
			tempV := OriginV[key]
			tempV.SaveVersion = value.SaveVersion
			tempV.Save = value.Save
			OriginV[key] = tempV
		}
		if OriginV[key].CheckVersion < value.CheckVersion {
			err := BCstate.CheckingStore.Set([]byte(key), IntToBytes(value.Check))
			if err != nil {
				fmt.Println(err)
			}
			tempV := OriginV[key]
			tempV.CheckVersion = value.CheckVersion
			tempV.Check = value.Check
			OriginV[key] = tempV
		}
	}
}

func (BCstate *BlockchainState) MergeSV(NewV map[string]AccountVersion) {
	for key, value := range NewV {
		err := BCstate.SavingStore.Set([]byte(key), IntToBytes(value.Save))
		if err != nil {
			fmt.Println(err)
		}
		err = BCstate.CheckingStore.Set([]byte(key), IntToBytes(value.Check))
		if err != nil {
			fmt.Println(err)
		}
	}
}
