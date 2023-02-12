package abci

import (
	"encoding/json"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"onffchain/smallbankapplication/application"
)

// 实现abci接口
var _ abcitypes.Application = (*SmallBankApplication)(nil)

// 定义KVStore程序的结构体
type SmallBankApplication struct {
	abcitypes.BaseApplication
	Node *application.BlockchainState
}

// 创建一个 ABCI APP
func NewSmallBankApplication(node *application.BlockchainState) *SmallBankApplication {
	return &SmallBankApplication{
		Node: node,
	}
}

//func (app *SmallBankApplication) SetOption(option abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
//	//TODO implement me
//
//	panic("implement me")
//}

// 检查交易是否符合自己的要求，返回0时代表有效交易
func (app *SmallBankApplication) isValid(tx []byte) (code uint32) {
	/*// 格式校验，如果不是k=v格式的返回码为1
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 2 {
		return 1
	}

	key, value := parts[0], parts[1]

	//检查是否存在相同的KV
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == nil {
			return item.Value(func(val []byte) error {
				if bytes.Equal(val, value) {
					code = 2
				}
				return nil
			})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}*/

	return code
}

func (app *SmallBankApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	//app.currentBatch = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

// 当新的交易被添加到Tendermint Core时，它会要求应用程序进行检查(验证格式、签名等)，当返回0时才通过
func (app SmallBankApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	// Leader execute and send the sub graph
	//var events []abcitypes.Event
	var mSub []byte
	var mSubV []byte
	if app.Node.Leader {
		Sub, SubV := app.Node.ResolveAndExecuteTx(req.Tx)
		mSub, err := json.Marshal(Sub)
		if err != nil {
			return abcitypes.ResponseCheckTx{}
		}
		mSubV, err := json.Marshal(SubV)
		if err != nil {
			return abcitypes.ResponseCheckTx{}
		}
		/*events = []abcitypes.Event{
			{
				Type: "G",
				Attributes: []abcitypes.EventAttribute{
					{Key: "S", Value: string(mSub), Index: true},
					{Key: "SV", Value: string(mSubV), Index: true},
				},
			},
		}*/

	}

	return abcitypes.ResponseCheckTx{Code: 0, GasUsed: 0}
}

// 这里我们创建了一个batch，它将存储block的交易。
func (app *SmallBankApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	/*code := app.isValid(req.Tx)
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}

	parts := bytes.Split(req.Tx, []byte("="))
	key, value := parts[0], parts[1]

	err := app.currentBatch.Set(key, value)
	if err != nil {
		panic(err)
	}
	*/

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *SmallBankApplication) Commit() abcitypes.ResponseCommit {
	// 往数据库中提交事务，当 Tendermint core 提交区块时，会调用这个函数
	/*app.currentBatch.Commit()*/
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *SmallBankApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	/*resQuery.Key = reqQuery.Data
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(reqQuery.Data)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == badger.ErrKeyNotFound {
			resQuery.Log = "does not exist"
		} else {
			return item.Value(func(val []byte) error {
				resQuery.Log = "exists"
				resQuery.Value = val
				return nil
			})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}*/
	return
}

func (SmallBankApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (SmallBankApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (SmallBankApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (SmallBankApplication) ListSnapshots(abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
	return abcitypes.ResponseListSnapshots{}
}

func (SmallBankApplication) OfferSnapshot(abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
	return abcitypes.ResponseOfferSnapshot{}
}

func (SmallBankApplication) LoadSnapshotChunk(abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
	return abcitypes.ResponseLoadSnapshotChunk{}
}

func (SmallBankApplication) ApplySnapshotChunk(abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
	return abcitypes.ResponseApplySnapshotChunk{}
}
