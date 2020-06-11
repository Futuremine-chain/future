package types

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IMessageHeader interface {
	Type() int
	Hash() arry.Hash
	From() arry.Address
	Nonce() uint64
	Fee() uint64
	Time() int64
	IsCoinBase() bool
}
