package dpos

import (
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/types"
)

type IDPos interface {
	CheckTime(time int64) error
	CheckSigner(chain blockchain.IBlockChain, header types.IHeader) error
	SuperIds() []string
	Confirmed() uint64
}
