package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

type IHeader interface {
	GetHash() arry.Hash
	GetPreHash() arry.Hash
	GetMsgRoot() arry.Hash
	GetActRoot() arry.Hash
	GetDPosRoot() arry.Hash
	GetTokenRoot() arry.Hash
	GetSigner() arry.Address
	GetSignature() ISignature
	GetHeight() uint64
	GetTime() int64
	GetCycle() int64
	ToRlpHeader() IRlpHeader
}
