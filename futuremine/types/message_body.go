package types

import (
	"errors"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit"
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

const PeerLength = 53

type Peer [PeerLength]byte

func (p Peer) String() string {
	return string(p[:])
}

func (p Peer) Bytes() []byte {
	return p[:]
}

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

func (t *TransactionBody) MsgToken() arry.Address {
	return t.TokenAddress
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

func (t *TokenBody) MsgToken() arry.Address {
	return t.TokenAddress
}

type CandidateBody struct {
	Peer Peer
}

func (c *CandidateBody) MsgTo() arry.Address {
	return arry.Address{}
}

func (c *CandidateBody) CheckBody() error {
	return nil
}

func (c *CandidateBody) MsgAmount() uint64 {
	return 0
}

func (c *CandidateBody) MsgToken() arry.Address {
	return config.Param.MainToken
}

type CancelBody struct {
}

func (c *CancelBody) MsgTo() arry.Address {
	return arry.Address{}
}

func (c *CancelBody) CheckBody() error {
	return nil
}

func (c *CancelBody) MsgToken() arry.Address {
	return config.Param.MainToken
}

func (c *CancelBody) MsgAmount() uint64 {
	return 0
}

type VoteBody struct {
	To arry.Address
}

func (v *VoteBody) MsgTo() arry.Address {
	return v.To
}

func (v *VoteBody) CheckBody() error {
	if !kit.CheckAddress(config.Param.Name, v.To) {
		return errors.New("wrong to address")
	}
	return nil
}

func (v *VoteBody) MsgToken() arry.Address {
	return config.Param.MainToken
}

func (v *VoteBody) MsgAmount() uint64 {
	return 0
}
