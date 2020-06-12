package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	"github.com/Futuremine-chain/futuremine/types"
)

type Block struct {
	*Header
	*Body
}

func (b *Block) Hash() arry.Hash {
	return b.Header.Hash()
}

func (b *Block) MsgRoot() arry.Hash {
	return b.Header.MsgRoot()
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

func (b *Block) Add(message types.IMessage) {
	b.Body.Add(message)
}

func (b *Block) Msgs() []types.IMessage {
	return b.Body.Msgs()
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

func (b *Block) SetHash() {
	b.Header.SetHash()
}

func (b *Block) Sign(key *secp256k1.PrivateKey) error {
	return b.Header.Sign(key)
}

type Blocks []*Block

func (b Blocks) Blocks() []*Block {
	return b
}
