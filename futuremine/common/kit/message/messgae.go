package message

import (
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"time"
)

func NewTransaction(from, to, token arry.Address, amount, fee, nonce uint64) *types.Message {
	tx := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Transaction,
			Hash:      arry.Hash{},
			From:      from,
			Nonce:     nonce,
			Fee:       fee,
			Time:      uint64(time.Now().Unix()),
			Signature: &types.Signature{},
		},
		Body: &types.TransactionBody{
			TokenAddress: token,
			Receiver:     to,
			Amount:       amount,
		},
	}
	tx.SetHash()
	return tx
}

func NewCandidate(from arry.Address, peerStr string, fee, nonce uint64) *types.Message {
	var peerID types.Peer
	copy(peerID[:], peerStr)
	can := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Candidate,
			Hash:      arry.Hash{},
			From:      from,
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

func NewCancel(from arry.Address, fee, nonce uint64) *types.Message {
	cancel := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Cancel,
			Hash:      arry.Hash{},
			From:      from,
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

func NewVote(from, to arry.Address, fee, nonce uint64) *types.Message {
	vote := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Vote,
			Hash:      arry.Hash{},
			From:      from,
			Nonce:     nonce,
			Fee:       fee,
			Time:      uint64(time.Now().Unix()),
			Signature: &types.Signature{},
		},
		Body: &types.VoteBody{to},
	}
	vote.SetHash()
	return vote
}

func NewToken(from, to, tokenAddr arry.Address, amount, fee, nonce uint64, name, shorthand string) *types.Message {
	token := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Token,
			Hash:      arry.Hash{},
			From:      from,
			Nonce:     nonce,
			Fee:       fee,
			Time:      uint64(time.Now().Unix()),
			Signature: &types.Signature{},
		},
		Body: &types.TokenBody{
			TokenAddress: tokenAddr,
			Receiver:     to,
			Name:         name,
			Shorthand:    shorthand,
			Amount:       amount,
		},
	}
	token.SetHash()
	return token
}
