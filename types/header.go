package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

type IHeader interface {
	Hash() arry.Hash
	PreHash() arry.Hash
	MsgRoot() arry.Hash
	ActRoot() arry.Hash
	DPosRoot() arry.Hash
	TokenRoot() arry.Hash
	Signer() arry.Address
	Signature() ISignature
	Height() uint64
	Time() int64
	Cycle() int64
	ToRlpHeader() IRlpHeader
}
