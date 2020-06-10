package status

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IStatus interface {
	InitRoots(actRoot, dPosRoot, tokenRoot arry.Hash) error
	SetConfirmed(confirmed uint64)
}
