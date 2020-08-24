package types

import "github.com/Futuremine-chain/future/future/types"

type RpcBody struct {
	Messages []*RpcMessage `json:"transactions"`
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
