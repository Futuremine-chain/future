package types

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/futuremine/common/param"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	hash2 "github.com/Futuremine-chain/futuremine/tools/crypto/hash"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/types"
	"time"
)

type Header struct {
	version   uint32
	hash      arry.Hash
	preHash   arry.Hash
	msgRoot   arry.Hash
	actRoot   arry.Hash
	dPosRoot  arry.Hash
	tokenRoot arry.Hash
	height    uint64
	time      time.Time
	cycle     uint64
	signer    arry.Address
	signature *Signature
}

func NewHeader(preHash, msgRoot, actRoot, dPosRoot, tokenRoot arry.Hash, height uint64,
	blockTime int64, signer arry.Address) *Header {
	return &Header{
		preHash:   preHash,
		msgRoot:   msgRoot,
		actRoot:   actRoot,
		dPosRoot:  dPosRoot,
		tokenRoot: tokenRoot,
		height:    height,
		time:      time.Unix(blockTime, 0),
		cycle:     uint64(blockTime) / uint64(param.CycleInterval),
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
	bytes, err := rlp.EncodeToBytes(h)
	fmt.Println(err)
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
	return h.time.Unix()
}

func (h *Header) Cycle() int64 {
	return int64(h.cycle)
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
