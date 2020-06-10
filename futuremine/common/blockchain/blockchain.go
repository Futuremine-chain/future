package blockchain

import (
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/status"
	"github.com/Futuremine-chain/futuremine/futuremine/db/chain"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
	"sync"
)

const chain_db = "chain_db"

type FMCChain struct {
	mutex      sync.RWMutex
	status     status.IStatus
	db         IChainDB
	actRoot    arry.Hash
	dPosRoot   arry.Hash
	tokenRoot  arry.Hash
	lastHeight uint64
}

func NewFMCChain(status status.IStatus) (*FMCChain, error) {
	db, err := chain.OpenChainDB(config.App.Setting().Data + "/" + chain_db)
	if err != nil {
		return nil, err
	}
	return &FMCChain{db: db, status: status}, nil
}

func (b *FMCChain) LastHeight() uint64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.lastHeight
}

func (b *FMCChain) NextBlock(txs types.ITransactions) types.IBlock {
	return nil
}

func (b *FMCChain) LastConfirmed() uint64                              { return 0 }
func (b *FMCChain) GetBlockHeight(uint64) (types.IBlock, error)        { return nil, nil }
func (b *FMCChain) GetBlockHash(arry.Hash) (types.IBlock, error)       { return nil, nil }
func (b *FMCChain) GetHeaderHeight(uint64) (types.IHeader, error)      { return nil, nil }
func (b *FMCChain) GetHeaderHash(arry.Hash) (types.IHeader, error)     { return nil, nil }
func (b *FMCChain) GetRlpBlockHeight(uint64) (types.IRlpBlock, error)  { return nil, nil }
func (b *FMCChain) GetRlpBlockHash(arry.Hash) (types.IRlpBlock, error) { return nil, nil }
func (b *FMCChain) Insert(block types.IBlock) error                    { return nil }
func (b *FMCChain) Roll() error                                        { return nil }
