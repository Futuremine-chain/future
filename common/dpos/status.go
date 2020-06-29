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
	CycleSupers(cycle uint64) (*fmctypes.Supers, error)
	SaveCycle(cycle uint64, supers *fmctypes.Supers)
	CheckMessage(msg types.IMessage) error
	AddCandidate(msg types.IMessage) error
	CancelCandidate(msg types.IMessage) error
	Voter(msg types.IMessage) error
	AddSuperBlockCount(cycle uint64, signer arry.Address)
	SuperBlockCount(cycle uint64, signer arry.Address) uint32
	Commit() (arry.Hash, error)
}
