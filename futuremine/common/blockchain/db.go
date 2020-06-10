package blockchain

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IChainDB interface {
	ActRoot() (arry.Hash, error)
	DPosRoot() (arry.Hash, error)
	TokenRoot() (arry.Hash, error)
	LastHeight() uint64
}
