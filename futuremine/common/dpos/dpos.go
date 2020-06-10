package dpos

import (
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/types"
)

type DPos struct {
}

func NewDPos() *DPos {
	return &DPos{}
}

func (dpos *DPos) CheckSigner(chain blockchain.IBlockChain, header types.IHeader) error {
	return nil
}

func (dpos *DPos) SuperIds() []string {
	return nil
}

func (dpos *DPos) Confirmed() uint64 {
	return 0
}
