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
			Signer: arry.StringToAddress("xC8RqvGNhQ8sEpKrBHqnxJQh2rrtiJCXZrH"),
			PeerId: "16Uiu2HAm2YzMYX2Uw3aWGkKxsv7uLsRD7nnMkr4JWeGh79Q8b7tn",
			Weight: 0,
		},
		&types.Member{
			Signer: arry.StringToAddress("xBx2VZpE5RdC7cVfMgPet7rtWcmvtfRu7vH"),
			PeerId: "16Uiu2HAmBmDihnWL2XkrrQoGgHRfxWuACczXUWgL7o9FGwPQzF1q",
			Weight: 0,
		},
		&types.Member{
			Signer: arry.StringToAddress("xC3iVTuN6uupTwTZa2mFfHLiYKvGiQubiQA"),
			PeerId: "16Uiu2HAmNudiRRgeFcoN9nT2Q5YoHatrhTBiPiMGm5y9UGPjqYaU",
			Weight: 0,
		},
	},
}
