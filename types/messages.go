package types

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IMessages interface {
	Add(msg IMessage)
	Msgs() []IMessage
	MsgRoot() arry.Hash
	CalculateFee() uint64
}
