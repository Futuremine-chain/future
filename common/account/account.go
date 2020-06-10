package account

import "github.com/Futuremine-chain/futuremine/types"

type IAccount interface {
	NeedUpdate() bool
	UpdateLocked(confirmed uint64) error
	FromTransaction(tx types.ITransaction, height uint64) error
	ToTransaction(tx types.ITransaction, height uint64) error
	Check(tx types.ITransaction) error
}
