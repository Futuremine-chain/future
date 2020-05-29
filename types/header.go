package types

import "github.com/Futuremine-chain/futuremine/futuremine/common/arry"

type IHeader interface {
	Hash() arry.Hash
	Signer() arry.Address
	Height() uint64
}
