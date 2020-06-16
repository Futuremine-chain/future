package types

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IMessageBody interface {
	MsgTo() arry.Address
	MsgToken()arry.Address
	MsgAmount() uint64
	CheckBody() error
}
