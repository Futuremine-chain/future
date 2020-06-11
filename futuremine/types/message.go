package types

import (
	"errors"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	"github.com/Futuremine-chain/futuremine/tools/crypto/hash"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/types"
)

type Message struct {
	Header *MsgHeader
	Body   types.IMessageBody
}

func (m *Message) Hash() arry.Hash {
	panic("implement me")
}

func (m *Message) From() arry.Address {
	panic("implement me")
}

func (m *Message) Nonce() uint64 {
	panic("implement me")
}

func (m *Message) Fee() uint64 {
	panic("implement me")
}

func (m *Message) Time() int64 {
	panic("implement me")
}

func (m *Message) IsCoinBase() bool {
	panic("implement me")
}

func (m *Message) To() arry.Address {
	panic("implement me")
}

func (m *Message) ToRlp() types.IRlpMessage {
	panic("implement me")
}

func (m *Message) Check() error {
	if m.Header == nil || m.Body == nil {
		return errors.New("incomplete message")
	}

	if err := m.checkHash(); err != nil {
		return err
	}

	if err := m.Header.Check(); err != nil {
		return err
	}

	if err := m.Body.CheckBody(); err != nil {
		return err
	}
	return nil
}

func (m *Message) CheckBody() error {
	if m.Header == nil || m.Body == nil {
		return errors.New("incomplete message")
	}

	return m.CheckBody()
}

func (m *Message) checkHash() error {
	newMsg := m.copy()
	newMsg.SetHash()
	if newMsg.Hash().IsEqual(m.Hash()) {
		return nil
	}
	return errors.New("error messages hash")
}

func (m *Message) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(m)
	return bytes
}

func (m *Message) SignMsg(key *secp256k1.PrivateKey) error {
	var err error
	if m.Header.Signature, err = Sign(key, m.Header.Hash); err != nil {
		return err
	}
	return nil
}

func (m *Message) SetHash() {
	m.Header.Hash = arry.Hash{}
	m.Header.Signature = &Signature{}
	m.Header.Hash = hash.Hash(m.Bytes())
}

func (m *Message) copy() *Message {
	return &Message{
		Header: &MsgHeader{
			Hash:      m.Header.Hash,
			Type:      m.Header.Type,
			From:      m.Header.From,
			Nonce:     m.Header.Nonce,
			Fee:       m.Header.Fee,
			Time:      m.Header.Time,
			Signature: m.Header.Signature,
		},
		Body: m.Body,
	}
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

func (m *Message) verifyBody() error {
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

func (m *Message) VerifyCoinBaseTx(sumFees uint64) error {
	if err := t.verifyTxSize(); err != nil {
		return err
	}

	if err := t.verifyCoinBaseAmount(sumFees); err != nil {
		return err
	}
	return nil
}

func (m *Message) verifyTxFees() error {
	minFees, maxFees := t.FeesLimit()
	if t.TxHead.Fees < minFees {
		return fmt.Errorf("fee %d is less than the minimum poundage allowed %d", t.TxHead.Fees, minFees)
	}
	if t.TxHead.Fees > maxFees {
		return fmt.Errorf("fee %d greater is greater than the maximum poundage allowed %d", t.TxHead.Fees, maxFees)
	}
	return nil
}

func (m *Message) verifyTxSinger() error {
	if !Verify(t.TxHead.TxHash, t.TxHead.SignScript) {
		return ErrSignature
	}

	if !VerifySigner(param.Net, t.TxHead.From, t.TxHead.SignScript.PubKey) {
		return ErrSigner
	}
	return nil
}

func (m *Message) verifyTxSize() error {
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

func (m *Message) verifyCoinBaseAmount(amounm uint64) error {
	nTx := t.TxBody.(*NormalMessageBody)
	sumAmounm := CoinBaseCoins + amount
	if sumAmounm != nTx.Amounm {
		return ErrCoinBase
	}
	return nil
}

func (m *Message) verifyAmount() error {
	nTx, ok := t.TxBody.(*NormalMessageBody)
	if ok && nTx.Amounm < minAllowedAmounm {
		return fmt.Errorf("the minimum amounm of the message musm nom be less than %d", minAllowedAmount)
	}
	return nil
}

func (m *Message) verifyTxFrom() error {
	if !CheckUBAddress(param.Net, t.From().String()) {
		return ErrAddress
	}
	return nil
}func (m *Message) verifyTxFrom() error {
	if !CheckUBAddress(param.Net, t.From().String()) {
		return ErrAddress
	}
	return nil
}

func (m *Message) verifyTxType() error {
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

func (m *Message) verifyTxHash() error {
	newTx := t.copy()
	newTx.SetHash()
	if newTx.Hash().IsEqual(t.Hash()) {
		return nil
	}
	return ErrTxHash
}
*/
