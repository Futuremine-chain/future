package types

import (
	"bytes"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/hash"
	"github.com/Futuremine-chain/futuremine/types"
)

type Messages []*Message

func (m Messages) MsgList() []types.IMessage {
	iTxs := make([]types.IMessage, len(m))
	for i, msg := range m {
		iTxs[i] = msg
	}
	return iTxs
}

func (m Messages) Add(iMsg types.IMessage) {
	iMsg = new(Message)
	msg := iMsg.(*Message)
	m = append(m, msg)
}

func (m Messages) Count() int {
	return len(m)
}

func (m Messages) MsgRoot() arry.Hash {
	var hashes [][]byte
	for _, msg := range m {
		hashes = append(hashes, msg.Hash().Bytes())
	}
	hashBytes := bytes.Join(hashes, []byte{})
	return hash.Hash(hashBytes)
}

func (m Messages) CalculateFee() uint64 {
	var sum uint64
	for _, msg := range m {
		sum += msg.Fee()
	}
	return sum
}
