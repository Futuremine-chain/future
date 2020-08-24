package validator

import "github.com/Futuremine-chain/future/types"

type IValidator interface {
	CheckMsg(types.IMessage, bool) error
}
