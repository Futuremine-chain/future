package types

import "github.com/Futuremine-chain/futuremine/futuremine/common/arry"

type ITransactionHeader interface {
	Hash() arry.Hash
}
