package dpos_db

import (
	"github.com/Futuremine-chain/futuremine/common/db/base"
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/tools/trie"
)

const (
	_cycleSupers = "cycleSupers"
	_candidates  = "candidates"
	_voters      = "voters"
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

func (d *DPosDB) CycleSupers(cycle int64) (*types.Supers, error) {
	var supers *types.Supers
	cycleBytes, err := rlp.EncodeToBytes(cycle)
	if err != nil {
		return nil, err
	}
	bytes, err := d.base.GetFromBucket(_cycleSupers, cycleBytes)
	if err := rlp.DecodeBytes(bytes, &supers); err != nil {
		return nil, err
	}
	return supers, nil
}

func (d *DPosDB) SaveCycle(cycle int64, supers *types.Supers) {
	value, _ := rlp.EncodeToBytes(supers)
	key, _ := rlp.EncodeToBytes(cycle)
	d.base.PutInBucket(_cycleSupers, key, value)
}

func (d *DPosDB) Candidates() (*types.Candidates, error) {
	var candidates *types.Candidates
	bytes := d.trie.Get(base.Key(_candidates, []byte(_candidates)))
	if err := rlp.DecodeBytes(bytes, &candidates); err != nil {
		return nil, err
	}
	return candidates, nil
}

func (d *DPosDB) Voters(addr arry.Address) []arry.Address {
	rs := make([]arry.Address, 0)
	iter := d.trie.PrefixIterator([]byte(_voters))
	if iter.Leaf() {
		key := iter.LeafKey()
		from := arry.BytesToAddress(key)
		value := iter.LeafKey()
		to := arry.BytesToAddress(value)
		if to.IsEqual(addr) {
			rs = append(rs, from)
		}
		if !iter.Next(true) {
			return rs
		}
	}
	return rs
}
