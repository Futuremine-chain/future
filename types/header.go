package types

import (
	"github.com/Futuremine-chain/futuremine/futuremine/common/arry"
	"time"
)

type IHeader interface {
	Hash() arry.Hash
	Signer() arry.Address
	Height() uint64
	Time() time.Time
}
