package blockchain

import (
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

type IChainDB interface {
	ActRoot() (arry.Hash, error)
	DPosRoot() (arry.Hash, error)
	TokenRoot() (arry.Hash, error)
	LastHeight() uint64
	GetMessages(txRoot arry.Hash) ([]*types.RlpMessage, error)
	GetHeaderHeight(height uint64) (*types.Header, error)
	GetHeaderHash(hash arry.Hash) (*types.Header, error)
}
