package types

type IBlock interface {
	IHeader
	IBody
	BlockHeader() IHeader
	ToRlpBlock() IRlpBlock
}

type IBlocks interface {
	Blocks() []IBlock
}
