package types

type IBlock interface {
	IHeader
	IBody
	Header() IHeader
	ToRlp() IRlpTransaction
}

type IBlocks interface {
	Blocks() []IBlock
}
