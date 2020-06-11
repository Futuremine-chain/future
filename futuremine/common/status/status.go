package status

import (
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/common/token"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type FMCStatus struct {
	actStatus   account.IActStatus
	dPosStatus  dpos.IDPosStatus
	tokenStatus token.ITokenStatus
}

func NewFMCStatus(actStatus account.IActStatus, dPosStatus dpos.IDPosStatus, tokenStatus token.ITokenStatus) *FMCStatus {
	return &FMCStatus{
		actStatus:   actStatus,
		dPosStatus:  dPosStatus,
		tokenStatus: tokenStatus,
	}
}

func (f *FMCStatus) InitRoots(actRoot, dPosRoot, tokenRoot arry.Hash) error {
	if err := f.actStatus.SetTrieRoot(actRoot); err != nil {
		return err
	}
	if err := f.dPosStatus.SetTrieRoot(dPosRoot); err != nil {
		return err
	}
	if err := f.tokenStatus.SetTrieRoot(tokenRoot); err != nil {
		return err
	}
	return nil
}

func (f *FMCStatus) SetConfirmed(confirmed uint64) {
	f.actStatus.SetConfirmed(confirmed)
}

func (f *FMCStatus) Check(tx types.ITransaction) error {
	if err := tx.Check(); err != nil {
		return err
	}

	if err := f.dPosStatus.CheckTransaction(tx); err != nil {
		return err
	}

	if err := f.actStatus.CheckTransaction(tx); err != nil {
		return err
	}

	if err := f.tokenStatus.CheckTransaction(tx); err != nil {
		return err
	}
	return nil
}
