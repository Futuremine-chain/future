package dpos_status

import (
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

type IDPosDB interface {
	SetRoot(hash arry.Hash) error
	Root() arry.Hash
	Commit() (arry.Hash, error)
	CandidatesCount() int
	Candidates() (*types.Candidates, error)
	AddCandidate(member *types.Member)
	CancelCandidate(signer arry.Address)
	CycleSupers(cycle int64) (*types.Supers, error)
	SaveCycle(cycle int64, supers *types.Supers)
	Voters() map[arry.Address][]arry.Address
	Confirmed() (uint64, error)
	SetConfirmed(uint64)
	Voter(from, to arry.Address)
	AddSuperBlockCount(cycle int64, signer arry.Address)
	SuperBlockCount(cycle int64, signer arry.Address) int
}
