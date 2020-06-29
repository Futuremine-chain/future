package status

import (
	"errors"
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/common/token"
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

const module = "chain"

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

func (f *FMCStatus) Account(address arry.Address) account.IAccount {
	return f.actStatus.Account(address)
}

func (f *FMCStatus) CheckMsg(msg types.IMessage, strict bool) error {
	if err := msg.Check(); err != nil {
		return err
	}

	if err := f.dPosStatus.CheckMessage(msg); err != nil {
		return err
	}

	if err := f.actStatus.CheckMessage(msg, strict); err != nil {
		return err
	}

	if err := f.tokenStatus.CheckMessage(msg); err != nil {
		return err
	}
	return nil
}

func (f *FMCStatus) Change(msgs []types.IMessage, block types.IBlock) error {
	for _, msg := range msgs {
		switch fmctypes.MessageType(msg.Type()) {
		case fmctypes.Transaction:
			if err := f.actStatus.ToMessage(msg, block.GetHeight()); err != nil {
				return err
			}
		case fmctypes.Token:
			if err := f.actStatus.ToMessage(msg, block.GetHeight()); err != nil {
				return err
			}
			if err := f.tokenStatus.UpdateToken(msg, block.GetHeight()); err != nil {
				return err
			}
		case fmctypes.Vote:
			if err := f.dPosStatus.Voter(msg); err != nil {
				return nil
			}
		case fmctypes.Candidate:
			if err := f.dPosStatus.AddCandidate(msg); err != nil {
				return nil
			}
		case fmctypes.Cancel:
			if err := f.dPosStatus.CancelCandidate(msg); err != nil {
				return nil
			}
		default:
			return errors.New("wrong message type")
		}
		if err := f.actStatus.FromMessage(msg, block.GetHeight()); err != nil {
			return err
		}

	}
	f.dPosStatus.AddSuperBlockCount(block.GetCycle(), block.GetSigner())
	return nil
}

func (f *FMCStatus) Commit() (arry.Hash, arry.Hash, arry.Hash, error) {
	actRoot, err := f.actStatus.Commit()
	if err != nil {
		return arry.Hash{}, arry.Hash{}, arry.Hash{}, err
	}
	tokenRoot, err := f.tokenStatus.Commit()
	if err != nil {
		return arry.Hash{}, arry.Hash{}, arry.Hash{}, err
	}
	dPosRoot, err := f.dPosStatus.Commit()
	if err != nil {
		return arry.Hash{}, arry.Hash{}, arry.Hash{}, err
	}
	return actRoot, tokenRoot, dPosRoot, nil
}

func (f *FMCStatus) Candidates() types.ICandidates {
	cans, _ := f.dPosStatus.Candidates()
	return cans
}

func (f *FMCStatus) CycleSupers(cycle uint64) types.ICandidates {
	supers, err := f.dPosStatus.CycleSupers(cycle)
	if err != nil {
		return fmctypes.NewSupers()
	}
	for i, s := range supers.Candidates {
		supers.Candidates[i].MntCount = f.dPosStatus.SuperBlockCount(cycle, s.Signer)
	}
	return supers
}
