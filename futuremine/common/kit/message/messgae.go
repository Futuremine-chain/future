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
	tx := &types.Message{
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
	tx.SetHash()
	return tx
}

func NewCancel(from arry.Address, fee, nonce uint64) *types.Message {
	tx := &types.Message{
		Header: &types.MsgHeader{
			Type:      types.Cancel,
			Hash:      arry.Hash{},
			From:      from,
			Nonce:     nonce,
			Fee:       fee,
			Time:      uint64(time.Now().Unix()),
			Signature: &types.Signature{},
		},
		Body: &types.CandidateBody{},
	}
	tx.SetHash()
	return tx
}

func NewVote(from, to arry.Address, fee, nonce uint64) *types.Message {
	tx := &types.Message{
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
	tx.SetHash()
	return tx
}
