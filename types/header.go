package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

type IHeader interface {
	Hash() arry.Hash
	TxRoot() arry.Hash
	Signer() arry.Address
	Height() uint64
	Time() int64
	ToRlpHeader() IRlpHeader
}
