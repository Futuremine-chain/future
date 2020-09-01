package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Futuremine-chain/future/tools/arry"
	"github.com/Futuremine-chain/future/types"
)

type IRpcMessageBody interface {
}

type RpcMessageHeader struct {
	MsgHash   string        `json:"msghash"`
	Type      MessageType   `json:"type"`
	From      string        `json:"from"`
	Nonce     uint64        `json:"nonce"`
	Fee       uint64        `json:"fee"`
	Time      uint64        `json:"time"`
	Signature *RpcSignature `json:"signscript"`
}

type RpcMessage struct {
	MsgHeader *RpcMessageHeader `json:"msgheader"`
	MsgBody   IRpcMessageBody   `json:"msgbody"`
	/*TxBody        *RpcTransactionBody `json:"txbody"`
	TokenBody     *RpcTokenBody       `json:"tokenbody"`
	CandidateBody *RpcCandidateBody   `json:"candidatebody"`
	CancelBody    *RpcCancelBody      `json:"cancelbody"`
	VoteBody      *RpcVoteBody        `json:"votebody"`*/
}

type RpcSignature struct {
	Signature string `json:"signature"`
	PubKey    string `json:"pubkey"`
}

func RpcMsgToMsg(rpcMsg *RpcMessage) (*Message, error) {
	var err error
	if rpcMsg.MsgHeader == nil {
		return nil, errors.New("message header is nil")
	}
	signScript, err := RpcSignatureToSignature(rpcMsg.MsgHeader.Signature)
	if err != nil {
		return nil, err
	}
	var msgBody types.IMessageBody
	switch rpcMsg.MsgHeader.Type {
	case Transaction:
		body := &RpcTransactionBody{}
		bytes, err := json.Marshal(rpcMsg.MsgBody)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, body)
		if err != nil {
			return nil, err
		}
		msgBody, err = RpcTransactionBodyToBody(body)
	case Token:
		body := &RpcTokenBody{}
		bytes, err := json.Marshal(rpcMsg.MsgBody)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, body)
		if err != nil {
			return nil, err
		}
		msgBody, err = RpcTokenBodyToBody(body)
	case Candidate:
		body := &RpcCandidateBody{}
		bytes, err := json.Marshal(rpcMsg.MsgBody)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, body)
		if err != nil {
			return nil, err
		}
		msgBody, err = RpcCandidateBodyToBody(body)
	case Cancel:
		msgBody = &CancelBody{}
	case Vote:
		body := &RpcVoteBody{}
		bytes, err := json.Marshal(rpcMsg.MsgBody)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, body)
		if err != nil {
			return nil, err
		}
		msgBody, err = RpcVoteBodyToBody(body)
	}
	hash, err := arry.StringToHash(rpcMsg.MsgHeader.MsgHash)
	if err != nil {
		return nil, fmt.Errorf("wrong message hash %s", rpcMsg.MsgHeader.MsgHash)
	}
	tx := &Message{
		Header: &MsgHeader{
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

func MsgToRpcMsg(msg types.IMessage) (*RpcMessage, error) {
	rpcMsg := &RpcMessage{
		MsgHeader: &RpcMessageHeader{
			MsgHash: msg.Hash().String(),
			Type:    MessageType(msg.Type()),
			From:    addressToString(msg.From()),
			Nonce:   msg.Nonce(),
			Fee:     msg.Fee(),
			Time:    msg.Time(),
			Signature: &RpcSignature{
				Signature: msg.Signature(),
				PubKey:    msg.PublicKey(),
			}},
		MsgBody: nil,
	}
	switch MessageType(msg.Type()) {
	case Transaction:
		rpcMsg.MsgBody = &RpcTransactionBody{
			Token:  msg.MsgBody().MsgToken().String(),
			To:     msg.MsgBody().MsgTo().String(),
			Amount: msg.MsgBody().MsgAmount(),
		}
	case Token:
		body, ok := msg.MsgBody().(*TokenBody)
		if !ok {
			return nil, errors.New("message type error")
		}

		rpcMsg.MsgBody = &RpcTokenBody{
			Address:        msg.MsgBody().MsgToken().String(),
			Receiver:       msg.MsgBody().MsgTo().String(),
			Name:           body.Name,
			Shorthand:      body.Shorthand,
			IncreaseIssues: body.IncreaseIssues,
			Amount:         msg.MsgBody().MsgAmount(),
		}
	case Candidate:
		body, ok := msg.MsgBody().(*CandidateBody)
		if !ok {
			return nil, errors.New("message type error")
		}
		rpcMsg.MsgBody = &RpcCandidateBody{
			PeerId: body.Peer.String(),
		}
	case Cancel:
		rpcMsg.MsgBody = &RpcCancelBody{}
	case Vote:
		rpcMsg.MsgBody = &RpcVoteBody{To: msg.MsgBody().MsgTo().String()}

	}

	return rpcMsg, nil
}

func RpcSignatureToSignature(rpcSignScript *RpcSignature) (*Signature, error) {
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
	return &Signature{
		Bytes:  signature,
		PubKey: pubKey,
	}, nil
}

func RpcTransactionBodyToBody(rpcBody *RpcTransactionBody) (*TransactionBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong transaction body")
	}

	return &TransactionBody{
		TokenAddress: arry.StringToAddress(rpcBody.Token),
		Receiver:     arry.StringToAddress(rpcBody.To),
		Amount:       rpcBody.Amount,
	}, nil
}

func RpcTokenBodyToBody(rpcBody *RpcTokenBody) (*TokenBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong token body")
	}

	return &TokenBody{
		TokenAddress:   arry.StringToAddress(rpcBody.Address),
		Receiver:       arry.StringToAddress(rpcBody.Receiver),
		Name:           rpcBody.Name,
		Shorthand:      rpcBody.Shorthand,
		IncreaseIssues: rpcBody.IncreaseIssues,
		Amount:         rpcBody.Amount,
	}, nil
}

func RpcCandidateBodyToBody(rpcBody *RpcCandidateBody) (*CandidateBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong candidate body")
	}
	body := &CandidateBody{}
	copy(body.Peer[:], rpcBody.PeerIdBytes())
	return body, nil
}

func RpcVoteBodyToBody(rpcBody *RpcVoteBody) (*VoteBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong vote body")
	}

	return &VoteBody{To: arry.StringToAddress(rpcBody.To)}, nil
}

func addressToString(address arry.Address) string {
	if address.IsEqual(CoinBase) {
		return CoinBase.String()
	}
	return address.String()
}
