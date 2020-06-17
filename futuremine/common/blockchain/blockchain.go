package blockchain

import (
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/common/status"
	"github.com/Futuremine-chain/futuremine/futuremine/db/chain_db"
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/service/pool"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/Futuremine-chain/futuremine/types"
	"sync"
	"time"
)

const chainDB = "chain_db"
const module = "module"

type FMCChain struct {
	mutex      sync.RWMutex
	status     status.IStatus
	db         IChainDB
	msgPool    *pool.Pool
	dPos       dpos.IDPos
	actRoot    arry.Hash
	dPosRoot   arry.Hash
	tokenRoot  arry.Hash
	lastHeight uint64
	confirmed  uint64
}

func NewFMCChain(status status.IStatus, dPos dpos.IDPos, msgPool *pool.Pool) (*FMCChain, error) {
	var err error
	fmc := &FMCChain{status: status, dPos: dPos, msgPool: msgPool}
	fmc.db, err = chain_db.Open(config.Param.Data + "/" + chainDB)
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
	if fmc.lastHeight, err = fmc.db.LastHeight(); err != nil {
		fmc.saveGenesisBlock(fmc.dPos.GenesisBlock())
	}
	fmc.UpdateConfirmed(fmc.dPos.Confirmed())
	return fmc, nil
}

func (b *FMCChain) LastHeight() uint64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.lastHeight
}

func (b *FMCChain) NextHeader(time int64) (types.IHeader, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	preHeader, err := b.GetHeaderHeight(b.lastHeight)
	if err != nil {
		return nil, err
	}
	// Build block header
	header := fmctypes.NewHeader(
		preHeader.GetHash(),
		arry.Hash{},
		b.actRoot,
		b.dPosRoot,
		b.tokenRoot,
		b.lastHeight+1,
		time,
		config.Param.IPrivate.Address(),
	)

	return header, nil
}

func (b *FMCChain) NextBlock(msgs []types.IMessage, blockTime int64) (types.IBlock, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	coinBase := &fmctypes.Message{
		Header: &fmctypes.MsgHeader{
			Type: fmctypes.Transaction,
			From: fmctypes.CoinBase,
			Time: time.Unix(blockTime, 0),
		},
		Body: &fmctypes.TransactionBody{
			TokenAddress: config.Param.TokenParam.MainToken,
			Receiver:     config.Param.IPrivate.Address(),
			Amount:       config.Param.TokenParam.Proportion + fmctypes.CalculateFee(msgs),
		},
	}
	coinBase.SetHash()
	fmcMsgs := msgs
	fmcMsgs = append(fmcMsgs, coinBase)
	lastHeader, err := b.GetHeaderHeight(b.lastHeight)
	if err != nil {
		return nil, err
	}
	// Build block header
	header := fmctypes.NewHeader(
		lastHeader.GetHash(),
		fmctypes.MsgRoot(msgs),
		b.actRoot,
		b.dPosRoot,
		b.tokenRoot,
		b.lastHeight+1,
		blockTime,
		config.Param.IPrivate.Address(),
	)
	body := &fmctypes.Body{fmcMsgs}
	newBlock := &fmctypes.Block{
		Header: header,
		Body:   body,
	}
	newBlock.SetHash()
	if err := newBlock.Sign(config.Param.IPrivate.PrivateKey()); err != nil {
		return nil, err
	}
	return newBlock, nil
}

func (b *FMCChain) LastConfirmed() uint64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.confirmed
}

func (b *FMCChain) SetConfirmed(confirmed uint64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.confirmed = confirmed
}

func (b *FMCChain) LastHeader() (types.IHeader, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.db.GetHeaderHeight(b.lastHeight)
}

func (b *FMCChain) GetBlockHeight(height uint64) (types.IBlock, error) {
	header, err := b.getHeaderHeight(height)
	if err != nil {
		return nil, err
	}
	txs, err := b.db.GetMessages(header.MsgRoot)
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
	txs, err := b.db.GetMessages(header.MsgRoot)
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

func (b *FMCChain) CycleLastHash(cycle int64) (arry.Hash, error) {
	return b.db.CycleLastHash(cycle)
}

func (b *FMCChain) GetRlpBlockHeight(height uint64) (types.IRlpBlock, error) {
	header, err := b.db.GetHeaderHeight(height)
	if err != nil {
		return nil, err
	}
	txs, err := b.db.GetMessages(header.MsgRoot)
	if err != nil {
		return nil, err
	}
	rlpBody := &fmctypes.RlpBody{txs}
	rlpHeader := header.ToRlpHeader().(*fmctypes.Header)
	block := &fmctypes.RlpBlock{rlpHeader, rlpBody}
	return block, nil
}

func (b *FMCChain) GetRlpBlockHash(hash arry.Hash) (types.IRlpBlock, error) {
	header, err := b.db.GetHeaderHash(hash)
	if err != nil {
		return nil, err
	}
	txs, err := b.db.GetMessages(header.MsgRoot)
	if err != nil {
		return nil, err
	}
	rlpBody := &fmctypes.RlpBody{txs}
	rlpHeader := header.ToRlpHeader().(*fmctypes.Header)
	block := &fmctypes.RlpBlock{rlpHeader, rlpBody}
	return block, nil
}

func (b *FMCChain) Insert(block types.IBlock) error {
	if err := b.checkBlock(block); err != nil {
		return err
	}
	if err := b.status.Change(block.BlockBody().MsgList(), block); err != nil {
		return err
	}
	b.saveBlock(block)
	return nil
}

func (b *FMCChain) saveBlock(block types.IBlock) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	bk := block.(*fmctypes.Block)
	rlpBlock := bk.ToRlpBlock().(*fmctypes.RlpBlock)
	b.db.SaveHeader(bk.Header)
	b.db.SaveMessages(block.GetMsgRoot(), rlpBlock.RlpBody.MsgList())
	b.db.SaveMsgIndex(bk.GetMsgIndexs())
	b.db.SaveHeightHash(block.GetHeight(), block.GetHash())
	b.db.SaveConfirmedHeight(block.GetHeight(), b.confirmed)
	b.db.SaveCycleLastHash(block.GetCycle(), block.GetHash())
	b.actRoot, b.tokenRoot, b.dPosRoot, _ = b.status.Commit()
	b.db.SaveActRoot(b.actRoot)
	b.db.SaveDPosRoot(b.dPosRoot)
	b.db.SaveTokenRoot(b.tokenRoot)

	b.lastHeight = block.GetHeight()
	b.db.SaveLastHeight(b.lastHeight)

	log.Info("Save block", "module", "module",
		"height", block.GetHeight(),
		"hash", block.GetHash().String(),
		"actroot", block.GetActRoot().String(),
		"tokenroot", block.GetTokenRoot().String(),
		"dposroot", block.GetDPosRoot().String(),
		"signer", block.GetSigner().String(),
		"msgcount", len(block.BlockBody().MsgList()),
		"time", block.GetTime(),
		"cycle", block.GetCycle())
}

func (b *FMCChain) saveGenesisBlock(block types.IBlock) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.status.Change(block.BlockBody().MsgList(), block)
	bk := block.(*fmctypes.Block)
	rlpBlock := bk.ToRlpBlock().(*fmctypes.RlpBlock)
	b.db.SaveHeader(bk.Header)
	b.db.SaveMessages(block.GetMsgRoot(), rlpBlock.RlpBody.MsgList())
	b.db.SaveMsgIndex(bk.GetMsgIndexs())
	b.db.SaveHeightHash(block.GetHeight(), block.GetHash())
	b.lastHeight = block.GetHeight()
	b.db.SaveConfirmedHeight(block.GetHeight(), b.confirmed)
	b.status.SetConfirmed(0)
	b.actRoot, b.tokenRoot, b.dPosRoot, _ = b.status.Commit()
	b.db.SaveActRoot(b.actRoot)
	b.db.SaveDPosRoot(b.dPosRoot)
	b.db.SaveTokenRoot(b.tokenRoot)
	b.db.SaveLastHeight(b.lastHeight)

	log.Info("Save block", "module", "module",
		"height", block.GetHeight(),
		"hash", block.GetHash().String(),
		"actroot", block.GetActRoot().String(),
		"tokenroot", block.GetTokenRoot().String(),
		"dposroot", block.GetDPosRoot().String(),
		"signer", block.GetSigner().String(),
		"msgcount", len(block.BlockBody().MsgList()),
		"time", block.GetTime(),
		"cycle", block.GetCycle())
}

func (b *FMCChain) checkBlock(block types.IBlock) error {
	lastHeight := b.LastHeight()

	if lastHeight != block.GetHeight()-1 {
		return fmt.Errorf("last height is %d, the current block height is %d", lastHeight, block.GetHeight())
	}

	if !block.CheckMsgRoot() {
		log.Warn("the message root hash verification failed", "module", module,
			"height", block.GetHeight(), "msgroot", block.GetMsgRoot().String())
		return errors.New("the message root hash verification failed")
	}
	if !block.GetActRoot().IsEqual(b.actRoot) {
		log.Warn("the account status root hash verification failed", "module", module,
			"height", block.GetHeight(), "actroot", block.GetActRoot().String())
		return errors.New("the account status root hash verification failed")
	}
	if !block.GetDPosRoot().IsEqual(b.dPosRoot) {
		log.Warn("the dpos status root hash verification failed", "module", module,
			"height", block.GetHeight(), "dposroot", block.GetDPosRoot().String())
		return errors.New("wrong contract root")
	}
	if !block.GetTokenRoot().IsEqual(b.tokenRoot) {
		log.Warn("the token status root hash verification failed", "module", module,
			"height", block.GetHeight(), "tokenroot", block.GetTokenRoot().String())
		return errors.New("wrong consensus root")
	}
	preHeader, err := b.GetHeaderHash(block.GetPreHash())
	if err != nil {
		return fmt.Errorf("no previous block %s found", block.GetPreHash().String())
	}

	if err := b.dPos.CheckHeader(block.BlockHeader(), preHeader, b); err != nil {
		return err
	}
	if err := b.dPos.CheckSeal(block.BlockHeader(), preHeader, b); err != nil {
		return err
	}
	if err := b.checkMsgs(block.BlockBody().MsgList(), block.GetHeight()); err != nil {
		return err
	}
	return nil
}

func (b *FMCChain) checkMsgs(msgs []types.IMessage, blockHeight uint64) error {
	address := make(map[string]bool)
	for _, msg := range msgs {
		if msg.IsCoinBase() {
			if err := b.checkCoinBase(msg, fmctypes.CalculateFee(msgs)); err != nil {
				return err
			}
		} else {
			if err := b.checkMsg(msg); err != nil {
				return err
			}
		}
		from := msg.From().String()
		if _, ok := address[from]; !ok {
			address[from] = true
		} else {
			return errors.New("one address in a block can only send one transaction")
		}
	}
	b.msgPool.Delete(msgs)
	return nil
}

func (b *FMCChain) checkCoinBase(coinBase types.IMessage, fee uint64) error {
	msg, ok := coinBase.(*fmctypes.Message)
	if !ok {
		return errors.New("wrong message type")
	}

	if err := msg.CheckCoinBase(fee); err != nil {
		return err
	}
	return nil
}

func (b *FMCChain) checkMsg(msg types.IMessage) error {
	msg, ok := msg.(*fmctypes.Message)
	if !ok {
		return errors.New("wrong message type")
	}

	if err := msg.Check(); err != nil {
		return err
	}

	if err := b.status.CheckMsg(msg, true); err != nil {
		return err
	}
	return nil
}

func (b *FMCChain) Roll() error { return nil }

func (b *FMCChain) UpdateConfirmed(height uint64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.confirmed = height
	b.status.SetConfirmed(height)
}

func (b *FMCChain) Vote(address arry.Address) uint64 {
	var vote uint64
	act := b.status.Account(address)
	vote += act.GetBalance(config.Param.MainToken)
	return vote
}
