package token

import "github.com/Futuremine-chain/futuremine/tools/arry"

type ITokenStatus interface {
	SetTrieRoot(hash arry.Hash) error
}
