package types

import "github.com/Futuremine-chain/futuremine/futuremine/common/arry"

type ITransactionHeader interface {
	Hash() arry.Hash
	From() arry.Address
	Nonce() uint64
	Fee() uint64
	Time() int64
}
