package account

import "github.com/Futuremine-chain/futuremine/futuremine/common/arry"

type IActStatus interface {
	Nonce(arry.Address) uint64
}
