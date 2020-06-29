package types

import "github.com/Futuremine-chain/futuremine/futuremine/types"

type RpcBody struct {
	Transactions []*RpcMessage `json:"transactions"`
}

func BodyToRpcBody(body *types.Body) (*RpcBody, error) {
	var rpcMsgs []*RpcMessage
	for _, msg := range body.Messages {
		rpcMsg, err := MsgToRpcMsg(msg.(*types.Message))
		if err != nil {
			return nil, err
		}
		rpcMsgs = append(rpcMsgs, rpcMsg)
	}
	return &RpcBody{rpcMsgs}, nil
}
