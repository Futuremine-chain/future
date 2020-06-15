package token

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type ITokenStatus interface {
	SetTrieRoot(hash arry.Hash) error
	TrieRoot() arry.Hash
	CheckMessage(msg types.IMessage) error
}
