package types

import (
	"errors"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	"github.com/Futuremine-chain/futuremine/tools/crypto/hash"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/types"
)

const CoinBase = "CoinBase"

type Message struct {
	Header *MsgHeader
	Body   types.IMessageBody
}

func (m *Message) Type() int {
	return int(m.Header.Type)
}

func (m *Message) Hash() arry.Hash {
	return m.Header.Hash
}

func (m *Message) From() arry.Address {
	return m.Header.From
}

func (m *Message) Nonce() uint64 {
	return m.Header.Nonce
}

func (m *Message) Fee() uint64 {
	return m.Header.Fee
}

func (m *Message) Time() int64 {
	return m.Header.Time
}

func (m *Message) IsCoinBase() bool {
	return m.Header.From.IsEqual(arry.StringToAddress(CoinBase))
}

func (m *Message) To() arry.Address {
	return m.Body.To()
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
