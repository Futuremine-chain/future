package types

type RlpBody struct {
	Txs []*RlpMessage
}

func (r *RlpBody) ToBody() *Body {
	return nil
}
