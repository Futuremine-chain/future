package msglist

import "github.com/Futuremine-chain/futuremine/types"

type ITxListDB interface {
	Read() []types.IMessage
	Save(message types.IMessage)
	Delete(msg types.IMessage)
	Clear()
	Close() error
}
