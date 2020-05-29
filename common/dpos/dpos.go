package dpos

import (
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/types"
)

type IDPos interface {
	CheckSigner(chain blockchain.IBlockChain, header types.IHeader) error
}
