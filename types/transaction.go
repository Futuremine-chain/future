package types

type ITransaction interface {
	ITransactionHeader
	ITransactionBody
}
