package types

import "github.com/Futuremine-chain/future/tools/arry"

type IMessageBody interface {
	MsgTo() arry.Address
	MsgToken() arry.Address
	MsgAmount() uint64
	CheckBody(from arry.Address) error
}
