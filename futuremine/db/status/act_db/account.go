package act_db

import (
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/common/db/base"
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

type ActDB struct {
	base.Base
}

func Open(path string) (*ActDB, error) {
	return &ActDB{}, nil
}

func (a *ActDB) SetRoot(hash arry.Hash) error {
	panic("implement me")
}

func (a *ActDB) Root() arry.Hash {
	panic("implement me")
}

func (a *ActDB) Commit() (arry.Hash, error) {
	panic("implement me")
}

func (a *ActDB) Account(address arry.Address) account.IAccount {
	panic("implement me")
}

func (a *ActDB) SetAccount(account account.IAccount) {
	panic("implement me")
}

func (a *ActDB) Nonce(address arry.Address) uint64 {
	panic("implement me")
}
