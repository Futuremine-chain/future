package blockchain

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/common/status"
	"github.com/Futuremine-chain/futuremine/futuremine/db/chain_db"
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
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
		return nil, fmt.Errorf("failed to open chain db, %s", err.Error())
	}
	// Read the status tree root hash
	fmc.actRoot, _ = fmc.db.ActRoot()
	fmc.dPosRoot, _ = fmc.db.DPosRoot()
	fmc.tokenRoot, _ = fmc.db.TokenRoot()

	// Initializes the state root hash
	if err := fmc.status.InitRoots(fmc.actRoot, fmc.dPosRoot, fmc.tokenRoot); err != nil {
		return nil, fmt.Errorf("failed to init status root, %s", err.Error())
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

func (b *FMCChain) LastConfirmed() uint64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.confirmed
}

func (b *FMCChain) GetBlockHeight(height uint64) (types.IBlock, error) {
	header, err := b.getHeaderHeight(height)
	if err != nil {
		return nil, err
	}
	txs, err := b.db.GetTransactions(header.TxRoot())
	if err != nil {
		return nil, err
	}
	rlpBody := &fmctypes.RlpBody{txs}
	block := &fmctypes.Block{header, rlpBody.ToBody()}
	return block, nil
}

func (b *FMCChain) GetBlockHash(hash arry.Hash) (types.IBlock, error) {
	header, err := b.getHeaderHash(hash)
	if err != nil {
		return nil, err
	}
	txs, err := b.db.GetTransactions(header.TxRoot())
	if err != nil {
		return nil, err
	}
	rlpBody := &fmctypes.RlpBody{txs}
	block := &fmctypes.Block{header, rlpBody.ToBody()}
	return block, nil
}

func (b *FMCChain) GetHeaderHeight(height uint64) (types.IHeader, error) {
	return b.getHeaderHeight(height)
}

func (b *FMCChain) getHeaderHeight(height uint64) (*fmctypes.Header, error) {
	if height > b.LastHeight() {
		return nil, fmt.Errorf("%d block header is not exist", height)
	}
	return b.db.GetHeaderHeight(height)
}

func (b *FMCChain) GetHeaderHash(hash arry.Hash) (types.IHeader, error) {
	return b.getHeaderHash(hash)
}

func (b *FMCChain) getHeaderHash(hash arry.Hash) (*fmctypes.Header, error) {
	return b.db.GetHeaderHash(hash)
}

func (b *FMCChain) GetRlpBlockHeight(height uint64) (types.IRlpBlock, error) {
	header, err := b.db.GetHeaderHeight(height)
	if err != nil {
		return nil, err
	}
	txs, err := b.db.GetTransactions(header.TxRoot())
	if err != nil {
		return nil, err
	}
	rlpBody := &fmctypes.RlpBody{txs}
	rlpHeader := header.ToRlpHeader().(*fmctypes.RlpHeader)
	block := &fmctypes.RlpBlock{rlpHeader, rlpBody}
	return block, nil
}

func (b *FMCChain) GetRlpBlockHash(hash arry.Hash) (types.IRlpBlock, error) {
	header, err := b.db.GetHeaderHash(hash)
	if err != nil {
		return nil, err
	}
	txs, err := b.db.GetTransactions(header.TxRoot())
	if err != nil {
		return nil, err
	}
	rlpBody := &fmctypes.RlpBody{txs}
	rlpHeader := header.ToRlpHeader().(*fmctypes.RlpHeader)
	block := &fmctypes.RlpBlock{rlpHeader, rlpBody}
	return block, nil
}

func (b *FMCChain) Insert(block types.IBlock) error { return nil }
func (b *FMCChain) Roll() error                     { return nil }

func (b *FMCChain) UpdateConfirmed(height uint64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.confirmed = height
	b.status.SetConfirmed(height)
}
