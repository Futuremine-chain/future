package types

import (
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/math"
)

const (
	PeerLength = 53
	MaxName    = 100
)

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

func (t *TransactionBody) CheckBody(from arry.Address) error {
	if !kit.CheckAddress(config.Param.Name, t.Receiver) {
		return errors.New("receive address verification failed")
	}
	if !t.TokenAddress.IsEqual(config.Param.MainToken) {
		if !kit.CheckTokenAddress(config.Param.Name, t.TokenAddress) {
			return errors.New("token address verification failed")
		}
	}
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

func (t *TokenBody) CheckBody(from arry.Address) error {
	if !kit.CheckAddress(config.Param.Name, t.Receiver) {
		return errors.New("receive address verification failed")
	}
	if !kit.CheckTokenAddress(config.Param.Name, t.TokenAddress) {
		return errors.New("token address verification failed")
	}
	toKenAddr, err := kit.GenerateTokenAddress(config.Param.Name, from, t.Shorthand)
	if err != nil {
		return errors.New("token address verification failed")
	}
	if !toKenAddr.IsEqual(t.TokenAddress) {
		return errors.New("token address verification failed")
	}
	if err := kit.CheckShorthand(t.Shorthand); err != nil {
		return fmt.Errorf("shorthand verification failed, %s", err.Error())
	}
	if len(t.Name) > MaxName {
		return fmt.Errorf("the maximum length of the token name is %d", MaxName)
	}
	if t.Amount > math.MaxInt64 {
		return fmt.Errorf("amount cannot be greater than %d", math.MaxInt64)
	}
	fAmount := Amount(t.Amount).ToCoin()
	if fAmount < config.Param.MinCoinCount || fAmount > config.Param.MaxCoinCount {
		return fmt.Errorf("the quantity of coins must be between %f and %f", config.Param.MinCoinCount, config.Param.MaxCoinCount)
	}
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

func (c *CandidateBody) CheckBody(from arry.Address) error {
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

func (c *CancelBody) CheckBody(from arry.Address) error {
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

func (v *VoteBody) CheckBody(from arry.Address) error {
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
