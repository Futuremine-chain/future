package dpos_status

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IDPosDB interface {
	SetRoot(hash arry.Hash) error
	CandidatesCount() int
}
