package message

import (
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"time"
)

func NewTransaction(from, to, token string, amount, fee, nonce uint64) *types.Message {
	tx := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Transaction,
			Hash:      arry.Hash{},
			From:      arry.StringToAddress(from),
			Nonce:     nonce,
			Fee:       fee,
			Time:      uint64(time.Now().Unix()),
			Signature: &types.Signature{},
		},
		Body: &types.TransactionBody{
			TokenAddress: arry.StringToAddress(token),
			Receiver:     arry.StringToAddress(to),
			Amount:       amount,
		},
	}
	tx.SetHash()
	return tx
}

func NewCandidate(from string, peerStr string, fee, nonce uint64) *types.Message {
	var peerID types.Peer
	copy(peerID[:], peerStr)
	can := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Candidate,
			Hash:      arry.Hash{},
			From:      arry.StringToAddress(from),
			Nonce:     nonce,
			Fee:       fee,
			Time:      uint64(time.Now().Unix()),
			Signature: &types.Signature{},
		},
		Body: &types.CandidateBody{peerID},
	}
	can.SetHash()
	return can
}

func NewCancel(from string, fee, nonce uint64) *types.Message {
	cancel := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Cancel,
			Hash:      arry.Hash{},
			From:      arry.StringToAddress(from),
			Nonce:     nonce,
			Fee:       fee,
			Time:      uint64(time.Now().Unix()),
			Signature: &types.Signature{},
		},
		Body: &types.CancelBody{},
	}
	cancel.SetHash()
	return cancel
}

func NewVote(from, to string, fee, nonce uint64) *types.Message {
	vote := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Vote,
			Hash:      arry.Hash{},
			From:      arry.StringToAddress(from),
			Nonce:     nonce,
			Fee:       fee,
			Time:      uint64(time.Now().Unix()),
			Signature: &types.Signature{},
		},
		Body: &types.VoteBody{arry.StringToAddress(to)},
	}
	vote.SetHash()
	return vote
}

func NewToken(from, to, tokenAddr string, amount, fee, nonce uint64, name, shorthand string, allowIncrease bool) *types.Message {
	token := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Token,
			Hash:      arry.Hash{},
			From:      arry.StringToAddress(from),
			Nonce:     nonce,
			Fee:       fee,
			Time:      uint64(time.Now().Unix()),
			Signature: &types.Signature{},
		},
		Body: &types.TokenBody{
			TokenAddress:   arry.StringToAddress(tokenAddr),
			Receiver:       arry.StringToAddress(to),
			Name:           name,
			Shorthand:      shorthand,
			IncreaseIssues: allowIncrease,
			Amount:         amount,
		},
	}
	token.SetHash()
	return token
}
