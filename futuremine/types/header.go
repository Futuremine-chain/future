package types

import "github.com/Futuremine-chain/futuremine/tools/arry"

type Header struct {
	hash   arry.Hash
	height uint64
}

func DecodeHeader([]byte) (*Header, error) {
	return nil, nil
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
