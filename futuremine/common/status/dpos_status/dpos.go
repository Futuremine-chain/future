package dpos_status

import (
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/futuremine/db/status/dpos_db"
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

func (d *DPosStatus) CheckMessage(msg types.IMessage) error {
	return nil
}
