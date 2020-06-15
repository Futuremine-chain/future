package dpos_status

import (
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

type IDPosDB interface {
	SetRoot(hash arry.Hash) error
	Root() arry.Hash
	CandidatesCount() int
	Candidates() (*types.Candidates, error)
	CycleSupers(cycle int64) (*types.Supers, error)
	SaveCycle(cycle int64, supers *types.Supers)
	Voters() map[arry.Address][]arry.Address
}
