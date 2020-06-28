package types

import (
	"github.com/Futuremine-chain/futuremine/futuremine/types"
)

type Account struct {
	Address   string       `json:"address"`
	Nonce     uint64       `json:"nonce"`
	Tokens    Tokens `json:"tokens"`
	Confirmed uint64       `json:"confirmed"`
}

type TokenAccount struct {
	Address   string `json:"address"`
	Balance   float64 `json:"balance"`
	LockedIn  float64 `json:"locked"`
}

// List of secondary accounts
type Tokens []*TokenAccount

func ToRpcAccount(a *types.Account) *Account {
	tokens := make(Tokens, len(a.Tokens))
	for i, t := range a.Tokens{
		tokens[i] = &TokenAccount{
			Address:  t.Address,
			Balance:  types.Amount(t.Balance).ToCoin(),
			LockedIn: types.Amount(t.LockedIn).ToCoin(),
		}
	}
	return &Account{
		Address:   a.Address.String(),
		Nonce:     a.Nonce,
		Tokens:    tokens,
		Confirmed: a.Confirmed,
	}
}
