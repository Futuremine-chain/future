package types

import "github.com/Futuremine-chain/futuremine/tools/rlp"

type RlpHeader Header

func (r *RlpHeader) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(r)
	return bytes
}
