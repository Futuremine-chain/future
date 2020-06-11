package types

type IMessages interface {
	Add(msg IMessage)
	Msgs() []IMessage
}
