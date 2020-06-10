package types

import "github.com/Futuremine-chain/futuremine/types"

type Transactions []*Transaction

func (t Transactions) Txs() []types.ITransaction {
	iTxs := make([]types.ITransaction, len(t))
	for i, tx := range t.Txs() {
		iTxs[i] = tx
	}
	return iTxs
}

func (t Transactions) Add(iTx types.ITransaction) {
	tx := iTx.(*Transaction)
	t = append(t, tx)
}
