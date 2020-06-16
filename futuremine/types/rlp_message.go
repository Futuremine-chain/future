package types

import (
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/types"
)

type RlpMessage struct {
	MsgHeader *MsgHeader
	MsgBody   []byte
}

func (r *RlpMessage) ToMessage() types.IMessage {
	return nil
}

func (r *RlpMessage) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(r)
	return bytes
}

func DecodeMessage([]byte) (*RlpMessage, error) {
	return nil, nil
}
