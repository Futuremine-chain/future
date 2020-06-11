package types

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IMessageBody interface {
	To() arry.Address
	CheckBody() error
}
