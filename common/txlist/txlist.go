package txlist

import "github.com/Futuremine-chain/futuremine/types"

type ITxList interface {
	DeleteExpired(int64)
	DeleteEnd(types.ITransaction)
	DeleteAndUpdate(transactions types.ITransactions)
	Read() error
	Close() error
	Update()
	Exist(types.ITransaction) bool
	Put(types.ITransaction) error
	Count() int
}
