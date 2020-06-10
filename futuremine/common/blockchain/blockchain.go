package blockchain

import (
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/common/status"
	"github.com/Futuremine-chain/futuremine/futuremine/db/chain_db"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
	"sync"
)

const chainDB = "chain_db"

type FMCChain struct {
	mutex      sync.RWMutex
	status     status.IStatus
	db         IChainDB
	dPos       dpos.IDPos
	actRoot    arry.Hash
	dPosRoot   arry.Hash
	tokenRoot  arry.Hash
	lastHeight uint64
	confirmed  uint64
}

func NewFMCChain(status status.IStatus, dPos dpos.IDPos) (*FMCChain, error) {
	var err error
	fmc := &FMCChain{status: status, dPos: dPos}
	fmc.db, err = chain_db.Open(config.App.Setting().Data + "/" + chainDB)
	if err != nil {
		return nil, err
	}
	// Read the status tree root hash
	if fmc.actRoot, err = fmc.db.ActRoot(); err != nil {
		return nil, err
	}
	if fmc.dPosRoot, err = fmc.db.DPosRoot(); err != nil {
		return nil, err
	}
	if fmc.tokenRoot, err = fmc.db.TokenRoot(); err != nil {
		return nil, err
	}

	// Initializes the state root hash
	if err := fmc.status.InitRoots(fmc.actRoot, fmc.dPosRoot, fmc.tokenRoot); err != nil {
		return nil, err
	}

	// Initialize chain height
	fmc.lastHeight = fmc.db.LastHeight()
	fmc.UpdateConfirmed(fmc.dPos.Confirmed())
	return fmc, nil
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

func (b *FMCChain) UpdateConfirmed(height uint64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.confirmed = height
	b.status.SetConfirmed(height)
}
