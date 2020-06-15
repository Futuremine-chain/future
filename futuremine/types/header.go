package types

import (
	"github.com/Futuremine-chain/futuremine/futuremine/common/param"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	hash2 "github.com/Futuremine-chain/futuremine/tools/crypto/hash"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/types"
)

type Header struct {
	version   int32
	hash      arry.Hash
	preHash   arry.Hash
	msgRoot   arry.Hash
	actRoot   arry.Hash
	dPosRoot  arry.Hash
	tokenRoot arry.Hash
	height    uint64
	time      int64
	cycle     int64
	signer    arry.Address
	signature *Signature
}

func NewHeader(preHash, msgRoot, actRoot, dPosRoot, tokenRoot arry.Hash, height uint64,
	time int64, signer arry.Address) *Header {
	return &Header{
		preHash:   preHash,
		msgRoot:   msgRoot,
		actRoot:   actRoot,
		dPosRoot:  dPosRoot,
		tokenRoot: tokenRoot,
		height:    height,
		time:      time,
		cycle:     time / int64(param.CycleInterval),
		signer:    signer,
	}
}

func DecodeHeader(bytes []byte) (*Header, error) {
	var h *Header
	if err := rlp.DecodeBytes(bytes, h); err != nil {
		return h, err
	}
	return h, nil
}

func (h *Header) Signer() arry.Address {
	return h.signer
}

func (h *Header) Hash() arry.Hash {
	return h.hash
}

func (h *Header) PreHash() arry.Hash {
	return h.preHash
}

func (h *Header) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(h)
	return bytes
}

func (h *Header) Height() uint64 {
	return h.height
}

func (h *Header) MsgRoot() arry.Hash {
	return h.msgRoot
}

func (h *Header) ActRoot() arry.Hash {
	return h.actRoot
}

func (h *Header) DPosRoot() arry.Hash {
	return h.dPosRoot
}

func (h *Header) TokenRoot() arry.Hash {
	return h.tokenRoot
}

func (h *Header) Signature() types.ISignature {
	return h.signature
}

func (h *Header) Time() int64 {
	return h.time
}

func (h *Header) Cycle() int64 {
	return h.cycle
}

func (h *Header) SetHash() {
	h.hash = hash2.Hash(h.Bytes())
}

func (h *Header) Sign(key *secp256k1.PrivateKey) error {
	sig, err := Sign(key, h.hash)
	if err != nil {
		return err
	}
	h.signature = sig
	return nil
}

func (h *Header) ToRlpHeader() types.IRlpHeader {
	return nil
}
