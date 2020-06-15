package dpos

import (
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IDPosStatus interface {
	SetTrieRoot(hash arry.Hash) error
	TrieRoot() arry.Hash
	Confirmed() (uint64, error)
	SetConfirmed(height uint64)
	Candidates() (*fmctypes.Candidates, error)
	Voters() map[arry.Address][]arry.Address
	CycleSupers(cycle int64) (*fmctypes.Supers, error)
	SaveCycle(cycle int64, supers *fmctypes.Supers)
	CheckMessage(msg types.IMessage) error
}
