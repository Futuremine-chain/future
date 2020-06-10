package types

import "github.com/Futuremine-chain/futuremine/tools/arry"

type TxIndex struct {
	TxRoot  arry.Hash
	TxIndex uint32
}

func DecodeTxIndex([]byte) (*TxIndex, error) {
	return nil, nil
}

func (t *TxIndex) Bytes() []byte {
	return nil
}
