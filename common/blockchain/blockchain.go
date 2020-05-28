package blockchain

type IBlockChain interface {
	LastHeight() uint64
}
