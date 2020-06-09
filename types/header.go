package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"time"
)

type IHeader interface {
	Hash() arry.Hash
	Signer() arry.Address
	Height() uint64
	Time() time.Time
	ToRlpHeader() IRlpHeader
}
