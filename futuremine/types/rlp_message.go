package types

import "github.com/Futuremine-chain/futuremine/types"

type RlpMessage struct {
}

func (r *RlpMessage) ToMessage() types.IMessage {
	return nil
}

func DecodeMessage([]byte) (*RlpMessage, error) {
	return nil, nil
}
