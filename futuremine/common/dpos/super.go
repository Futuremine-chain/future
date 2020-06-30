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
			Signer: arry.StringToAddress("xCJkBm6ZZRMBEZiXbWbk8HCaR4mxz1oJxoH"),
			PeerId: "16Uiu2HAm89g96hN3TYiKnY2Vyf6XVyxsJQRwQyep659AeFW2xYTT",
			Weight: 0,
		},
		&types.Member{
			Signer: arry.StringToAddress("xCGEcommDH2N8Ev9pGxiWmyDVf9qkAHXqps"),
			PeerId: "16Uiu2HAmPqPQg98uYn3PJbfARj4PUJBboWpg3jeTJSZA6teGvV1H",
			Weight: 0,
		},
		&types.Member{
			Signer: arry.StringToAddress("xCBV6JvhbpjBdbwDVvuA9C26iyVqY13TWeS"),
			PeerId: "16Uiu2HAmBgafEgfCW4nqjLDg8eWjEFHhbpEFdM4j4jaj9nYDtrmr",
			Weight: 0,
		},
	},
}
