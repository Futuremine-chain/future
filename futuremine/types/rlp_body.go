package types

type RlpBody struct {
	Txs []*RlpTransaction
}

func (r *RlpBody) ToBody() *Body {
	return nil
}
