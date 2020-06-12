package dpos

import (
	"encoding/binary"
	"errors"
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/futuremine/common/param"
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/hash"
	"math/rand"
	"sort"
)

var Err_Elected = errors.New("the election has been passed")

type Cycle struct {
	DPosStatus dpos.IDPosStatus
}

func (c *Cycle) CheckCycle(chain blockchain.IBlockChain, preTime, time int64) error {
	currentTerm := time / param.BlockInterval

	supers, _ := c.DPosStatus.CycleSupers(currentTerm)
	if supers != nil && len(supers.Candidates) != 0 {
		return Err_Elected
	}
	return nil
}

func (c *Cycle) Elect(time int64, preHash arry.Hash, chain blockchain.IBlockChain) error {
	curCycle := time / param.BlockInterval
	voters, err := c.calVotes(chain)
	if err != nil {
		return err
	}
	candidates := types.SortableCandidates{}
	for _, candidate := range voters {
		candidates = append(candidates, candidate)
	}
	if len(candidates) < param.SafeSize {
		return errors.New("too few candidate")
	}

	sort.Sort(candidates)

	if len(candidates) > SuperCount {
		candidates = candidates[:SuperCount]
	}

	// Use the last block hash of the last cycle as a random number seed
	// to ensure that the election results of each node are consistent
	seed := int64(binary.LittleEndian.Uint32(hash.Hash(preHash.Bytes()).Bytes())) + int64(curCycle)
	r := rand.New(rand.NewSource(seed))
	for i := len(candidates) - 1; i > 0; i-- {
		j := int(r.Int31n(int32(i + 1)))
		candidates[i], candidates[j] = candidates[j], candidates[i]
	}

	supers := &types.Supers{Candidates: candidates, PreHash: preHash}
	c.DPosStatus.SaveCycle(curCycle, supers)
	return nil
}

func (c *Cycle) calVotes(chain blockchain.IBlockChain) ([]*types.Member, error) {
	candidates, err := c.DPosStatus.Candidates()
	if err != nil {
		return nil, errors.New("no candidate")
	}
	if len(candidates.Members) < SuperCount {
		return nil, errors.New("not enough candidates")
	}
	for index, candidate := range candidates.Members {
		voters := c.DPosStatus.Voters(candidate.Signer)
		for _, voter := range voters {
			candidates.Members[index].Weight += chain.Vote(voter)
		}
	}
	return candidates.Members, nil
}
