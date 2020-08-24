package types

import "github.com/Futuremine-chain/futuremine/types"

type Body struct {
	Messages
}

func (b *Body) Msgs() []types.IMessage {
	return b.Messages
}

func (b *Body) ToRlpBody() *RlpBody {
	rMsgs := make([]*RlpMessage, b.Messages.Count())
	for i, msg := range b.Messages {
		rMsgs[i] = msg.ToRlp().(*RlpMessage)
	}
	return &RlpBody{rMsgs}
}
