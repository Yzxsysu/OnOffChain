package application

import (
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

var Client *rpchttp.HTTP

func (BCstate *BlockchainState) Validate(s []SmallBankTransaction, GE [][]GraphEdge, u [][]uint16, groupNum int) {
	v := u[groupNum]
	version := make(map[string]AccountVersion)
	ch := make(chan bool, 2)
	go BCstate.VValidate(s, v, ch)
	go BCstate.GValidate(s, GE, groupNum, version, ch)
	// go OffChainExecute
	<-ch
	<-ch
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

func GetEventFromCheckTx() {

}

// GetGE get the GraphEdge from the Leader
func GetGE() {

}

// GetU get the GraphEdge from the Leader
func GetU() {

}
