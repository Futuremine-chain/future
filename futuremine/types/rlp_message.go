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
	msg := &Message{}
	msg.Header = r.MsgHeader
	switch r.MsgHeader.Type {
	case Transaction:
		var body *TransactionBody
		rlp.DecodeBytes(r.MsgBody, &body)
		msg.Body = body
	case Token:
		var body *TokenBody
		rlp.DecodeBytes(r.MsgBody, &body)
		msg.Body = body
	case Candidate:
		var body *CandidateBody
		rlp.DecodeBytes(r.MsgBody, &body)
		msg.Body = body
	case Cancel:
		var body *CancelBody
		rlp.DecodeBytes(r.MsgBody, &body)
		msg.Body = body
	case Vote:
		var body *VoteBody
		rlp.DecodeBytes(r.MsgBody, &body)
		msg.Body = body
	}
	return msg
}

func (r *RlpMessage) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(r)
	return bytes
}

func DecodeMessage([]byte) (*RlpMessage, error) {
	return nil, nil
}
