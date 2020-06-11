package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type Transaction struct {
}

func (t *Transaction) Hash() arry.Hash {
	panic("implement me")
}

func (t *Transaction) From() arry.Address {
	panic("implement me")
}

func (t *Transaction) Nonce() uint64 {
	panic("implement me")
}

func (t *Transaction) Fee() uint64 {
	panic("implement me")
}

func (t *Transaction) Time() int64 {
	panic("implement me")
}

func (t *Transaction) IsCoinBase() bool {
	panic("implement me")
}

func (t *Transaction) To() arry.Address {
	panic("implement me")
}

func (t *Transaction) ToRlp() types.IRlpTransaction {
	panic("implement me")
}

func (t *Transaction) Check() error {
	return nil
}
