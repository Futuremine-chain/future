package types

import (
	"errors"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type Message struct {
	Header *Message_Header
	Body   types.IMessageBody
}

func (t *Message) Hash() arry.Hash {
	panic("implement me")
}

func (t *Message) From() arry.Address {
	panic("implement me")
}

func (t *Message) Nonce() uint64 {
	panic("implement me")
}

func (t *Message) Fee() uint64 {
	panic("implement me")
}

func (t *Message) Time() int64 {
	panic("implement me")
}

func (t *Message) IsCoinBase() bool {
	panic("implement me")
}

func (t *Message) To() arry.Address {
	panic("implement me")
}

func (t *Message) ToRlp() types.IRlpMessage {
	panic("implement me")
}

func (t *Message) Check() error {
	if t.Header == nil || t.Body == nil {
		return errors.New("incomplete message")
	}

	if err := t.Header.Check(); err != nil {
		return err
	}

	if err := t.Body.CheckBody(); err != nil {
		return err
	}
	return nil
}

func (t *Message) CheckBody() error {
	if t.Header == nil || t.Body == nil {
		return errors.New("incomplete message")
	}

	return t.CheckBody()
}

/*
func (t *Message) verifyHead() error {


	if err := t.verifyTxType(); err != nil {
		return err
	}

	if err := t.verifyTxHash(); err != nil {
		return err
	}

	if err := t.verifyTxFrom(); err != nil {
		return err
	}

	if err := t.verifyTxFees(); err != nil {
		return err
	}

	if err := t.verifyTxSinger(); err != nil {
		return err
	}
	return nil
}

func (t *Message) verifyBody() error {
	if t.TxBody == nil {
		return ErrTxBody
	}

	if err := t.verifyAmount(); err != nil {
		return err
	}

	if err := t.TxBody.VerifyBody(t.TxHead.From); err != nil {
		return err
	}
	return nil
}

func (t *Message) VerifyCoinBaseTx(sumFees uint64) error {
	if err := t.verifyTxSize(); err != nil {
		return err
	}

	if err := t.verifyCoinBaseAmount(sumFees); err != nil {
		return err
	}
	return nil
}

func (t *Message) verifyTxFees() error {
	minFees, maxFees := t.FeesLimit()
	if t.TxHead.Fees < minFees {
		return fmt.Errorf("fee %d is less than the minimum poundage allowed %d", t.TxHead.Fees, minFees)
	}
	if t.TxHead.Fees > maxFees {
		return fmt.Errorf("fee %d greater is greater than the maximum poundage allowed %d", t.TxHead.Fees, maxFees)
	}
	return nil
}

func (t *Message) verifyTxSinger() error {
	if !Verify(t.TxHead.TxHash, t.TxHead.SignScript) {
		return ErrSignature
	}

	if !VerifySigner(param.Net, t.TxHead.From, t.TxHead.SignScript.PubKey) {
		return ErrSigner
	}
	return nil
}

func (t *Message) verifyTxSize() error {
	// TODO change maxsize
	switch t.TxHead.TxType {
	case NormalMessage:
		fallthrough
	case ContractMessage:
		fallthrough
	case LogoutCandidate:
		fallthrough
	case LoginCandidate:
		fallthrough
	case VoteToCandidate:
		if t.Size() > MaxNoDataTxSize {
			return ErrTxSize
		}
	}
	return nil
}

func (t *Message) verifyCoinBaseAmount(amount uint64) error {
	nTx := t.TxBody.(*NormalMessageBody)
	sumAmount := CoinBaseCoins + amount
	if sumAmount != nTx.Amount {
		return ErrCoinBase
	}
	return nil
}

func (t *Message) verifyAmount() error {
	nTx, ok := t.TxBody.(*NormalMessageBody)
	if ok && nTx.Amount < minAllowedAmount {
		return fmt.Errorf("the minimum amount of the message must not be less than %d", minAllowedAmount)
	}
	return nil
}

func (t *Message) verifyTxFrom() error {
	if !CheckUBAddress(param.Net, t.From().String()) {
		return ErrAddress
	}
	return nil
}

func (t *Message) verifyTxType() error {
	switch t.TxHead.TxType {
	case NormalMessage:
		return nil
	case ContractMessage:
		return nil
	case VoteToCandidate:
		return nil
	case LoginCandidate:
		return nil
	case LogoutCandidate:
		return nil
	}
	return ErrTxType
}

func (t *Message) verifyTxHash() error {
	newTx := t.copy()
	newTx.SetHash()
	if newTx.Hash().IsEqual(t.Hash()) {
		return nil
	}
	return ErrTxHash
}
*/
