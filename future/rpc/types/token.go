package types

import (
	"github.com/Futuremine-chain/future/future/types"
	"github.com/Futuremine-chain/future/tools/amount"
)

type RpcTokenBody struct {
	Address        string `json:"address"`
	Receiver       string `json:"receiver"`
	Name           string `json:"name"`
	Shorthand      string `json:"shorthand"`
	Amount         uint64 `json:"amount"`
	IncreaseIssues bool   `json:"allowedincrease"`
}

type RpcToken struct {
	Address        string    `json:"address"`
	Sender         string    `json:"sender"`
	Name           string    `json:"name"`
	Shorthand      string    `json:"shorthand"`
	IncreaseIssues bool      `json:"increaseissues"`
	Records        []*Record `json:"records"`
}

type Record struct {
	Height   uint64  `json:"height"`
	Receiver string  `json:"receiver"`
	MsgHash  string  `json:"msghash"`
	Time     uint64  `json:"time"`
	Amount   float64 `json:"amount"`
}

func TokenToRpcToken(token *types.TokenRecord) *RpcToken {
	rpcToken := &RpcToken{
		Address:        token.Address.String(),
		Sender:         token.Sender.String(),
		Name:           token.Name,
		Shorthand:      token.Shorthand,
		IncreaseIssues: token.IncreaseIssues,
		Records:        make([]*Record, token.Records.Len()),
	}
	for i, record := range *token.Records {
		rpcToken.Records[i] = &Record{
			Height:   record.Height,
			MsgHash:  record.MsgHash.String(),
			Receiver: record.Receiver.String(),
			Time:     record.Time,
			Amount:   amount.Amount(record.Amount).ToCoin(),
		}
	}
	return rpcToken
}
