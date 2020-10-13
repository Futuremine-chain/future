package blockchain

import (
	"errors"
	"fmt"
	"github.com/Futuremine-chain/future/common/config"
	"github.com/Futuremine-chain/future/common/dpos"
	"github.com/Futuremine-chain/future/common/status"
	"github.com/Futuremine-chain/future/future/common/kit"
	"github.com/Futuremine-chain/future/future/db/chain_db"
	fmctypes "github.com/Futuremine-chain/future/future/types"
	servicesync "github.com/Futuremine-chain/future/service/sync"
	"github.com/Futuremine-chain/future/tools/arry"
	log "github.com/Futuremine-chain/future/tools/log/log15"
	"github.com/Futuremine-chain/future/types"
	"sync"
)

const chainDB = "chain_db"
const module = "module"

type FMCChain struct {
	mutex         sync.RWMutex
	status        status.IStatus
	db            IChainDB
	dPos          dpos.IDPos
	actRoot       arry.Hash
	dPosRoot      arry.Hash
	tokenRoot     arry.Hash
	lastHeight    uint64
	confirmed     uint64
	poolDeleteMsg func(message types.IMessage)
}

func NewFMCChain(status status.IStatus, dPos dpos.IDPos) (*FMCChain, error) {
	var err error
	fmc := &FMCChain{status: status, dPos: dPos}
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

func (b *FMCChain) NextHeader(time uint64) (types.IHeader, error) {
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

func (b *FMCChain) NextBlock(msgs []types.IMessage, blockTime uint64) (types.IBlock, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	height := b.lastHeight + 1
	coinbase := kit.CalCoinBase(config.Param.Name, height)

	coinBase := &fmctypes.Message{
		Header: &fmctypes.MsgHeader{
			Type: fmctypes.Transaction,
			From: fmctypes.CoinBase,
			Time: blockTime,
		},
		Body: &fmctypes.TransactionBody{
			TokenAddress: config.Param.TokenParam.MainToken,
			Receiver:     config.Param.IPrivate.Address(),
			Amount:       coinbase + fmctypes.CalculateFee(msgs),
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
		fmctypes.MsgRoot(fmcMsgs),
		b.actRoot,
		b.dPosRoot,
		b.tokenRoot,
		height,
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
	b.status.SetConfirmed(confirmed)
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

func (b *FMCChain) CycleLastHash(cycle uint64) (arry.Hash, error) {
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

func (b *FMCChain) GetMessage(hash arry.Hash) (types.IMessage, error) {
	rlpTx, err := b.db.GetMessage(hash)
	if err != nil {
		return nil, err
	}
	return rlpTx.ToMessage(), nil
}

func (b *FMCChain) GetMessageIndex(hash arry.Hash) (types.IMessageIndex, error) {
	msgIndex, err := b.db.GetMsgIndex(hash)
	if err != nil {
		return nil, err
	}
	if msgIndex.Height > b.LastHeight() {
		return nil, errors.New("not exist")
	}
	return msgIndex, nil
}

func (b *FMCChain) Insert(block types.IBlock) error {
	if err := b.checkBlock(block); err != nil {
		return err
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.lastHeight >= block.GetHeight() {
		return errors.New("wrong block height")
	}
	if err := b.status.Change(block.BlockBody().MsgList(), block); err != nil {
		return err
	}
	msgs := block.BlockBody().MsgList()
	for _, msg := range msgs {
		if b.poolDeleteMsg != nil {
			b.poolDeleteMsg(msg)
		} else {
			log.Error("Need to register message pool delete function", "module", module)
		}
	}
	b.saveBlock(block)
	return nil
}

func (b *FMCChain) saveBlock(block types.IBlock) {
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
	/*log.Info("Save block", "module", "module",
	"height", block.GetHeight(),
	"hash", block.GetHash().String(),
	"actroot", block.GetActRoot().String(),
	"tokenroot", block.GetTokenRoot().String(),
	"dposroot", block.GetDPosRoot().String(),
	"signer", block.GetSigner().String(),
	"msgcount", len(block.BlockBody().MsgList()),
	"time", block.GetTime(),
	"cycle", block.GetCycle())*/
}

func (b *FMCChain) saveGenesisBlock(block types.IBlock) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if err := b.verifyGenesis(block); err != nil {
		return err
	}

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
	return nil
}

func (b *FMCChain) verifyGenesis(block types.IBlock) error {
	var sumCoins uint64
	for _, tx := range block.BlockBody().MsgList() {
		sumCoins += tx.MsgBody().MsgAmount()
	}
	if sumCoins != config.Param.PreCirculation {
		return fmt.Errorf("wrong genesis coins")
	}
	return nil
}

func (b *FMCChain) checkBlock(block types.IBlock) error {
	lastHeight := b.LastHeight()

	if block.GetHeight() == lastHeight {
		lastHeader, err := b.GetHeaderHeight(lastHeight)
		if err == nil && lastHeader.GetHash().IsEqual(block.GetHash()) {
			return servicesync.Err_RepeatBlock
		}
	}

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
		return errors.New("wrong dpos root")
	}
	if !block.GetTokenRoot().IsEqual(b.tokenRoot) {
		log.Warn("the token status root hash verification failed", "module", module,
			"height", block.GetHeight(), "tokenroot", block.GetTokenRoot().String())
		return errors.New("wrong token root")
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

func (b *FMCChain) checkMsgs(msgs []types.IMessage, height uint64) error {
	address := make(map[string]int)
	for i, msg := range msgs {
		if msg.IsCoinBase() {
			if err := b.checkCoinBase(msg, fmctypes.CalculateFee(msgs), height); err != nil {
				return err
			}
		} else {
			if err := b.checkMsg(msg); err != nil {
				return err
			}
		}
		from := msg.From().String()
		if lastIndex, ok := address[from]; !ok {
			address[from] = i
		} else {
			log.Warn("Repeat address block", "module", module,
				"preMsg", msgs[lastIndex],
				"curMsg", msg)
			return errors.New("one address in a block can only send one transaction")
		}
	}
	return nil
}

func (b *FMCChain) checkCoinBase(coinBase types.IMessage, fee, height uint64) error {
	msg, ok := coinBase.(*fmctypes.Message)
	if !ok {
		return errors.New("wrong message type")
	}
	coinbase := kit.CalCoinBase(config.Param.Name, height)

	if err := msg.CheckCoinBase(fee, coinbase); err != nil {
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

func (b *FMCChain) Confirmed() uint64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.confirmed
}

func (b *FMCChain) Roll() error {
	var curHeight uint64
	confirmed := b.Confirmed()
	if confirmed != 0 {
		curHeight = confirmed
	}
	return b.RollbackTo(curHeight)
}

func (b *FMCChain) RollbackTo(height uint64) error {
	confirmedHeight := b.confirmed
	if height > confirmedHeight && height != 0 {
		err := fmt.Sprintf("the height of the roolback must be less than or equal to %d and greater than %d", confirmedHeight, 0)
		log.Error("Roll back to block height", "height", height, "error", err)
		return errors.New(err)
	}

	var curBlockHeight, nextBlockHeight uint64
	curActRoot := arry.Hash{}
	curTokenRoot := arry.Hash{}
	curDPosRoot := arry.Hash{}

	nextBlockHeight = height + 1
	curBlockHeight = height

	// set new confirmed height and header
	hisConfirmedHeight, err := b.db.GetConfirmedHeight(curBlockHeight)
	if err != nil {
		log.Error("Fall back to block height", "height", height, "error", "can not find history confirmed height")
		return fmt.Errorf("fall back to block height %d failed! Can not find history confirmed height", height)
	}
	b.dPos.SetConfirmed(hisConfirmedHeight)

	log.Warn("Fall back to block height", "height", height)
	header, err := b.GetHeaderHeight(nextBlockHeight)
	if err != nil {
		log.Error("Fall back to block height", "height", height, "error", "can not find block")
		return fmt.Errorf("fall back to block height %d failed! Can not find block %d", height, nextBlockHeight)
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.confirmed = hisConfirmedHeight
	b.status.SetConfirmed(hisConfirmedHeight)

	// fall back to pre state root
	curActRoot = header.GetActRoot()
	curTokenRoot = header.GetTokenRoot()
	curDPosRoot = header.GetDPosRoot()
	err = b.status.InitRoots(curActRoot, curDPosRoot, curTokenRoot)
	if err != nil {
		log.Error("Fall back to block height", "height", height, "error", "init state trie failed")
		return fmt.Errorf("fall back to block height %d failed! nit state trie failed", height)
	}
	b.actRoot = curActRoot
	b.tokenRoot = curTokenRoot
	b.dPosRoot = curDPosRoot
	b.db.SaveActRoot(b.actRoot)
	b.db.SaveTokenRoot(b.tokenRoot)
	b.db.SaveDPosRoot(b.dPosRoot)

	b.lastHeight = curBlockHeight
	b.db.SaveLastHeight(curBlockHeight)
	return nil
}

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

func (b *FMCChain) RegisterMsgPoolDeleteFunc(fun func(message types.IMessage)) {
	b.poolDeleteMsg = fun
}
