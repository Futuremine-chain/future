package types

import "github.com/Futuremine-chain/futuremine/tools/rlp"

type RlpHeader struct {
}

func (r *RlpHeader) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(r)
	return bytes
}
