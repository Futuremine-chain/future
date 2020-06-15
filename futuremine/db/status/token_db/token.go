package token_db

import (
	"github.com/Futuremine-chain/futuremine/common/db/base"
	"github.com/Futuremine-chain/futuremine/futuremine/types"
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

func (t *TokenDB) Commit() (arry.Hash, error) {
	return t.trie.Commit()
}

func (t *TokenDB) Root() arry.Hash {
	return t.trie.Hash()
}

func (t *TokenDB) Close() error {
	return t.base.Close()
}

func (t *TokenDB) Token(address arry.Address) *types.TokenRecord {
	bytes := t.trie.Get(address.Bytes())
	token, err := types.DecodeToken(bytes)
	if err != nil {
		return nil
	}
	return token
}

func (t *TokenDB) SetToken(token *types.TokenRecord) {
	t.trie.Update(token.Address.Bytes(), token.Bytes())
}
