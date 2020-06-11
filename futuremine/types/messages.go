package types

import "github.com/Futuremine-chain/futuremine/types"

type Messages []*Message

func (t Messages) Msgs() []types.IMessage {
	iTxs := make([]types.IMessage, len(t))
	for i, msg := range t.Msgs() {
		iTxs[i] = msg
	}
	return iTxs
}

func (t Messages) Add(iMsg types.IMessage) {
	iMsg = new(Message)
	msg := iMsg.(*Message)
	t = append(t, msg)
}
