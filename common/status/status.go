package status

import (
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IStatus interface {
	InitRoots(actRoot, dPosRoot, tokenRoot arry.Hash) error
	SetConfirmed(confirmed uint64)
	Account(address arry.Address) account.IAccount
	CheckMsg(msg types.IMessage, strict bool) error
	Change(msgs []types.IMessage, block types.IBlock) error
	Commit() (arry.Hash, arry.Hash, arry.Hash, error)
}
