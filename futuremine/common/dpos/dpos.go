package dpos

import (
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/futuremine/common/param"
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"github.com/Futuremine-chain/futuremine/types"
)

const (
	SuperCount = 3
)

type DPos struct {
	cycle     *Cycle
	confirmed uint64
}

func NewDPos(dPosStatus dpos.IDPosStatus) *DPos {
	return &DPos{cycle: &Cycle{DPosStatus: dPosStatus}}
}

func (d *DPos) GenesisBlock() types.IBlock {
	block := &fmctypes.Block{
		Header: fmctypes.NewHeader(arry.Hash{},
			arry.Hash{},
			arry.Hash{},
			arry.Hash{},
			arry.Hash{},
			0,
			config.Param.GenesisTime,
			arry.Address{},
		),
		Body: &fmctypes.Body{Messages: make([]types.IMessage, 0)},
	}
	for _, super := range genesisSuperList.Members {
		var peerId fmctypes.Peer
		copy(peerId[:], super.PeerId)
		msg := &fmctypes.Message{
			Header: &fmctypes.MsgHeader{
				Type:      fmctypes.Candidate,
				Nonce:     1,
				From:      super.Signer,
				Time:     config.Param.GenesisTime,
				Signature: &fmctypes.Signature{},
			},
			Body: &fmctypes.CandidateBody{Peer: peerId},
		}
		msg.SetHash()
		block.Body.Messages = append(block.Body.Messages, msg)
	}
	block.SetHash()
	return block
}

func (d *DPos) CheckTime(header types.IHeader, chain blockchain.IChain) error {
	preHeader, err := chain.GetBlockHash(header.GetPreHash())
	if err != nil {
		return err
	}

	if err := d.cycle.CheckCycle(chain, preHeader.GetTime(), header.GetTime()); err != Err_Elected {
		if err := d.cycle.Elect(header.GetTime(), preHeader.GetHash(), chain); err != nil {
			return err
		}
	}

	// Check if the time of block production is correct
	if err := d.checkTime(preHeader, header); err != nil {
		return err
	}
	return nil
}

func (d *DPos) CheckSigner(header types.IHeader, chain blockchain.IChain) error {
	// Find the block address at that time
	super, err := d.lookupSuper(header.GetTime())
	if err != nil {
		return err
	}
	if !super.IsEqual(header.GetSigner()) {
		return errors.New("it's not the miner's turn")
	}
	return nil
}

func (d *DPos) CheckHeader(header types.IHeader, parent types.IHeader, chain blockchain.IChain) error {
	// If the block time is in the future, it will fail
	if header.GetTime() > uint64(utils.NowUnix()) {
		return errors.New("block in the future")
	}
	// Verify whether it is the time point of block generation
	if err := d.checkTime(parent, header); err != nil {
		return errors.New("time check failed")
	}
	if header.GetSignature() == nil {
		return errors.New("no signature")
	}
	if parent.GetTime()+param.BlockInterval > header.GetTime() {
		return errors.New("invalid timestamp")
	}
	return nil
}

func (d *DPos) CheckSeal(header types.IHeader, parent types.IHeader, chain blockchain.IChain) error {
	// Verifying the genesis block is not supported
	if header.GetHeight() == 0 {
		return errors.New("unknown block")
	}
	if header.GetHeight() <= d.confirmed {
		return errors.New("height error")
	}
	lastCycleHeader, err := d.preCycleLastHash(header, chain)
	if err != nil {
		return err
	}
	// Verify the block node
	if err := d.CheckCreator(header, lastCycleHeader, chain); err != nil {
		return err
	}
	// Update the height of the confirmed block
	return d.updateConfirmed(chain)
}

func (d *DPos) CheckCreator(header types.IHeader, parent types.IHeader, chain blockchain.IChain) error {
	signer, err := d.setAndLookupSuper(header.GetTime(), parent, chain)
	if err != nil {
		return err
	}
	if err := d.checkSigner(signer, header); err != nil {
		return err
	}
	return nil
}

func (d *DPos) SuperIds() []string {
	return nil
}

func (d *DPos) Confirmed() uint64 {
	return d.confirmed
}

func (d *DPos) checkTime(lastHeader types.IHeader, header types.IHeader) error {
	nextTime := nextTime(header.GetTime())
	if lastHeader.GetTime() >= nextTime {
		return errors.New("create the future block")
	}
	if nextTime-header.GetTime() >= 1 {
		return fmt.Errorf("wait for last block arrived, next slot = %d, block time = %d ", nextTime, header.GetTime)
	}
	if header.GetTime() == nextTime {
		return nil
	}
	return fmt.Errorf("wait for last block arrived, next slot = %d, block time = %d ", nextTime, header.GetTime)
}

func (d *DPos) lookupSuper(now uint64) (arry.Address, error) {
	offset := now % param.CycleInterval
	if offset%param.BlockInterval != 0 {
		return arry.Address{}, errors.New("invalid time to mint the block")
	}
	offset /= param.BlockInterval
	supers, err := d.cycle.DPosStatus.CycleSupers(now / param.CycleInterval)
	if err != nil {
		return arry.Address{}, err
	}
	if len(supers.Candidates) == 0 {
		return arry.Address{}, errors.New("no super to be found in storage")
	}
	offset %= uint64(len(supers.Candidates))
	super := supers.Candidates[offset]
	return super.Signer, nil
}

func (d *DPos) setAndLookupSuper(now uint64, parent types.IHeader, chain blockchain.IChain) (arry.Address, error) {
	offset := now % param.CycleInterval
	if offset%param.BlockInterval != 0 {
		return arry.Address{}, errors.New("invalid time to mint the block")
	}
	offset /= param.BlockInterval
	supers, err := d.setSupers(now, parent, chain)
	if err != nil {
		return arry.Address{}, err
	}
	if len(supers) == 0 {
		return arry.Address{}, errors.New("no winner to be found in storage")
	}
	offset %= uint64(len(supers))
	winner := supers[offset]
	return winner.Signer, nil
}

func (d *DPos) setSupers(time uint64, parent types.IHeader, chain blockchain.IChain) ([]*fmctypes.Member, error) {
	cycle := time / param.CycleInterval
	supers, err := d.cycle.DPosStatus.CycleSupers(cycle)

	// If the election result of the current cycle does not
	// exist, the current cycle of elections is conducted
	if err != nil || supers == nil || !parent.GetHash().IsEqual(supers.PreHash) {
		if err := d.cycle.Elect(time, parent.GetHash(), chain); err != nil {
			return nil, err
		}
		if supers, err = d.cycle.DPosStatus.CycleSupers(cycle); err != nil {
			return nil, err
		}
	}
	return supers.Candidates, nil
}

// Get the hash of the last block of the previous cycle
// as the random number seed of the new cycle.
func (d *DPos) preCycleLastHash(current types.IHeader, chain blockchain.IChain) (types.IHeader, error) {
	preTermLastHash, err := chain.CycleLastHash(current.GetCycle() - 1)
	if err == nil {
		header, _ := chain.GetBlockHash(preTermLastHash)
		tHeader, _ := chain.GetHeaderHeight(header.GetHeight())
		if header.GetHeight() < current.GetHeight() && header.GetHash().IsEqual(tHeader.GetHash()) {
			return header, nil
		}
	}

	// If the last block header of the last cycle cannot be obtained directly
	// from the chain, then look forward from the current block
	genesis, err := chain.GetHeaderHeight(0)
	if err != nil {
		return nil, err
	}
	header, err := chain.GetHeaderHeight(1)
	if err != nil {
		return genesis, nil
	}

	if header.GetCycle() >= current.GetCycle() {
		return genesis, nil
	}
	height := current.GetHeight()
	for height > 0 {
		height--
		header, err := chain.GetHeaderHeight(height)
		if err != nil {
			continue
		}
		if header.GetCycle() < current.GetCycle() {
			return header, nil
		}
	}
	return nil, errors.New("not found")
}

func (d *DPos) checkSigner(super arry.Address, header types.IHeader) error {
	if !fmctypes.VerifySigner(config.Param.Name, super, header.GetSignature().PubicKey()) {
		return errors.New("not the signature of the address")
	}
	if !fmctypes.Verify(header.GetHash(), header.GetSignature()) {
		return errors.New("verify seal failed")
	}
	return nil
}

// updateConfirmedBlockHeader Update the final confirmation block
func (d *DPos) updateConfirmed(chain blockchain.IChain) error {
	if d.confirmed == 0 {
		height, err := d.cycle.DPosStatus.Confirmed()
		if err != nil {
			header, err := chain.GetHeaderHeight(0)
			if err != nil {
				return err
			}
			height = header.GetHeight()
		}
		d.confirmed = height
	}
	curHeader, err := chain.LastHeader()
	if err != nil {
		return err
	}
	// If there are already more than two-thirds of different nodes generating blocks,
	// it means that the blocks before these blocks have been confirmed

	cycle := uint64(0)
	superMap := make(map[string]int)
	for d.confirmed < curHeader.GetHeight() {
		curCycle := curHeader.GetTime() / param.CycleInterval
		if curCycle != cycle {
			cycle = curCycle
			superMap = make(map[string]int)
		}
		// fast return
		// if block number difference less consensusSize-witnessNum
		// there is no need to check block is confirmed

		count := superMap[curHeader.GetSigner().String()]
		superMap[curHeader.GetSigner().String()] = count + 1

		if len(superMap) >= param.DPosSize /*dpos.checkWinnerMapCount(winnerMap, 1)*/ {
			d.cycle.DPosStatus.SetConfirmed(curHeader.GetHeight())
			d.confirmed = curHeader.GetHeight()
			chain.SetConfirmed(curHeader.GetHeight())
			//log.Info("DPos set confirmed block header", "currentHeader", curHeader.Height)
			return nil
		}
		curHeader, err = chain.GetHeaderHash(curHeader.GetPreHash())
		if err != nil {
			return errors.New("nil block header returned")
		}
	}
	return nil
}

func nextTime(now uint64) uint64 {
	return (now + param.BlockInterval - 1) / param.BlockInterval * param.BlockInterval
}
