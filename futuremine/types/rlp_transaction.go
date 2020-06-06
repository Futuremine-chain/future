package types

import "github.com/Futuremine-chain/futuremine/types"

type RlpTransaction struct {
}

func (r *RlpTransaction) ToTransaction() types.ITransaction {
	return nil
}

func DecodeTransaction([]byte) (*RlpTransaction, error) {
	return nil, nil
}
