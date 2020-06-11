package types

import "github.com/Futuremine-chain/futuremine/tools/arry"

type TransactionBody struct {
	TokenAddress arry.Address
	Receiver     arry.Address
	Amount       uint64
}

func (t *TransactionBody) MsgTo() arry.Address {
	return t.Receiver
}
func (t *TransactionBody) CheckBody() error {
	return nil
}
func (t *TransactionBody) MsgAmount() uint64 {
	return t.Amount
}

type TokenBody struct {
	TokenAddress arry.Address
	Receiver     arry.Address
	Name         string
	Shorthand    string
	Amount       uint64
}

func (t *TokenBody) MsgTo() arry.Address {
	return arry.Address{}
}
func (t *TokenBody) CheckBody() error {
	return nil
}
func (t *TokenBody) MsgAmount() uint64 {
	return t.Amount
}

type CandidateBody struct {
}

func (c *CandidateBody) To() arry.Address {
	return arry.Address{}
}
func (c *CandidateBody) CheckBody() error {
	return nil
}

type CancelBody struct {
}

func (c *CancelBody) To() arry.Address {
	return arry.Address{}
}
func (c *CancelBody) CheckBody() error {
	return nil
}

type VoteBody struct {
}

func (v *VoteBody) To() arry.Address {
	return arry.Address{}
}
func (v *VoteBody) CheckBody() error {
	return nil
}
