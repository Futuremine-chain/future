package param

import "github.com/Futuremine-chain/futuremine/tools/arry"

const (
	// Maximum amount limit for second-level tokens
	MaxCoinCount = 90000000000
	// The minimum amount limit of the second-level token
	MinCoinCount = 1000

	Proportion = 10000

	CoinBase = 1000000000
)

var MainToken = arry.StringToAddress("FMC")

const (
	// Block interval period
	BlockInterval = int64(5)
	// Re-election interval
	CycleInterval = 60 * 60 * 24
	// Maximum number of super nodes
	SuperSize = 9
	// The minimum number of nodes required to confirm the transaction
	SafeSize = SuperSize*2/3 + 1
	// The minimum threshold at which a block is valid
	DPosSize = SuperSize*2/3 + 1
)
