package blockchain

import (
	"github.com/Futuremine-chain/futuremine/futuremine/common/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type BlockChain struct {
}

func NewBlockChain() *BlockChain {
	return &BlockChain{}
}

func (b *BlockChain) LastHeight() uint64 {
	return 0
}

func (b *BlockChain) NextBlock(txs types.ITransactions) types.IBlock {
	return nil
}

func (b *BlockChain) LastConfirmed() uint64                              { return 0 }
func (b *BlockChain) GetBlockHeight(uint64) (types.IBlock, error)        { return nil, nil }
func (b *BlockChain) GetBlockHash(arry.Hash) (types.IBlock, error)       { return nil, nil }
func (b *BlockChain) GetHeaderHeight(uint64) (types.IHeader, error)      { return nil, nil }
func (b *BlockChain) GetHeaderHash(arry.Hash) (types.IHeader, error)     { return nil, nil }
func (b *BlockChain) GetRlpBlockHeight(uint64) (types.IRlpBlock, error)  { return nil, nil }
func (b *BlockChain) GetRlpBlockHash(arry.Hash) (types.IRlpBlock, error) { return nil, nil }
func (b *BlockChain) Insert(block types.IBlock) error                    { return nil }
func (b *BlockChain) Roll() error                                        { return nil }
