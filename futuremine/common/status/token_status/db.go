package token_status

import "github.com/Futuremine-chain/futuremine/tools/arry"

type ITokenDB interface {
	SetRoot(hash arry.Hash) error
}
