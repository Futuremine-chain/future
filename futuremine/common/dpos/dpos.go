package dpos

import (
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/futuremine/common/param"
	"github.com/Futuremine-chain/futuremine/types"
)

const (
	SuperCount = 3
)

type DPos struct {
	cycle *Cycle
}

func NewDPos(dPosStatus dpos.IDPosStatus) *DPos {
	return &DPos{cycle: &Cycle{DPosStatus: dPosStatus}}
}

func (d *DPos) CheckTime(header types.IHeader, chain blockchain.IChain) error {
	preHeader, err := chain.GetBlockHash(header.PreHash())
	if err != nil {
		return err
	}

	if err := d.cycle.CheckCycle(chain, preHeader.Time(), header.Time()); err != Err_Elected {
		if err := d.cycle.Elect(header.Time(), preHeader.Hash(), chain); err != nil {
			return err
		}
	}

	// Check if the time of block production is correct
	if err := d.checkTime(preHeader, header); err != nil {
		return err
	}
	return nil
}

func (d *DPos) CheckSigner(chain blockchain.IChain, header types.IHeader) error {
	return nil
}

func (d *DPos) SuperIds() []string {
	return nil
}

func (d *DPos) Confirmed() uint64 {
	return 0
}

func (d *DPos) checkTime(lastHeader types.IHeader, header types.IHeader) error {
	nextTime := nextTime(header.Time())
	if lastHeader.Time() >= nextTime {
		return errors.New("create the future block")
	}
	if nextTime-header.Time() >= 1 {
		return fmt.Errorf("wait for last block arrived, next slot = %d, block time = %d ", nextTime, header.Time)
	}
	if header.Time() == nextTime {
		return nil
	}
	return fmt.Errorf("wait for last block arrived, next slot = %d, block time = %d ", nextTime, header.Time)
}

func nextTime(now int64) int64 {
	return (now + param.BlockInterval - 1) / param.BlockInterval * param.BlockInterval
}
