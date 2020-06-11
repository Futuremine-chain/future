package kit

import "github.com/Futuremine-chain/futuremine/futuremine/common/param"

func CalConsumption(amount uint64) uint64 {
	if float64(amount)/param.Proportion < 1 {
		return 1
	} else if amount%param.Proportion != 0 {
		return amount/param.Proportion + 1
	}
	return amount / param.Proportion
}
