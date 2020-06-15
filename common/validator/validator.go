package validator

import "github.com/Futuremine-chain/futuremine/types"

type IValidator interface {
	CheckMsg(types.IMessage, bool) error
}
