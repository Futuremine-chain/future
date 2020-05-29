package blockchain

import (
	"github.com/Futuremine-chain/futuremine/types"
)

type IBlockChain interface {
	LastHeight() uint64
	NextBlock(types.ITransactions) types.IBlock
}
