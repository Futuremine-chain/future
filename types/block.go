package types

type IBlock interface {
	IHeader
	BlockHeader() IHeader
	ToRlpBlock() IRlpBlock
}

type IBlocks interface {
	Blocks() []IBlock
}
