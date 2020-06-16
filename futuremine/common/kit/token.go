package kit

func CalConsumption(amount uint64, proportion uint64) uint64 {
	if float64(amount)/float64(proportion) < 1 {
		return 1
	} else if amount%proportion != 0 {
		return amount/proportion + 1
	}
	return amount / proportion
}
