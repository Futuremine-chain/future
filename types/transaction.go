package types

import "github.com/Futuremine-chain/futuremine/types/rlp"

type ITransaction interface {
	ITransactionHeader
	ITransactionBody

	ToRlp() rlp.IRlpTransaction
}
