package status

import (
	"errors"
	"github.com/Futuremine-chain/future/common/dpos"
	fmctypes "github.com/Futuremine-chain/future/future/types"
	"github.com/Futuremine-chain/future/tools/arry"
	"github.com/Futuremine-chain/future/types"
)

const module = "chain"

type FMCStatus struct {
	actStatus   types.IActStatus
	dPosStatus  dpos.IDPosStatus
	tokenStatus types.ITokenStatus
}

func NewFMCStatus(actStatus types.IActStatus, dPosStatus dpos.IDPosStatus, tokenStatus types.ITokenStatus) *FMCStatus {
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

func (f *FMCStatus) Account(address arry.Address) types.IAccount {
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
	candidates, err := f.dPosStatus.CycleSupers(cycle)
	if err != nil {
		return fmctypes.NewSupers()
	}
	supers := candidates.(*fmctypes.Supers)
	for i, s := range supers.Candidates {
		supers.Candidates[i].MntCount = f.dPosStatus.SuperBlockCount(cycle, s.Signer)
	}
	return supers
}

func (f *FMCStatus) Token(address arry.Address) (types.IToken, error) {
	return f.tokenStatus.Token(address)
}
