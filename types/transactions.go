package types

type ITransactions interface {
	Add(transaction ITransaction)
	Txs() []ITransaction
}
