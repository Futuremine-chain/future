package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
)

type MsgIndex struct {
	MsgRoot arry.Hash
	Index   uint32
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
