package blockchain

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IChain interface {
	LastHeight() uint64
	LastHeader() (types.IHeader, error)
	LastConfirmed() uint64
	SetConfirmed(uint64)
	GetBlockHeight(uint64) (types.IBlock, error)
	GetBlockHash(arry.Hash) (types.IBlock, error)
	GetHeaderHeight(uint64) (types.IHeader, error)
	GetHeaderHash(arry.Hash) (types.IHeader, error)

	CycleLastHash(uint64) (arry.Hash, error)

	GetRlpBlockHeight(uint64) (types.IRlpBlock, error)
	GetRlpBlockHash(arry.Hash) (types.IRlpBlock, error)

	NextHeader(uint64) (types.IHeader, error)
	NextBlock([]types.IMessage, uint64) (types.IBlock, error)
	Insert(types.IBlock) error
	Roll() error
	Vote(arry.Address) uint64
}
