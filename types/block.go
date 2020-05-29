package types

import "github.com/Futuremine-chain/futuremine/types/rlp"

type IBlock interface {
	IHeader
	IBody

	ToRlp() rlp.IRlpTransaction
}
