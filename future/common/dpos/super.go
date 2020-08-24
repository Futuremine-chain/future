package dpos

import (
	"github.com/Futuremine-chain/future/future/types"
	"github.com/Futuremine-chain/future/tools/arry"
)

// initialCandidates the first super node of the block generation cycle.
// The first half is the address of the block, the second half is the id of the block node
var genesisSuperList = types.Candidates{
	Members: []*types.Member{
		{
			Signer: arry.StringToAddress("xC6nGfqHPt4KPoyhwCG3zNemNweYsBWXRR2"),
			PeerId: "16Uiu2HAkvG7xGvVZBkZR16sWwmqhqGNm69FSa61tbpphpgCcMmEo",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("xC3LLtxiCaAcG2KrJmzo75pdavgqNLooz1N"),
			PeerId: "16Uiu2HAmSUvLmVizqJYNraZ5DUQQGHCgDL6XvRgoTj9KY9oEZoba",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("xC2xC5DNc3Uz3jMnufsX9jX8a67mypdfxGH"),
			PeerId: "16Uiu2HAmDESGfzayFiBKnKaryyVgXKmJvePva8nFVNzFsGKkLdW2",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("xCJVPCgbiRy6tAtozAZwXLF5NzG9XBU16ou"),
			PeerId: "16Uiu2HAkyexJCRYMHrN2Ft3E5MmJr3eoNtEZNJjaZ4jbwXcA1FCK",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("xC554Aq8EiCz82hzVfi5PmnLCFGxAi9fBMp"),
			PeerId: "16Uiu2HAm9KiJJfE8mhwFuLv8Yco1YJkFtMekzUPowMW4t8FeLbw7",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("xC86k2HFGPzsJ1vubBoR7DXviN6UKpEKFRv"),
			PeerId: "16Uiu2HAm5VouADKhW55qLBiFoixegEtym72C9P6M2XKZ7bEMNUt8",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("xC1cGYmjAEU3yymozKs5jPLQPzkyzbvxJuU"),
			PeerId: "16Uiu2HAmLk1r8BHEpuTR9XRm1DCw5GpqaANF2RYJ3DRzQxtAtSc9",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("xC8hzrNRzMzpUcs68HcWqFTLBjWKeTatqcJ"),
			PeerId: "16Uiu2HAm92KMJiSLAYFCc7AdTBFwVgCBJ9zhWYU6UHM5vMAWH1hr",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("xBzSZn52pNuBsUhu3qnRYTqo8HetVpgrRVL"),
			PeerId: "16Uiu2HAmHBNYs8eFMnjpUwoJzCws77B8LuSXqVwKaQXcA3WVDvvQ",
			Weight: 0,
		},
	},
}
