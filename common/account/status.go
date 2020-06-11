package account

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IActStatus interface {
	Nonce(arry.Address) uint64
	SetTrieRoot(hash arry.Hash) error
	SetConfirmed(confirmed uint64)
	CheckTransaction(tx types.ITransaction) error
}
