package types

type Block struct {
	Header
	Body
}

type Blocks []*Block

func (b Blocks) Blocks() []*Block {
	return b
}
