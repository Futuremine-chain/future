package types

type IRlpBlock interface {
	Bytes() []byte
	BytesToBlock() IBlock
}
