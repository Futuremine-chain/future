package chain

import "github.com/Futuremine-chain/futuremine/tools/arry"

type ChainDB struct {
}

func OpenChainDB(path string) (*ChainDB, error) {
	return nil, nil
}

func (c *ChainDB) ActRoot() arry.Hash {
	return arry.Hash{}
}

func (c *ChainDB) DPosRoot() arry.Hash {
	return arry.Hash{}
}

func (c *ChainDB) TokenRoot() arry.Hash {
	return arry.Hash{}
}

func (c *ChainDB) LastHeight() uint64 {
	return 0
}
