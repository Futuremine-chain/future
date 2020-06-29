package types

import (
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/types"
)

type RpcMember struct {
	Signer   string `json:"address"`
	PeerId   string `json:"peerid"`
	Weight   uint64 `json:"votes"`
	MntCount uint32 `json:"mntcount"`
}

type RpcCandidates struct {
	Members []*RpcMember `json:"members"`
}

func CandidatesToRpcCandidates(candidates types.ICandidates) *RpcCandidates {
	rpcMems := &RpcCandidates{Members: make([]*RpcMember, 0)}
	cas := candidates.(*fmctypes.Candidates)
	for _, candidate := range cas.Members {
		rpcMem := &RpcMember{
			Signer: candidate.Signer.String(),
			PeerId: candidate.PeerId,
			Weight: candidate.Weight,
		}
		rpcMems.Members = append(rpcMems.Members, rpcMem)
	}
	return rpcMems
}

func SupersToRpcCandidates(candidates types.ICandidates) *RpcCandidates {
	rpcMems := &RpcCandidates{Members: make([]*RpcMember, 0)}
	supers := candidates.(*fmctypes.Supers)
	for _, candidate := range supers.Candidates {
		rpcMem := &RpcMember{
			Signer:   candidate.Signer.String(),
			PeerId:   candidate.PeerId,
			Weight:   candidate.Weight,
			MntCount: candidate.MntCount,
		}
		rpcMems.Members = append(rpcMems.Members, rpcMem)
	}
	return rpcMems
}
