package account

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IAccount interface {
	NeedUpdate() bool
	UpdateLocked(confirmed uint64) error
	FromMessage(msg types.IMessage, height uint64) error
	ToMessage(msg types.IMessage, height uint64) error
	Check(msg types.IMessage) error
	Bytes() []byte
	Address() arry.Address
}
