package act_db

import (
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/common/db/base"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/trie"
)

type ActDB struct {
	base *base.Base
	trie *trie.Trie
}

func Open(path string) (*ActDB, error) {
	return &ActDB{}, nil
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
	panic("implement me")
}

func (a *ActDB) SetAccount(account account.IAccount) {
	panic("implement me")
}

func (a *ActDB) Nonce(address arry.Address) uint64 {
	panic("implement me")
}

func (a *ActDB) Close() error {
	return a.base.Close()
}

/*
func (s *StateStorage) InitTrie(stateRoot hasharry.Hash) error {
	stateTrie, err := trie.New(stateRoot, s.trieDB)
	if err != nil {
		return err
	}
	s.stateTrie = stateTrie
	return nil
}

func (s *StateStorage) Open() error {
	return s.trieDB.Open()
}



func (s *StateStorage) GetAccountState(stateKey hasharry.Hash) types.IAccount {
	account := types.NewAccount()
	bytes := s.stateTrie.Get(stateKey.Bytes())
	err := rlp.DecodeBytes(bytes, &account)
	if err != nil {
		return types.NewAccount()
	}
	return account
}

func (s *StateStorage) SetAccountState(account types.IAccount) {
	bytes, err := rlp.EncodeToBytes(account.(*types.Account))
	if err != nil {
		return
	}
	s.stateTrie.Update(account.StateKey().Bytes(), bytes)
}

func (s *StateStorage) GetAccountBalance(stateKey hasharry.Hash, contract string) uint64 {
	account := types.NewAccount()
	bytes := s.stateTrie.Get(stateKey.Bytes())
	err := rlp.DecodeBytes(bytes, &account)
	if err != nil {
		return 0
	}
	return account.GetBalance(contract)
}

func (s *StateStorage) GetAccountNonce(stateKey hasharry.Hash) uint64 {
	account := types.NewAccount()
	bytes := s.stateTrie.Get(stateKey.Bytes())
	err := rlp.DecodeBytes(bytes, &account)
	if err != nil {
		return 0
	}
	return account.GetNonce()
}

func (s *StateStorage) DeleteAccount(stateKey hasharry.Hash) {
	s.stateTrie.Delete(stateKey.Bytes())
}

func (s *StateStorage) Commit() (hasharry.Hash, error) {
	return s.stateTrie.Commit()
}

func (s *StateStorage) RootHash() hasharry.Hash {
	return s.stateTrie.Hash()
}*/
