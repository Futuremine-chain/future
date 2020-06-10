package types

import "github.com/Futuremine-chain/futuremine/tools/arry"

type ITransactionBody interface {
	To() arry.Address
}
