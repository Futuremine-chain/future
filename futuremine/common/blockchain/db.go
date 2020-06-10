package blockchain

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IChainDB interface {
	ActRoot() arry.Hash
	DPosRoot() arry.Hash
	TokenRoot() arry.Hash
	LastHeight() uint64
}
