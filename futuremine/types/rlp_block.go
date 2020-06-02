package types

import "github.com/Futuremine-chain/futuremine/types"

type RlpBlock struct {
}

func (r *RlpBlock) ToBlock() types.IBlock {
	return nil
}

func DecodeRlpBlock([]byte) (*RlpBlock, error) {
	return nil, nil
}

type RlpBlocks []RlpBlock

func (r *RlpBlocks) ToBlocks() []types.IBlock {
	return nil
}

func DecodeRlpBlocks([]byte) (*RlpBlocks, error) {
	return nil, nil
}
