package types

import (
	fmctypes "github.com/Futuremine-chain/future/future/types"
	"github.com/Futuremine-chain/future/types"
)

type TxPool struct {
	MsgsCount  int                    `json:"msgs"`
	ReadyCount int                    `json:"ready"`
	CacheCount int                    `json:"cache"`
	ReadyMsgs  []*fmctypes.RpcMessage `json:"readymsgs"`
	CacheMsgs  []*fmctypes.RpcMessage `json:"cachemsgs"`
}

func MsgsToRpcMsgsPool(readyMsgs []types.IMessage, cacheMsgs []types.IMessage) *TxPool {
	var readyRpcMsgs, cacheRpcMsgs []*fmctypes.RpcMessage
	for _, msg := range readyMsgs {
		rpcMsg, _ := fmctypes.MsgToRpcMsg(msg.(*fmctypes.Message))
		readyRpcMsgs = append(readyRpcMsgs, rpcMsg)
	}

	for _, msg := range cacheMsgs {
		rpcMsg, _ := fmctypes.MsgToRpcMsg(msg.(*fmctypes.Message))
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
