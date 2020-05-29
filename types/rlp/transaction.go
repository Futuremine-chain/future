package rlp

type IRlpTransaction interface {
	Bytes() []byte
}
