package types

import (
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/types"
)

type TxPool struct {
	MsgsCount  int           `json:"msgs"`
	ReadyCount int           `json:"ready"`
	CacheCount int           `json:"cache"`
	ReadyMsgs  []*RpcMessage `json:"readymsgs"`
	CacheMsgs  []*RpcMessage `json:"cachemsgs"`
}

func MsgsToRpcMsgsPool(readyMsgs []types.IMessage, cacheMsgs []types.IMessage) *TxPool {
	var readyRpcMsgs, cacheRpcMsgs []*RpcMessage
	for _, msg := range readyMsgs {
		rpcMsg, _ := MsgToRpcMsg(msg.(*fmctypes.Message))
		readyRpcMsgs = append(readyRpcMsgs, rpcMsg)
	}

	for _, msg := range cacheMsgs {
		rpcMsg, _ := MsgToRpcMsg(msg.(*fmctypes.Message))
		cacheRpcMsgs = append(cacheRpcMsgs, rpcMsg)
	}

	readyCount := len(readyRpcMsgs)
	cacheCount := len(cacheRpcMsgs)

	return &TxPool{
		MsgsCount:  readyCount + cacheCount,
		ReadyCount: readyCount,
		CacheCount: cacheCount,
		ReadyMsgs:  readyRpcMsgs,
		CacheMsgs:  cacheRpcMsgs,
	}
}
