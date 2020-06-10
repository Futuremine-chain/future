package account

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IActStatus interface {
	Nonce(arry.Address) uint64
}
