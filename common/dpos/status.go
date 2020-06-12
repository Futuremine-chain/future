package dpos

import (
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IDPosStatus interface {
	SetTrieRoot(hash arry.Hash) error
	CheckMessage(msg types.IMessage) error
	Candidates() (*fmctypes.Candidates, error)
	Voters(addr arry.Address) []arry.Address
	CycleSupers(cycle int64) (*fmctypes.Supers, error)
	SaveCycle(cycle int64, supers *fmctypes.Supers)
}
