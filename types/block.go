package types

type IBlock interface {
	IHeader
	IBody
	Header() IHeader
	ToRlpBlock() IRlpBlock
}

type IBlocks interface {
	Blocks() []IBlock
}
