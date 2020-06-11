package token_status

import (
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/futuremine/db/status/token_db"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

const tokenDB = "token_db"

type TokenStatus struct {
	db ITokenDB
}

func NewTokenStatus() (*TokenStatus, error) {
	db, err := token_db.Open(config.App.Setting().Data + "/" + tokenDB)
	if err != nil {
		return nil, err
	}
	return &TokenStatus{db: db}, nil
}

func (t *TokenStatus) SetTrieRoot(hash arry.Hash) error {
	return t.db.SetRoot(hash)
}

func (t *TokenStatus) CheckMessage(msg types.IMessage) error {
	return nil
}
