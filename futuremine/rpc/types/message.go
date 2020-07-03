package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type IRpcMessageBody interface {
}

type RpcMessageHeader struct {
	MsgHash   string               `json:"txhash"`
	Type      fmctypes.MessageType `json:"type"`
	From      string               `json:"from"`
	Nonce     uint64               `json:"nonce"`
	Fee       uint64               `json:"fee"`
	Time      uint64               `json:"time"`
	Signature *RpcSignature        `json:"signscript"`
}

type RpcMessage struct {
	MsgHeader     *RpcMessageHeader   `json:"txhead"`
	TxBody        *RpcTransactionBody `json:"txbody"`
	TokenBody     *RpcTokenBody       `json:"tokenbody"`
	CandidateBody *RpcCandidateBody   `json:"candidatebody"`
	CancelBody    *RpcCancelBody      `json:"cancelbody"`
	VoteBody      *RpcVoteBody        `json:"votebody"`
}

type RpcSignature struct {
	Signature string `json:"signature"`
	PubKey    string `json:"pubkey"`
}

func RpcMsgToMsg(rpcMsg *RpcMessage) (*fmctypes.Message, error) {
	var err error
	signScript, err := RpcSignatureToSignature(rpcMsg.MsgHeader.Signature)
	if err != nil {
		return nil, err
	}
	var msgBody types.IMessageBody
	switch rpcMsg.MsgHeader.Type {
	case fmctypes.Transaction:
		msgBody, err = RpcTransactionBodyToBody(rpcMsg.TxBody)
	case fmctypes.Token:
		msgBody, err = RpcTokenBodyToBody(rpcMsg.TokenBody)
	case fmctypes.Candidate:
		msgBody, err = RpcCandidateBodyToBody(rpcMsg.CandidateBody)
	case fmctypes.Cancel:
		msgBody = &fmctypes.CancelBody{}
	case fmctypes.Vote:
		msgBody, err = RpcVoteBodyToBody(rpcMsg.VoteBody)
	}
	hash, err := arry.StringToHash(rpcMsg.MsgHeader.MsgHash)
	if err != nil {
		return nil, fmt.Errorf("wrong message hash %s", rpcMsg.MsgHeader.MsgHash)
	}
	tx := &fmctypes.Message{
		Header: &fmctypes.MsgHeader{
			Hash:      hash,
			Type:      rpcMsg.MsgHeader.Type,
			From:      arry.StringToAddress(rpcMsg.MsgHeader.From),
			Nonce:     rpcMsg.MsgHeader.Nonce,
			Fee:       rpcMsg.MsgHeader.Fee,
			Time:      rpcMsg.MsgHeader.Time,
			Signature: signScript,
		},
		Body: msgBody,
	}
	return tx, nil
}

func MsgToRpcMsg(msg *fmctypes.Message) (*RpcMessage, error) {
	rpcMsg := &RpcMessage{
		MsgHeader: &RpcMessageHeader{
			MsgHash: msg.Hash().String(),
			Type:    fmctypes.MessageType(msg.Type()),
			From:    addressToString(msg.From()),
			Nonce:   msg.Nonce(),
			Fee:     msg.Fee(),
			Time:    msg.Time(),
			Signature: &RpcSignature{
				Signature: hex.EncodeToString(msg.Header.Signature.SignatureBytes()),
				PubKey:    hex.EncodeToString(msg.Header.Signature.PubicKey()),
			}},
		TxBody:        nil,
		TokenBody:     nil,
		CandidateBody: nil,
		CancelBody:    nil,
		VoteBody:      nil,
	}
	switch fmctypes.MessageType(msg.Type()) {
	case fmctypes.Transaction:
		rpcMsg.TxBody = &RpcTransactionBody{
			Token:  msg.Body.MsgToken().String(),
			To:     msg.Body.MsgTo().String(),
			Amount: msg.Body.MsgAmount(),
		}
	case fmctypes.Token:
		body, ok := msg.Body.(*fmctypes.TokenBody)
		if !ok {
			return nil, errors.New("message type error")
		}
		rpcMsg.TokenBody = &RpcTokenBody{
			Address:   msg.Body.MsgToken().String(),
			Receiver:  msg.Body.MsgTo().String(),
			Name:      body.Name,
			Shorthand: body.Shorthand,
			Amount:    msg.Body.MsgAmount(),
		}
	case fmctypes.Candidate:
		body, ok := msg.Body.(*fmctypes.CandidateBody)
		if !ok {
			return nil, errors.New("message type error")
		}
		rpcMsg.CandidateBody = &RpcCandidateBody{
			PeerId: body.Peer.String(),
		}
	case fmctypes.Cancel:
		rpcMsg.CancelBody = &RpcCancelBody{}
	case fmctypes.Vote:
		rpcMsg.VoteBody = &RpcVoteBody{To: msg.Body.MsgTo().String()}

	}

	return rpcMsg, nil
}

func RpcSignatureToSignature(rpcSignScript *RpcSignature) (*fmctypes.Signature, error) {
	if rpcSignScript == nil {
		return nil, errors.New("signature is nil")
	}
	if rpcSignScript.Signature == "" || rpcSignScript.PubKey == "" {
		return nil, errors.New("signature content is nil")
	}
	signature, err := hex.DecodeString(rpcSignScript.Signature)
	if err != nil {
		return nil, err
	}
	pubKey, err := hex.DecodeString(rpcSignScript.PubKey)
	if err != nil {
		return nil, err
	}
	return &fmctypes.Signature{
		Bytes:  signature,
		PubKey: pubKey,
	}, nil
}

func RpcTransactionBodyToBody(rpcBody *RpcTransactionBody) (*fmctypes.TransactionBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong transaction body")
	}

	return &fmctypes.TransactionBody{
		TokenAddress: arry.StringToAddress(rpcBody.Token),
		Receiver:     arry.StringToAddress(rpcBody.To),
		Amount:       rpcBody.Amount,
	}, nil
}

func RpcTokenBodyToBody(rpcBody *RpcTokenBody) (*fmctypes.TokenBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong token body")
	}

	return &fmctypes.TokenBody{
		TokenAddress: arry.StringToAddress(rpcBody.Address),
		Receiver:     arry.StringToAddress(rpcBody.Receiver),
		Name:         rpcBody.Name,
		Shorthand:    rpcBody.Shorthand,
		Amount:       rpcBody.Amount,
	}, nil
}

func RpcCandidateBodyToBody(rpcBody *RpcCandidateBody) (*fmctypes.CandidateBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong candidate body")
	}
	body := &fmctypes.CandidateBody{}
	copy(body.Peer[:], rpcBody.PeerIdBytes())
	return body, nil
}

func RpcVoteBodyToBody(rpcBody *RpcVoteBody) (*fmctypes.VoteBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong vote body")
	}

	return &fmctypes.VoteBody{To: arry.StringToAddress(rpcBody.To)}, nil
}

func addressToString(address arry.Address) string {
	if address.IsEqual(fmctypes.CoinBase) {
		return fmctypes.CoinBase.String()
	}
	return address.String()
}
