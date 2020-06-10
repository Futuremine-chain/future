package token_db

import (
	"github.com/Futuremine-chain/futuremine/common/db/base"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/trie"
)

type TokenDB struct {
	base *base.Base
	trie *trie.Trie
}

func Open(path string) (*TokenDB, error) {
	var err error
	baseDB, err := base.Open(path)
	if err != nil {
		return nil, err
	}
	return &TokenDB{base: baseDB}, nil
}

func (t *TokenDB) SetRoot(hash arry.Hash) error {
	tri, err := trie.New(hash, t.base)
	if err != nil {
		return err
	}
	t.trie = tri
	return nil
}
