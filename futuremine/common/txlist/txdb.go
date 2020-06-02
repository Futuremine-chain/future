package txlist

import "github.com/Futuremine-chain/futuremine/types"

type ITxListDB interface {
	Read() []types.ITransaction
	Save(transaction types.ITransaction)
	Delete(tx types.ITransaction)
	Clear()
	Close() error
}
