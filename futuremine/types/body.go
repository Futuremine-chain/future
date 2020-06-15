package types

import "github.com/Futuremine-chain/futuremine/types"

type Body struct {
	Messages
}

func (b *Body) Msgs() types.IMessages {
	return b.Messages
}
