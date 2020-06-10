package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type Block struct {
	*Header
	*Body
}

func (b *Block) Hash() arry.Hash {
	return b.Header.Hash()
}

func (b *Block) TxRoot() arry.Hash {
	return b.Header.TxRoot()
}

func (b *Block) Signer() arry.Address {
	return b.Header.Signer()
}

func (b *Block) Height() uint64 {
	return b.Header.height
}

func (b *Block) Time() int64 {
	return b.Header.time
}

func (b *Block) Add(transaction types.ITransaction) {
	b.Body.Add(transaction)
}

func (b *Block) Txs() []types.ITransaction {
	return b.Body.Txs()
}

func (b *Block) ToRlpHeader() types.IRlpHeader {
	panic("implement me")
}

func (b *Block) BlockHeader() types.IHeader {
	return b.Header
}

func (b *Block) ToRlpBlock() types.IRlpBlock {
	panic("implement me")
}

type Blocks []*Block

func (b Blocks) Blocks() []*Block {
	return b
}
