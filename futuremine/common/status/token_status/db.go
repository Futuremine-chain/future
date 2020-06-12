package token_status

import (
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

type ITokenDB interface {
	SetRoot(hash arry.Hash) error
	Token(addr arry.Address) *types.TokenRecord
}
