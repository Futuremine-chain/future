package types

type RpcTransactionBody struct {
	Token  string `json:"token"`
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

func (r *RpcTransactionBody) GetTo() []byte {
	return []byte(r.To)
}
