package dpos

import (
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IDPos interface {
	GenesisBlock() types.IBlock
	CheckTime(header types.IHeader, chain blockchain.IChain) error
	CheckSigner(header types.IHeader, chain blockchain.IChain) error
	CheckHeader(header types.IHeader, parent types.IHeader, chain blockchain.IChain) error
	CheckSeal(header types.IHeader, parent types.IHeader, chain blockchain.IChain) error
	Confirmed() uint64
	SetConfirmed(uint64)
}

type IDPosStatus interface {
	SetTrieRoot(hash arry.Hash) error
	TrieRoot() arry.Hash
	Confirmed() (uint64, error)
	SetConfirmed(height uint64)
	Candidates() (types.ICandidates, error)
	Voters() map[arry.Address][]arry.Address
	CycleSupers(cycle uint64) (types.ICandidates, error)
	SaveCycle(cycle uint64, supers types.ICandidates)
	CheckMessage(msg types.IMessage) error
	AddCandidate(msg types.IMessage) error
	CancelCandidate(msg types.IMessage) error
	Voter(msg types.IMessage) error
	AddSuperBlockCount(cycle uint64, signer arry.Address)
	SuperBlockCount(cycle uint64, signer arry.Address) uint32
	Commit() (arry.Hash, error)
}
