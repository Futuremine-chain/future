package act_status

import (
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

type IActDB interface {
	SetRoot(hash arry.Hash) error
	Root() arry.Hash
	Commit() (arry.Hash, error)
	Close() error

	Account(address arry.Address) account.IAccount
	SetAccount(account account.IAccount)
	Nonce(address arry.Address) uint64
}
