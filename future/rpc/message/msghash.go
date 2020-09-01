package message

import (
	"encoding/json"
	rpctypes "github.com/Futuremine-chain/future/future/rpc/types"
	"github.com/Futuremine-chain/future/tools/arry"
	"github.com/Futuremine-chain/future/tools/crypto/hash"
	"github.com/Futuremine-chain/future/types"
)

func MessageHash(msg types.IMessage) (arry.Hash, error) {
	rpcMsg, err := rpctypes.MsgToRpcMsg(msg)
	if err != nil {
		return arry.Hash{}, err
	}
	mBytes, err := json.Marshal(rpcMsg)
	if err != nil {
		return arry.Hash{}, err
	}
	return hash.Hash(mBytes), nil
}
