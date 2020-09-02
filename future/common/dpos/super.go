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
			Signer: arry.StringToAddress("FMn1ryQna7NZLypDwFjVgQB6gGoRss4PMCz"),
			PeerId: "16Uiu2HAm9g3LZexdKiBCyYfbuZxuzeKa4TJQvmWQuaeYrAENHddz",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMoZybzhqLUgAvVty8YTTyA8kno3SHH2vZQ"),
			PeerId: "16Uiu2HAkzAYyVa6j21SAWfoWRKxYr4VPjAMeLmAyWYC9zShJLcJp",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMaVkP6SyLTb8vnNZkLU3hDL1hmUC9CBv8o"),
			PeerId: "16Uiu2HAm8KRurykapQpkb9GbvprhFxVVNhRnL43zDQqR6ar7SyvZ",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMqLJjrTXZpppJgM57xepzJyrvLU37HdZEm"),
			PeerId: "16Uiu2HAmLJYziKsVcixpkpnhTPoBVfnsLjJochL9x2Cy7hdDFVm5",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMe5S5EDzQeQg7WArGNskoB4w21dADGrhu6"),
			PeerId: "16Uiu2HAkuqzf694CenVRHjFMh55QKzBEEC7wYTjSR8RHHkQq6j2R",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMgMC6w8K2Ka44o3fcKsGENJ6aHtPoRusyG"),
			PeerId: "16Uiu2HAmTtL9sqbvodKZBCrqH2Pc9vdQrWcKWs7u7sirss9RuREh",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMn9swpExZL4XTrUG5jG51YUn12vuxZB2YK"),
			PeerId: "16Uiu2HAkvtMecWCxcxgFfN3MkBpVAbQ3TAXvyiCSddd2NtA9TT9F",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMa1ZNLstWs3PbPJC8LdbxdeCXidyxfjPn6"),
			PeerId: "16Uiu2HAmFVgtcVysvSDLvWgSsUQd772y17KXtZ1vjsJdvwCb6UG1",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMYpiZ1trVD9ThejVvXEkiPwrm8T8HydKiu"),
			PeerId: "16Uiu2HAmNck2dmR4qaRnmGjQm4vsz1iuaon1YvcUB51WYhJieqd9",
			Weight: 0,
		},
	},
}
