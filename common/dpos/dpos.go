package dpos

import (
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/types"
)

type IDPos interface {
	CheckTime(header types.IHeader, chain blockchain.IChain) error
	CheckSigner(header types.IHeader, chain blockchain.IChain) error
	CheckHeader(header types.IHeader, parent types.IHeader, chain blockchain.IChain) error
	CheckSeal(header types.IHeader, parent types.IHeader, chain blockchain.IChain) error
	SuperIds() []string
	Confirmed() uint64
}
