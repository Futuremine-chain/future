package types

import (
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/types"
)

type RlpBlock struct {
	*RlpHeader
	*RlpBody
}

func (r *RlpBlock) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(r)
	return bytes
}

func (r *RlpBlock) ToBlock() types.IBlock {
	return nil
}

func DecodeRlpBlock([]byte) (*RlpBlock, error) {
	return nil, nil
}

type RlpBlocks []types.IRlpBlock

func (r *RlpBlocks) ToBlocks() []types.IBlock {
	return nil
}

func (r *RlpBlocks) Add(block types.IRlpBlock) {
}

func (r *RlpBlocks) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(r)
}

func DecodeRlpBlocks([]byte) (*RlpBlocks, error) {
	return nil, nil
}
