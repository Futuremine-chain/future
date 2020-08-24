package types

import (
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
)

type RpcBlock struct {
	RpcHeader *RpcHeader `json:"header"`
	RpcBody   *RpcBody   `json:"body"`
	Confirmed bool       `json:"confirmed"`
}

func BlockToRpcBlock(block *fmctypes.Block, confirmed uint64) (*RpcBlock, error) {
	rpcHeader := HeaderToRpcHeader(block.Header)
	rpcBody, err := BodyToRpcBody(block.Body)
	if err != nil {
		return nil, err
	}
	return &RpcBlock{
		RpcHeader: rpcHeader,
		RpcBody:   rpcBody,
		Confirmed: confirmed >= rpcHeader.Height,
	}, nil
}
