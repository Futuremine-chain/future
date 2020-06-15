package token_status

import (
	"errors"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/futuremine/db/status/token_db"
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
	"sync"
)

const tokenDB = "token_db"

type TokenStatus struct {
	db    ITokenDB
	mutex sync.RWMutex
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

func (t *TokenStatus) TrieRoot() arry.Hash {
	return t.db.Root()
}

func (t *TokenStatus) Commit() (arry.Hash, error) {
	return t.db.Commit()
}

func (t *TokenStatus) CheckMessage(msg types.IMessage) error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if fmctypes.MessageType(msg.Type()) != fmctypes.Token {
		return nil
	}
	body, ok := msg.MsgBody().(*fmctypes.TokenBody)
	if !ok {
		return errors.New("incorrect message type and message body")
	}
	token := t.db.Token(body.TokenAddress)
	if token != nil {
		return token.Check(msg)
	}
	return nil
}

// Update contract status
func (t *TokenStatus) UpdateToken(msg types.IMessage, height uint64) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	msgBody, ok := msg.MsgBody().(*fmctypes.TokenBody)
	if !ok {
		return errors.New("wrong message type")
	}
	record := &fmctypes.Record{
		Height:  height,
		MsgHash: msg.Hash(),
		Time:    msg.Time(),
		Amount:  msgBody.MsgAmount(),
	}
	tokenAddr := msgBody.TokenAddress
	token := t.db.Token(tokenAddr)
	if token != nil {
		token.AddContract(record)
	} else {
		token = &fmctypes.TokenRecord{
			Address:   tokenAddr,
			Name:      msgBody.Name,
			Shorthand: msgBody.Shorthand,
			Records: &fmctypes.RecordList{
				record,
			},
		}
	}
	t.db.SetToken(token)
	return nil
}
