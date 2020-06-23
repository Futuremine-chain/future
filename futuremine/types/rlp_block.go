package types

import (
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/types"
)

type RlpBlock struct {
	RlpHeader *Header
	RlpBody   *RlpBody
}

func (r *RlpBlock) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(r)
	return bytes
}

func (r *RlpBlock) ToBlock() types.IBlock {
	return &Block{
		Header: r.RlpHeader,
		Body:   r.RlpBody.ToBody(),
	}
	return nil
}

func DecodeRlpBlock(bytes []byte) (*RlpBlock, error) {
	var rlpBlock *RlpBlock
	if err := rlp.DecodeBytes(bytes, &rlpBlock); err != nil {
		return nil, err
	}
	return rlpBlock, nil
}

type RlpBlocks []types.IRlpBlock

func (r *RlpBlocks) ToBlocks() []types.IBlock {
	rs := make([]types.IBlock, len(*r))
	for i, rlpBlock := range *r {
		rs[i] = rlpBlock.ToBlock()
	}
	return rs
}

func (r *RlpBlocks) Add(block types.IRlpBlock) {
}

func (r *RlpBlocks) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(r)
}

func DecodeRlpBlocks(bytes []byte) (*RlpBlocks, error) {
	var rlpBlocks *RlpBlocks
	err := rlp.DecodeBytes(bytes, &rlpBlocks)
	if err != nil {
		return nil, err
	}
	return rlpBlocks, nil
}
