package dpos

import (
	"github.com/Futuremine-chain/future/common/blockchain"
	"github.com/Futuremine-chain/future/types"
)

type IDPos interface {
	GenesisBlock() types.IBlock
	CheckTime(header types.IHeader, chain blockchain.IChain) error
	CheckSigner(header types.IHeader, chain blockchain.IChain) error
	CheckHeader(header types.IHeader, parent types.IHeader, chain blockchain.IChain) error
	CheckSeal(header types.IHeader, parent types.IHeader, chain blockchain.IChain) error
	Confirmed() uint64
	SetConfirmed(uint64)
}
