package dpos

import (
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/types"
)

type IDPos interface {
	CheckTime(header types.IHeader, chain blockchain.IChain) error
	CheckSigner(chain blockchain.IChain, header types.IHeader) error
	SuperIds() []string
	Confirmed() uint64
}
