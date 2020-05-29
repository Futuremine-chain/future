package validator

import "github.com/Futuremine-chain/futuremine/types"

type IValidator interface {
	Check(transaction types.ITransaction) error
}
