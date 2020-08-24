package token_status

import (
	"github.com/Futuremine-chain/future/future/types"
	"github.com/Futuremine-chain/future/tools/arry"
)

type ITokenDB interface {
	SetRoot(hash arry.Hash) error
	Root() arry.Hash
	Commit() (arry.Hash, error)
	Token(addr arry.Address) *types.TokenRecord
	SetToken(token *types.TokenRecord)
}
