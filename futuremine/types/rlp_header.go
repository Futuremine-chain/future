package types

import "github.com/Futuremine-chain/futuremine/tools/rlp"

type RlpHeader Header

func (r *RlpHeader) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(r)
	return bytes
}

func DecodeRlpHeader(bytes []byte) (*RlpHeader, error) {
	var h = new(RlpHeader)
	if err := rlp.DecodeBytes(bytes, h); err != nil {
		return h, err
	}
	return h, nil
}
