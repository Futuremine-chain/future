package dpos_db

import (
	"github.com/Futuremine-chain/futuremine/common/db/base"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/trie"
)

type DPosDB struct {
	base *base.Base
	trie *trie.Trie
}

func Open(path string) (*DPosDB, error) {
	var err error
	baseDB, err := base.Open(path)
	if err != nil {
		return nil, err
	}
	return &DPosDB{base: baseDB}, nil
}

func (d *DPosDB) SetRoot(hash arry.Hash) error {
	t, err := trie.New(hash, d.base)
	if err != nil {
		return err
	}
	d.trie = t
	return nil
}

func (d *DPosDB) CandidatesCount() int {
	return 0
}
