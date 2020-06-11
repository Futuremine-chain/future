package types

import "github.com/Futuremine-chain/futuremine/tools/arry"

type MsgIndex struct {
	MsgRoot arry.Hash
	Index   uint32
}

func DecodeTxIndex([]byte) (*MsgIndex, error) {
	return nil, nil
}

func (t *MsgIndex) Bytes() []byte {
	return nil
}
