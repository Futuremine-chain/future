package types

import (
	"github.com/Futuremine-chain/future/tools/arry"
	"github.com/Futuremine-chain/future/tools/rlp"
)

type MsgIndex struct {
	MsgRoot arry.Hash
	Index   uint32
	Height  uint64
}

func DecodeTxIndex(bytes []byte) (*MsgIndex, error) {
	var msgIndex *MsgIndex
	err := rlp.DecodeBytes(bytes, &msgIndex)
	if err != nil {
		return nil, err
	}
	return msgIndex, nil
}

func (t *MsgIndex) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(t)
	return bytes
}

func (t *MsgIndex) GetHeight() uint64 {
	return t.Height
}
