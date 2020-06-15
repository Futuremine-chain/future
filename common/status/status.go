package status

import (
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IStatus interface {
	InitRoots(actRoot, dPosRoot, tokenRoot arry.Hash) error
	SetConfirmed(confirmed uint64)
	Account(address arry.Address) account.IAccount
	CheckBlock(block types.IBlock, chain blockchain.IChain) error
	CheckMsg(msg types.IMessage, strict bool) error
}
