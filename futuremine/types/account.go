package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/types"
)

type Account struct {
	address arry.Address
	nonce   uint64
}

func (a *Account) NeedUpdate() bool {
	panic("implement me")
}

func (a *Account) UpdateLocked(confirmed uint64) error {
	panic("implement me")
}

func (a *Account) FromMessage(msg types.IMessage, height uint64) error {
	panic("implement me")
}

func (a *Account) ToMessage(msg types.IMessage, height uint64) error {
	panic("implement me")
}

func (a *Account) Check(msg types.IMessage) error {
	panic("implement me")
}

func (a *Account) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(a)
	return bytes
}

func (a *Account) Address() arry.Address {
	return a.address
}

func (a *Account) Nonce() uint64 {
	return a.nonce
}

func NewAccount() *Account {
	return &Account{}
}

func DecodeAccount(bytes []byte) (*Account, error) {
	var account *Account
	err := rlp.DecodeBytes(bytes, account)
	return account, err
}
