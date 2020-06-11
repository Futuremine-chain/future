package dpos_status

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/futuremine/db/status/dpos_db"
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

const dPosDB = "dpos_db"

type DPosStatus struct {
	db IDPosDB
}

func NewDPosStatus() (*DPosStatus, error) {
	db, err := dpos_db.Open(config.App.Setting().Data + "/" + dPosDB)
	if err != nil {
		return nil, err
	}
	return &DPosStatus{db: db}, nil
}

func (d *DPosStatus) SetTrieRoot(hash arry.Hash) error {
	return d.db.SetRoot(hash)
}

// If the current number of candidates is less than or equal to the
// number of super nodes, it is not allowed to withdraw candidates.
func (d *DPosStatus) CheckMessage(msg types.IMessage) error {
	switch fmctypes.MessageType(msg.Type()) {
	case fmctypes.Cancel:
		if d.db.CandidatesCount() <= dpos.SuperCount {
			return fmt.Errorf("candidate nodes are already in the minimum number. Cannot cancel the candidate status now, please wait")
		}
	}
	return nil

}
