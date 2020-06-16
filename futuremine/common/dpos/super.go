package dpos

import (
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

// initialCandidates the first super node of the block generation cycle.
// The first half is the address of the block, the second half is the id of the block node
var genesisSuperList = types.Candidates{
	Members: []*types.Member{
		&types.Member{
			Signer: arry.Address{},
			PeerId: "",
			Weight: 0,
		},
		&types.Member{
			Signer: arry.Address{},
			PeerId: "",
			Weight: 0,
		},
		&types.Member{
			Signer: arry.Address{},
			PeerId: "",
			Weight: 0,
		},
	},
}
