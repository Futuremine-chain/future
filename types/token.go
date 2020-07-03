package types

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IToken interface {
}

type ITokenStatus interface {
	SetTrieRoot(hash arry.Hash) error
	TrieRoot() arry.Hash
	CheckMessage(msg IMessage) error
	UpdateToken(msg IMessage, height uint64) error
	Token(address arry.Address) (IToken, error)
	Commit() (arry.Hash, error)
}
