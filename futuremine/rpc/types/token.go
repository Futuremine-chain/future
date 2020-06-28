package types

type RpcTokenBody struct {
	Address   string `json:"address"`
	To        string `json:"to"`
	Name      string `json:"name"`
	Shorthand string `json:"shorthand"`
	Amount    uint64 `json:"amount"`
}
