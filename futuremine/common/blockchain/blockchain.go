package blockchain

type BlockChain struct {
}

func NewBlockChain() *BlockChain {
	return &BlockChain{}
}

func (b *BlockChain) LastHeight() uint64 {
	return 0
}
