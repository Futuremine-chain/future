package types

type ITransaction interface {
	ITransactionHeader
	ITransactionBody

	ToRlp() IRlpTransaction
}
