package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/types"
)

type Header struct {
	hash   arry.Hash
	txHash arry.Hash
	height uint64
	time   int64
}

func DecodeHeader([]byte) (*Header, error) {
	return nil, nil
}

func (h *Header) Signer() arry.Address {
	panic("implement me")
}

func (h *Header) Hash() arry.Hash {
	return h.Hash()
}

func (h *Header) Bytes() []byte {
	return nil
}

func (h *Header) Height() uint64 {
	return h.height
}

func (h *Header) MsgRoot() arry.Hash {
	return h.txHash
}

func (h *Header) Time() int64 {
	return h.time
}

func (h *Header) ToRlpHeader() types.IRlpHeader {
	return nil
}
