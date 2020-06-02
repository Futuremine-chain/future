package blockchain

import (
	"github.com/Futuremine-chain/futuremine/futuremine/common/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IBlockChain interface {
	LastHeight() uint64
	LastConfirmed() uint64
	GetBlockHeight(uint64) (types.IBlock, error)
	GetBlockHash(arry.Hash) (types.IBlock, error)
	GetHeaderHeight(uint64) (types.IHeader, error)
	GetHeaderHash(arry.Hash) (types.IHeader, error)
	NextBlock(types.ITransactions) types.IBlock
	Insert(block types.IBlock) error
	Roll() error
}
