package types

type RlpBody struct {
	Msgs []*RlpMessage
}

func (r *RlpBody) ToBody() *Body {
	return nil
}

func (r *RlpBody) MsgList() []*RlpMessage {
	return r.Msgs
}
