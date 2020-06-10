package act_db

import (
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/common/db/base"
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/trie"
)

type ActDB struct {
	base *base.Base
	trie *trie.Trie
}

func Open(path string) (*ActDB, error) {
	var err error
	baseDB, err := base.Open(path)
	if err != nil {
		return nil, err
	}
	return &ActDB{base: baseDB}, nil
}

func (a *ActDB) SetRoot(hash arry.Hash) error {
	t, err := trie.New(hash, a.base)
	if err != nil {
		return err
	}
	a.trie = t
	return nil
}

func (a *ActDB) Root() arry.Hash {
	panic("implement me")
}

func (a *ActDB) Commit() (arry.Hash, error) {
	panic("implement me")
}

func (a *ActDB) Account(address arry.Address) account.IAccount {
	bytes := a.trie.Get(address.Bytes())
	if account, err := types.DecodeAccount(bytes); err != nil {
		return types.NewAccount()
	} else {
		return account
	}
}

func (a *ActDB) SetAccount(account account.IAccount) {
	a.trie.Update(account.Address().Bytes(), account.Bytes())
}

func (a *ActDB) Nonce(address arry.Address) uint64 {
	bytes := a.trie.Get(address.Bytes())
	account, err := types.DecodeAccount(bytes)
	if err != nil {
		return 0
	}
	return account.Nonce()
}

func (a *ActDB) Close() error {
	return a.base.Close()
}
