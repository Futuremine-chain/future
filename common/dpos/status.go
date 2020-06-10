package dpos

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IDPosStatus interface {
	SetTrieRoot(hash arry.Hash) error
}
