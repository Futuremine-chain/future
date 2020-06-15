package account

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IActStatus interface {
	Nonce(arry.Address) uint64
	SetTrieRoot(hash arry.Hash) error
	TrieRoot() arry.Hash
	SetConfirmed(confirmed uint64)
	Account(address arry.Address) IAccount
	CheckMessage(msg types.IMessage, strict bool) error
	FromMessage(msg types.IMessage, height uint64) error
	ToMessage(msg types.IMessage, height uint64) error
	Commit() (arry.Hash, error)
}
