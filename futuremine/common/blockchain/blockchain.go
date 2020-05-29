package blockchain

import "github.com/Futuremine-chain/futuremine/types"

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
