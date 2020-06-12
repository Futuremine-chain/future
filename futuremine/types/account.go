package types

import (
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit"
	"github.com/Futuremine-chain/futuremine/futuremine/common/param"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/types"
)

const MaxAddressTxs = 1000

type Account struct {
	address arry.Address
	nonce   uint64
	Tokens  Tokens
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

func (a *Account) Balance(tokenAddr arry.Address) uint64 {
	token, ok := a.Tokens.Get(tokenAddr.String())
	if !ok {
		return 0
	}
	return token.Balance
}

func (a *Account) Check(msg types.IMessage) error {
	if !a.Exist() {
		a.address = msg.MsgBody().MsgTo()
	}

	if msg.Nonce() <= a.nonce {
		return fmt.Errorf("the nonce value of the message must be greater than %d", a.nonce)
	}

	// The nonce value cannot be greater than the
	// maximum number of address transactions
	if msg.Nonce() > a.nonce+MaxAddressTxs {
		return fmt.Errorf("the nonce value of the message cannot be greater "+
			"than the nonce value of the account %d", MaxAddressTxs)
	}

	// Verify the balance of the token
	switch MessageType(msg.Type()) {
	case Transaction:
		body, ok := msg.MsgBody().(*TransactionBody)
		if !ok {
			return errors.New("incorrect message type and message body")
		}
		if body.TokenAddress.IsEqual(param.MainToken) {
			return a.checkMainBalance(msg)
		} else {
			return a.checkTokenBalance(msg, body)
		}
	case Token:
		return a.checkTokenAmount(msg)
	default:
		if msg.MsgBody().MsgAmount() != 0 {
			return errors.New("wrong amount")
		}
		return a.checkFees(msg)
	}
}

// Verification contract amount
func (a *Account) checkTokenAmount(msg types.IMessage) error {
	body, ok := msg.MsgBody().(*TokenBody)
	if !ok {
		return errors.New("incorrect message type and message body")
	}
	amount := body.Amount
	consumption := kit.CalConsumption(amount)
	mainAddress := param.MainToken.String()
	main, ok := a.Tokens.Get(mainAddress)
	if !ok {
		return fmt.Errorf("it takes %f %s and %f %s as a handling fee to issue %f, and the balance is insufficient",
			Amount(consumption).ToCoin(),
			mainAddress,
			Amount(msg.Fee()).ToCoin(),
			mainAddress,
			Amount(amount).ToCoin())
	} else if main.Balance < msg.Fee()+consumption {
		return fmt.Errorf("it takes %f %s and %f %s as a handling fee to issue %f, and the balance is insufficient",
			Amount(consumption).ToCoin(),
			mainAddress,
			Amount(msg.Fee()).ToCoin(),
			mainAddress,
			Amount(amount).ToCoin())
	}
	return nil
}

// Verify the account balance of the primary transaction, the transaction
// value and transaction fee cannot be greater than the balance.
func (a *Account) checkMainBalance(msg types.IMessage) error {
	main := param.MainToken.String()
	token, ok := a.Tokens.Get(main)
	if !ok {
		return fmt.Errorf("%s does not have enough balance", main)
	} else if token.Balance < msg.Fee()+msg.MsgBody().MsgAmount() {
		return fmt.Errorf("%s does not have enough balance", main)
	}
	return nil
}

// Verify the account balance of the secondary transaction, the transaction
// value cannot be greater than the balance.
func (a *Account) checkTokenBalance(msg types.IMessage, body *TransactionBody) error {
	if err := a.checkFees(msg); err != nil {
		return err
	}

	coinAccount, ok := a.Tokens.Get(body.TokenAddress.String())
	if !ok {
		return fmt.Errorf("%s does not have enough balance", body.TokenAddress.String())
	} else if coinAccount.Balance < body.Amount {
		return fmt.Errorf("%s does not have enough balance", body.TokenAddress.String())
	}
	return nil
}

// Verification fee
func (a *Account) checkFees(msg types.IMessage) error {
	main := param.MainToken.String()
	token, ok := a.Tokens.Get(main)
	if !ok {
		return fmt.Errorf("%s does not have enough balance to pay the handling fee", main)
	} else if token.Balance < msg.Fee() {
		return fmt.Errorf("%s does not have enough balance to pay the handling fee", main)
	}
	return nil
}

func (a *Account) Exist() bool {
	return !arry.EmptyAddress(a.address)
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

type TokenAccount struct {
	Address   string
	Balance   uint64
	LockedIn  uint64
	LockedOut uint64
}

// List of secondary accounts
type Tokens []*TokenAccount

func (t *Tokens) Get(contract string) (*TokenAccount, bool) {
	for _, coin := range *t {
		if coin.Address == contract {
			return coin, true
		}
	}
	return &TokenAccount{}, false
}

func (t *Tokens) Set(newCoin *TokenAccount) {
	for i, coin := range *t {
		if coin.Address == newCoin.Address {
			(*t)[i] = newCoin
			return
		}
	}
	*t = append(*t, newCoin)
}
