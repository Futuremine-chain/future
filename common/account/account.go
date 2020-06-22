package account

import (
	rpctypes "github.com/Futuremine-chain/futuremine/service/rpc/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IAccount interface {
	NeedUpdate() bool
	UpdateLocked(confirmed uint64) error
	FromMessage(msg types.IMessage, height uint64) error
	ToMessage(msg types.IMessage, height uint64) error
	Check(msg types.IMessage, strict bool) error
	Bytes() []byte
	GetAddress() arry.Address
	GetBalance(tokenAddr arry.Address) uint64
}
